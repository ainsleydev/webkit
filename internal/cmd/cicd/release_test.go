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
}
