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
					Name:       "db",
					Type:       appdef.ResourceTypePostgres,
					Provider:   appdef.ResourceProviderDigitalOcean,
					Monitoring: ptr.BoolPtr(true),
					Backup: appdef.ResourceBackupConfig{
						Enabled: ptr.BoolPtr(true),
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

		// Verify correct PROD prefix in Peekaping ping URLs (with vars. prefix for repository variables).
		content := string(file)
		assert.Contains(t, content, "${{ vars.PROD_DB_BACKUP_PING_URL }}")
		assert.Contains(t, content, "${{ vars.PROD_CODEBASE_BACKUP_PING_URL }}")
		assert.NotContains(t, content, "${{ vars._DB_BACKUP_PING_URL }}")
		assert.NotContains(t, content, "${{ vars._CODEBASE_BACKUP_PING_URL }}")
	})

	t.Run("S3", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Resources: []appdef.Resource{
				{
					Name:       "store",
					Type:       appdef.ResourceTypeS3,
					Provider:   appdef.ResourceProviderDigitalOcean,
					Monitoring: ptr.BoolPtr(true),
					Config: map[string]any{
						"key": "value",
					},
					Backup: appdef.ResourceBackupConfig{
						Enabled: ptr.BoolPtr(true),
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

		// Verify correct PROD prefix in Peekaping ping URLs (with vars. prefix for repository variables).
		content := string(file)
		assert.Contains(t, content, "${{ vars.PROD_STORE_BACKUP_PING_URL }}")
		assert.Contains(t, content, "${{ vars.PROD_CODEBASE_BACKUP_PING_URL }}")
		assert.NotContains(t, content, "${{ vars._STORE_BACKUP_PING_URL }}")
		assert.NotContains(t, content, "${{ vars._CODEBASE_BACKUP_PING_URL }}")
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
						Enabled: ptr.BoolPtr(true),
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
					Name:       "db",
					Type:       appdef.ResourceTypePostgres,
					Provider:   appdef.ResourceProviderDigitalOcean,
					Monitoring: ptr.BoolPtr(true),
					Backup: appdef.ResourceBackupConfig{
						Enabled: ptr.BoolPtr(true),
					},
				},
				{
					Name:       "store",
					Type:       appdef.ResourceTypeS3,
					Provider:   appdef.ResourceProviderDigitalOcean,
					Monitoring: ptr.BoolPtr(true),
					Backup: appdef.ResourceBackupConfig{
						Enabled: ptr.BoolPtr(true),
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

		// Verify correct PROD prefix in Peekaping ping URLs for all resources (with vars. prefix for repository variables).
		content := string(file)
		assert.Contains(t, content, "${{ vars.PROD_DB_BACKUP_PING_URL }}")
		assert.Contains(t, content, "${{ vars.PROD_STORE_BACKUP_PING_URL }}")
		assert.Contains(t, content, "${{ vars.PROD_CODEBASE_BACKUP_PING_URL }}")
		assert.NotContains(t, content, "${{ vars._DB_BACKUP_PING_URL }}")
		assert.NotContains(t, content, "${{ vars._STORE_BACKUP_PING_URL }}")
		assert.NotContains(t, content, "${{ vars._CODEBASE_BACKUP_PING_URL }}")
	})

	t.Run("SQLite Turso", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Resources: []appdef.Resource{
				{
					Name:       "db",
					Type:       appdef.ResourceTypeSQLite,
					Provider:   appdef.ResourceProviderTurso,
					Monitoring: ptr.BoolPtr(true),
					Backup: appdef.ResourceBackupConfig{
						Enabled: ptr.BoolPtr(true),
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

		// Verify correct PROD prefix in Peekaping ping URLs (with vars. prefix for repository variables).
		content := string(file)
		assert.Contains(t, content, "${{ vars.PROD_DB_BACKUP_PING_URL }}")
		assert.Contains(t, content, "${{ vars.PROD_CODEBASE_BACKUP_PING_URL }}")
		assert.NotContains(t, content, "${{ vars._DB_BACKUP_PING_URL }}")
		assert.NotContains(t, content, "${{ vars._CODEBASE_BACKUP_PING_URL }}")
	})

	t.Run("Multiple Resources Including SQLite", func(t *testing.T) {
		t.Parallel()

		// This verifies that mixing SQLite with other resource types
		// generates valid YAML (i.e., sync-to-gdrive depends on all backup jobs).
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Resources: []appdef.Resource{
				{
					Name:       "db",
					Type:       appdef.ResourceTypeSQLite,
					Provider:   appdef.ResourceProviderTurso,
					Monitoring: ptr.BoolPtr(true),
					Backup: appdef.ResourceBackupConfig{
						Enabled: ptr.BoolPtr(true),
					},
				},
				{
					Name:       "store",
					Type:       appdef.ResourceTypeS3,
					Provider:   appdef.ResourceProviderDigitalOcean,
					Monitoring: ptr.BoolPtr(true),
					Backup: appdef.ResourceBackupConfig{
						Enabled: ptr.BoolPtr(true),
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

		// Verify correct PROD prefix in Peekaping ping URLs for all resources (with vars. prefix for repository variables).
		content := string(file)
		assert.Contains(t, content, "${{ vars.PROD_DB_BACKUP_PING_URL }}")
		assert.Contains(t, content, "${{ vars.PROD_STORE_BACKUP_PING_URL }}")
		assert.Contains(t, content, "${{ vars.PROD_CODEBASE_BACKUP_PING_URL }}")
		assert.NotContains(t, content, "${{ vars._DB_BACKUP_PING_URL }}")
		assert.NotContains(t, content, "${{ vars._STORE_BACKUP_PING_URL }}")
		assert.NotContains(t, content, "${{ vars._CODEBASE_BACKUP_PING_URL }}")
	})

	t.Run("Global monitoring disabled", func(t *testing.T) {
		t.Parallel()

		enabled := false
		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Monitoring: appdef.Monitoring{
				Enabled: &enabled,
			},
			Resources: []appdef.Resource{
				{
					Name:       "db",
					Type:       appdef.ResourceTypePostgres,
					Provider:   appdef.ResourceProviderDigitalOcean,
					Monitoring: ptr.BoolPtr(true),
					Backup: appdef.ResourceBackupConfig{
						Enabled: ptr.BoolPtr(true),
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

		content := string(file)
		assert.NotContains(t, content, "Ping Peekaping Heartbeat Monitor")
		assert.NotContains(t, content, "PROD_DB_BACKUP_PING_URL")
		assert.NotContains(t, content, "PROD_CODEBASE_BACKUP_PING_URL")
	})

	t.Run("Per-resource monitoring disabled", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Resources: []appdef.Resource{
				{
					Name:       "db",
					Type:       appdef.ResourceTypePostgres,
					Provider:   appdef.ResourceProviderDigitalOcean,
					Monitoring: ptr.BoolPtr(false),
					Backup: appdef.ResourceBackupConfig{
						Enabled: ptr.BoolPtr(true),
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

		content := string(file)
		assert.NotContains(t, content, "PROD_DB_BACKUP_PING_URL")
		assert.Contains(t, content, "PROD_CODEBASE_BACKUP_PING_URL")
	})

	t.Run("Mixed monitoring enabled and disabled resources", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Resources: []appdef.Resource{
				{
					Name:       "db",
					Type:       appdef.ResourceTypePostgres,
					Provider:   appdef.ResourceProviderDigitalOcean,
					Monitoring: ptr.BoolPtr(true),
					Backup: appdef.ResourceBackupConfig{
						Enabled: ptr.BoolPtr(true),
					},
				},
				{
					Name:       "store",
					Type:       appdef.ResourceTypeS3,
					Provider:   appdef.ResourceProviderDigitalOcean,
					Monitoring: ptr.BoolPtr(false),
					Backup: appdef.ResourceBackupConfig{
						Enabled: ptr.BoolPtr(true),
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

		content := string(file)
		assert.Contains(t, content, "PROD_DB_BACKUP_PING_URL")
		assert.NotContains(t, content, "PROD_STORE_BACKUP_PING_URL")
		assert.Contains(t, content, "PROD_CODEBASE_BACKUP_PING_URL")
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
			},
			Resources: []appdef.Resource{
				{
					Name:       "db",
					Type:       appdef.ResourceTypePostgres,
					Provider:   appdef.ResourceProviderDigitalOcean,
					Monitoring: ptr.BoolPtr(true),
					Backup: appdef.ResourceBackupConfig{
						Enabled: ptr.BoolPtr(true),
					},
				},
			},
		}

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), appDef)

		got := BackupWorkflow(t.Context(), input)
		assert.Error(t, got)
	})
}
