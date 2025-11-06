package cicd

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestDeployAppWorkflow(t *testing.T) {
	t.Parallel()

	t.Run("No Apps", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := DeployAppWorkflow(t.Context(), input)
		assert.NoError(t, err)

		// No workflow should be created when no apps exist.
		exists, err := afero.Exists(input.FS, filepath.Join(workflowsPath, "deploy-app.yaml"))
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("App Without Dockerfile", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Build: appdef.Build{
						Port: 3000,
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "container",
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := DeployAppWorkflow(t.Context(), input)
		assert.NoError(t, err)

		// No workflow should be created when app has no Dockerfile.
		exists, err := afero.Exists(input.FS, filepath.Join(workflowsPath, "deploy-app.yaml"))
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Single Container App", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:  "web",
					Title: "Web",
					Type:  appdef.AppTypeSvelteKit,
					Path:  "web",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       3000,
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "container",
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := DeployAppWorkflow(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "deploy-app.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		content := string(file)

		t.Log("Workflow metadata")
		{
			assert.Contains(t, content, "name: Deploy App")
			assert.Contains(t, content, "workflow_dispatch:")
		}

		t.Log("Inputs")
		{
			assert.Contains(t, content, "app_name:")
			assert.Contains(t, content, "image_tag:")
			assert.Contains(t, content, "default: 'latest'")
		}

		t.Log("App in options")
		{
			assert.Contains(t, content, "- web")
		}

		t.Log("Router job for container deployment")
		{
			assert.Contains(t, content, "deploy-container:")
			assert.Contains(t, content, "uses: ./.github/workflows/deploy-container.yaml")
			assert.Contains(t, content, "app_name: ${{ github.event.inputs.app_name }}")
		}

		t.Log("No setup-webkit job for container-only apps")
		{
			assert.NotContains(t, content, "setup-webkit:")
		}

		t.Log("No direct deployment logic (router pattern)")
		{
			assert.NotContains(t, content, "curl -X POST")
			assert.NotContains(t, content, "dawidd6/action-ansible-playbook")
		}
	})

	t.Run("Single VM App", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:  "cms",
					Title: "Cms",
					Type:  appdef.AppTypePayload,
					Path:  "cms",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       3000,
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "vm",
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := DeployAppWorkflow(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "deploy-app.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		content := string(file)

		t.Log("Workflow metadata")
		{
			assert.Contains(t, content, "name: Deploy App")
			assert.Contains(t, content, "workflow_dispatch:")
		}

		t.Log("Inputs")
		{
			assert.Contains(t, content, "app_name:")
			assert.Contains(t, content, "image_tag:")
			assert.Contains(t, content, "webkit_version:")
			assert.Contains(t, content, "default: 'latest'")
		}

		t.Log("App in options")
		{
			assert.Contains(t, content, "- cms")
		}

		t.Log("Setup webkit job for VM apps")
		{
			assert.Contains(t, content, "setup-webkit:")
			assert.Contains(t, content, "if: github.event.inputs.app_name == 'cms'")
		}

		t.Log("Router job for VM deployment")
		{
			assert.Contains(t, content, "deploy-vm:")
			assert.Contains(t, content, "uses: ./.github/workflows/deploy-vm.yaml")
			assert.Contains(t, content, "webkit_version: ${{ needs.setup-webkit.outputs.version }}")
		}

		t.Log("No direct deployment logic (router pattern)")
		{
			assert.NotContains(t, content, "dawidd6/action-ansible-playbook@v4")
		}
	})

	t.Run("Mixed Container and VM Apps", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:  "web",
					Title: "Web",
					Type:  appdef.AppTypeSvelteKit,
					Path:  "web",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       3000,
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "container",
					},
				},
				{
					Name:  "cms",
					Title: "Cms",
					Type:  appdef.AppTypePayload,
					Path:  "cms",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       3001,
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "vm",
					},
				},
				{
					Name:  "api",
					Title: "Api",
					Type:  appdef.AppTypeGoLang,
					Path:  "api",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       8080,
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "container",
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := DeployAppWorkflow(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "deploy-app.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		content := string(file)

		t.Log("All apps in options")
		{
			assert.Contains(t, content, "- web")
			assert.Contains(t, content, "- cms")
			assert.Contains(t, content, "- api")
		}

		t.Log("Router job for container deployments")
		{
			assert.Contains(t, content, "deploy-container:")
			assert.Contains(t, content, "uses: ./.github/workflows/deploy-container.yaml")
			assert.Contains(t, content, "if: github.event.inputs.app_name == 'web' || github.event.inputs.app_name == 'api'")
		}

		t.Log("Router job for VM deployments")
		{
			assert.Contains(t, content, "deploy-vm:")
			assert.Contains(t, content, "uses: ./.github/workflows/deploy-vm.yaml")
			assert.Contains(t, content, "if: github.event.inputs.app_name == 'cms'")
		}

		t.Log("Setup webkit job exists for VM apps")
		{
			assert.Contains(t, content, "setup-webkit:")
			assert.Contains(t, content, "if: github.event.inputs.app_name == 'cms'")
		}

		t.Log("No direct deployment logic (router pattern)")
		{
			assert.NotContains(t, content, "curl -X POST")
			assert.NotContains(t, content, "dawidd6/action-ansible-playbook@v4")
		}
	})

	t.Run("Non-DigitalOcean Provider", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       3000,
					},
					Infra: appdef.Infra{
						Provider: "aws",
						Type:     "container",
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := DeployAppWorkflow(t.Context(), input)
		assert.NoError(t, err)

		// No workflow should be created for non-DigitalOcean providers.
		exists, err := afero.Exists(input.FS, filepath.Join(workflowsPath, "deploy-app.yaml"))
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "container",
					},
				},
			},
		}

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), appDef)

		err := DeployAppWorkflow(t.Context(), input)
		assert.Error(t, err)
	})
}
