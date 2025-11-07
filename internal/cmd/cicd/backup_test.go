package cicd

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestBackupWorkflow(t *testing.T) {
	t.Parallel()

	t.Run("No Resources", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Resources: []appdef.Resource{},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		got := BackupWorkflow(t.Context(), input)
		assert.NoError(t, got)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "backup.yaml"))
		require.NoError(t, err)
		require.NotEmpty(t, file)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)
	})

	t.Run("Postgres", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Resources: []appdef.Resource{
				{
					Name:     "db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
					Backup: appdef.ResourceBackupConfig{
						Enabled: true,
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		got := BackupWorkflow(t.Context(), input)
		assert.NoError(t, got)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "backup.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)
	})

	t.Run("S3", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Resources: []appdef.Resource{
				{
					Name:     "store",
					Type:     appdef.ResourceTypeS3,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"key": "value",
					},
					Backup: appdef.ResourceBackupConfig{
						Enabled: true,
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		got := BackupWorkflow(t.Context(), input)
		assert.NoError(t, got)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "backup.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)
	})

	t.Run("S3 Non DigitalOcean Skipped", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Resources: []appdef.Resource{
				{
					Name:     "b2store",
					Type:     appdef.ResourceTypeS3,
					Provider: appdef.ResourceProviderBackBlaze,
					Backup: appdef.ResourceBackupConfig{
						Enabled: true,
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		got := BackupWorkflow(t.Context(), input)
		assert.NoError(t, got)

		// File should be created but empty (no jobs) since B2 provider is not supported
		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "backup.yaml"))
		require.NoError(t, err)
		assert.NotEmpty(t, file)
	})

	t.Run("Multiple Resources", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Resources: []appdef.Resource{
				{
					Name:     "db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
					Backup: appdef.ResourceBackupConfig{
						Enabled: true,
					},
				},
				{
					Name:     "store",
					Type:     appdef.ResourceTypeS3,
					Provider: appdef.ResourceProviderDigitalOcean,
					Backup: appdef.ResourceBackupConfig{
						Enabled: true,
					},
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		got := BackupWorkflow(t.Context(), input)
		assert.NoError(t, got)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "backup.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Resources: []appdef.Resource{
				{
					Name:     "db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
					Backup: appdef.ResourceBackupConfig{
						Enabled: true,
					},
				},
			},
		}

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), appDef)

		got := BackupWorkflow(t.Context(), input)
		assert.Error(t, got)
	})
}
