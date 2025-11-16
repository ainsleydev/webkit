package cicd

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func TestReleaseWorkflow(t *testing.T) {
	t.Parallel()

	t.Run("No Apps", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := ReleaseWorkflow(t.Context(), input)
		assert.NoError(t, err)

		// No workflow should be created when no apps exist.
		exists, err := afero.Exists(input.FS, filepath.Join(workflowsPath, "release.yaml"))
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Apps Without Dockerfile", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "cms",
					Build: appdef.Build{
						Dockerfile: "",
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := ReleaseWorkflow(t.Context(), input)
		assert.NoError(t, err)

		// No workflow should be created when apps have no Dockerfile.
		exists, err := afero.Exists(input.FS, filepath.Join(workflowsPath, "release.yaml"))
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Apps With Release Disabled", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "cms",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Release:    ptr.BoolPtr(false),
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := ReleaseWorkflow(t.Context(), input)
		assert.NoError(t, err)

		// No workflow should be created when all apps have release disabled.
		exists, err := afero.Exists(input.FS, filepath.Join(workflowsPath, "release.yaml"))
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Single App", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "cms",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       3000,
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := ReleaseWorkflow(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "release.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		content := string(file)

		t.Log("Workflow metadata")
		{
			assert.Contains(t, content, "name: Release")
			assert.Contains(t, content, "on:")
			assert.Contains(t, content, "push:")
			assert.Contains(t, content, "branches:")
			assert.Contains(t, content, "- main")
		}

		t.Log("Jobs")
		{
			assert.Contains(t, content, "build-and-push:")
			assert.Contains(t, content, "cleanup-containers:")
		}

		t.Log("App in matrix")
		{
			assert.Contains(t, content, "- name: cms")
			assert.Contains(t, content, "context: ./cms")
			assert.Contains(t, content, "dockerfile: ./cms/Dockerfile")
		}

		t.Log("Docker build steps")
		{
			assert.Contains(t, content, "Set up QEMU")
			assert.Contains(t, content, "Set up Docker Buildx")
			assert.Contains(t, content, "Log in to GitHub Container Registry")
			assert.Contains(t, content, "Build and push Docker image")
		}

		t.Log("Cleanup job")
		{
			assert.Contains(t, content, "Delete old images from GHCR")
			assert.Contains(t, content, "min-versions-to-keep: 5")
		}
	})

	t.Run("Multiple Apps", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "cms",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       3000,
					},
				},
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "web",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       3001,
					},
				},
				{
					Name: "api",
					Type: appdef.AppTypeGoLang,
					Path: "api",
					Build: appdef.Build{
						Dockerfile: "Dockerfile.production",
						Port:       8080,
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := ReleaseWorkflow(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "release.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		content := string(file)

		t.Log("All apps in matrix")
		{
			assert.Contains(t, content, "- name: cms")
			assert.Contains(t, content, "context: ./cms")
			assert.Contains(t, content, "dockerfile: ./cms/Dockerfile")

			assert.Contains(t, content, "- name: web")
			assert.Contains(t, content, "context: ./web")
			assert.Contains(t, content, "dockerfile: ./web/Dockerfile")

			assert.Contains(t, content, "- name: api")
			assert.Contains(t, content, "context: ./api")
			assert.Contains(t, content, "dockerfile: ./api/Dockerfile.production")
		}

		t.Log("Cleanup matrix includes all apps")
		{
			assert.Contains(t, content, "service: [cms, web, api]")
		}
	})

	t.Run("Mixed Release Flags", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "cms",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Release:    ptr.BoolPtr(true),
					},
				},
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "web",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Release:    ptr.BoolPtr(false),
					},
				},
				{
					Name: "api",
					Type: appdef.AppTypeGoLang,
					Path: "api",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Release:    nil, // Defaults to true.
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := ReleaseWorkflow(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "release.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		content := string(file)

		t.Log("Only enabled apps in matrix")
		{
			assert.Contains(t, content, "- name: cms")
			assert.Contains(t, content, "- name: api")
			assert.NotContains(t, content, "- name: web")
		}

		t.Log("Cleanup matrix only includes enabled apps")
		{
			assert.Contains(t, content, "service: [cms, api]")
		}
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "cms",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
					},
				},
			},
		}

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), appDef)

		err := ReleaseWorkflow(t.Context(), input)
		assert.Error(t, err)
	})

	t.Run("DigitalOcean App Name Uses Project Name", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "my-awesome-project",
				Repo: appdef.GitHubRepo{
					Owner: "acme",
					Name:  "myawesomeproject",
				},
			},
			Apps: []appdef.App{
				{
					Name:  "api",
					Title: "API",
					Type:  appdef.AppTypeGoLang,
					Path:  "api",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       8080,
					},
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "container",
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := ReleaseWorkflow(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "release.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		content := string(file)

		t.Log("DigitalOcean app_name uses project.name, NOT repo name")
		{
			// The app_name should use project.name (my-awesome-project) to match Terraform.
			// Terraform constructs the app name as: "${var.project_name}-${var.name}".
			// This ensures GitHub Actions and Terraform use the same naming convention.
			assert.Contains(t, content, "app_name: 'my-awesome-project-api'")

			// Should NOT use the GitHub repo name (myawesomeproject).
			assert.NotContains(t, content, "app_name: 'myawesomeproject-api'")
		}

		t.Log("Docker image repository uses repo name variable")
		{
			// The image repository should use github.event.repository.name to match GHCR.
			// GitHub Actions publishes to: ghcr.io/{owner}/{repo-name}-{app-name}.
			// We check for the template variable, not the hardcoded value.
			assert.Contains(t, content, "images: ghcr.io/${{ github.repository_owner }}/${{ github.event.repository.name }}-${{ matrix.service.name }}")
		}
	})

	t.Run("Terraform Apply Job", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:  "web",
					Title: "Web",
					Type:  appdef.AppTypeGoLang,
					Path:  "web",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       8080,
					},
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "container",
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := ReleaseWorkflow(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "release.yaml"))
		require.NoError(t, err)

		content := string(file)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		t.Log("Terraform Apply Job")
		{
			assert.Contains(t, content, "terraform-apply-production:")
			assert.Contains(t, content, "Run Terraform Plan")
			assert.Contains(t, content, "Run Terraform Apply")
			assert.Contains(t, content, "./webkit infra plan")
			assert.Contains(t, content, "./webkit infra apply")
			assert.Contains(t, content, "needs: [setup-webkit, build-and-push]")
			assert.Contains(t, content, "Setup Infrastructure Dependencies")
			assert.Contains(t, content, "./.github/actions/setup-infra")
			assert.Contains(t, content, "Send Slack Notification")
			assert.Contains(t, content, "Infra Plan")
		}

		t.Log("Diff-based change detection")
		{
			assert.Contains(t, content, "Check if Terraform needed")
			assert.Contains(t, content, "id: diff")
			assert.Contains(t, content, "./webkit infra diff --format=github")
			assert.Contains(t, content, "if: steps.diff.outputs.skip_terraform != 'true'")
			assert.Contains(t, content, "Skipping Terraform - no infrastructure changes detected")
			assert.Contains(t, content, "Running Terraform - infrastructure changes detected")
		}

		t.Log("Deploy jobs depend on terraform apply")
		{
			assert.Contains(t, content, "needs: [build-and-push, terraform-apply-production]")
		}
	})
}
