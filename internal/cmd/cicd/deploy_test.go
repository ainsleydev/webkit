package cicd

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestDeployContainerWorkflow(t *testing.T) {
	t.Parallel()

	t.Run("No Container Apps", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "api",
					Type: appdef.AppTypeGoLang,
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       8080,
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "vm",
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := DeployContainerWorkflow(t.Context(), input)
		assert.NoError(t, err)

		// No workflow should be created when no container apps exist.
		exists, err := afero.Exists(input.FS, filepath.Join(workflowsPath, "deploy-container.yaml"))
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

		err := DeployContainerWorkflow(t.Context(), input)
		assert.NoError(t, err)

		// No workflow should be created when app has no Dockerfile.
		exists, err := afero.Exists(input.FS, filepath.Join(workflowsPath, "deploy-container.yaml"))
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Single Container App", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "web",
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

		err := DeployContainerWorkflow(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "deploy-container.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		content := string(file)

		t.Log("Workflow metadata")
		{
			assert.Contains(t, content, "name: Deploy Container App")
			assert.Contains(t, content, "workflow_dispatch:")
			assert.Contains(t, content, "workflow_call:")
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

		t.Log("Deploy job")
		{
			assert.Contains(t, content, "deploy-web:")
			assert.Contains(t, content, "uses: digitalocean/app_action/deploy@v2")
		}
	})

	t.Run("Multiple Container Apps", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "web",
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
					Name: "api",
					Type: appdef.AppTypeGoLang,
					Path: "api",
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

		err := DeployContainerWorkflow(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "deploy-container.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		content := string(file)

		t.Log("All apps in options")
		{
			assert.Contains(t, content, "- web")
			assert.Contains(t, content, "- api")
		}

		t.Log("All deploy jobs")
		{
			assert.Contains(t, content, "deploy-web:")
			assert.Contains(t, content, "deploy-api:")
		}
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

		err := DeployContainerWorkflow(t.Context(), input)
		assert.Error(t, err)
	})
}

func TestDeployVMWorkflow(t *testing.T) {
	t.Parallel()

	t.Run("No VM Apps", func(t *testing.T) {
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
						Provider: "digitalocean",
						Type:     "container",
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := DeployVMWorkflow(t.Context(), input)
		assert.NoError(t, err)

		// No workflow should be created when no VM apps exist.
		exists, err := afero.Exists(input.FS, filepath.Join(workflowsPath, "deploy-vm.yaml"))
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("App Without Dockerfile", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Build: appdef.Build{
						Port: 3000,
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "vm",
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := DeployVMWorkflow(t.Context(), input)
		assert.NoError(t, err)

		// No workflow should be created when app has no Dockerfile.
		exists, err := afero.Exists(input.FS, filepath.Join(workflowsPath, "deploy-vm.yaml"))
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Single VM App", func(t *testing.T) {
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
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "vm",
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := DeployVMWorkflow(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "deploy-vm.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		content := string(file)

		t.Log("Workflow metadata")
		{
			assert.Contains(t, content, "name: Deploy VM App")
			assert.Contains(t, content, "workflow_dispatch:")
			assert.Contains(t, content, "workflow_call:")
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

		t.Log("Setup and deploy jobs")
		{
			assert.Contains(t, content, "setup-webkit:")
			assert.Contains(t, content, "deploy-cms:")
			assert.Contains(t, content, "uses: dawidd6/action-ansible-playbook@v4")
		}
	})

	t.Run("Multiple VM Apps", func(t *testing.T) {
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
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "vm",
					},
				},
				{
					Name: "admin",
					Type: appdef.AppTypePayload,
					Path: "admin",
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
						Port:       3001,
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "vm",
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := DeployVMWorkflow(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "deploy-vm.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		content := string(file)

		t.Log("All apps in options")
		{
			assert.Contains(t, content, "- cms")
			assert.Contains(t, content, "- admin")
		}

		t.Log("All deploy jobs")
		{
			assert.Contains(t, content, "deploy-cms:")
			assert.Contains(t, content, "deploy-admin:")
		}
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Build: appdef.Build{
						Dockerfile: "Dockerfile",
					},
					Infra: appdef.Infra{
						Provider: "digitalocean",
						Type:     "vm",
					},
				},
			},
		}

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), appDef)

		err := DeployVMWorkflow(t.Context(), input)
		assert.Error(t, err)
	})
}
