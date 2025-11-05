package cicd

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestVMMaintenanceWorkflow(t *testing.T) {
	t.Parallel()

	t.Run("No Apps", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Apps: []appdef.App{},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		got := VMMaintenanceWorkflow(t.Context(), input)
		assert.NoError(t, got)

		// Workflow file should not be created when there are no apps
		exists, err := afero.Exists(input.FS, filepath.Join(workflowsPath, "server-maintenance.yaml"))
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("No Digital Ocean VM Apps", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeGoLang,
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "container",
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		got := VMMaintenanceWorkflow(t.Context(), input)
		assert.NoError(t, got)

		// Workflow file should not be created when there are no VM apps
		exists, err := afero.Exists(input.FS, filepath.Join(workflowsPath, "server-maintenance.yaml"))
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Single Digital Ocean VM App", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Apps: []appdef.App{
				{
					Name:  "api",
					Title: "API Server",
					Type:  appdef.AppTypeGoLang,
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "vm",
						Config: map[string]any{
							"admin_email": "admin@example.com",
						},
					},
					Domains: []appdef.Domain{
						{
							Name: "api.example.com",
							Type: appdef.DomainTypePrimary,
						},
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		got := VMMaintenanceWorkflow(t.Context(), input)
		assert.NoError(t, got)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "server-maintenance.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		// Verify the file contains expected content
		content := string(file)
		assert.Contains(t, content, "name: Server Maintenance")
		assert.Contains(t, content, "maintenance-vm-api")
		assert.Contains(t, content, "API Server")
	})

	t.Run("Multiple Apps Mixed Types", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Apps: []appdef.App{
				{
					Name:  "api",
					Title: "API Server",
					Type:  appdef.AppTypeGoLang,
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "vm",
					},
					Domains: []appdef.Domain{
						{
							Name: "api.example.com",
							Type: appdef.DomainTypePrimary,
						},
					},
				},
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "container",
					},
				},
				{
					Name:  "worker",
					Title: "Background Worker",
					Type:  appdef.AppTypeGoLang,
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "vm",
					},
					Domains: []appdef.Domain{
						{
							Name: "worker.example.com",
							Type: appdef.DomainTypePrimary,
						},
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		got := VMMaintenanceWorkflow(t.Context(), input)
		assert.NoError(t, got)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "server-maintenance.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		// Verify both VM apps are included but container app is not
		content := string(file)
		assert.Contains(t, content, "maintenance-vm-api")
		assert.Contains(t, content, "maintenance-vm-worker")
		assert.NotContains(t, content, "maintenance-vm-web")
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Apps: []appdef.App{
				{
					Name: "api",
					Type: appdef.AppTypeGoLang,
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "vm",
					},
				},
			},
		}

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), appDef)

		got := VMMaintenanceWorkflow(t.Context(), input)
		assert.Error(t, got)
	})
}
