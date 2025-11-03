package cicd

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestPR(t *testing.T) {
	t.Parallel()

	t.Run("Creates Workflow", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:  "cms",
					Title: "CMS",
					Path:  "./cms",
					Type:  appdef.AppTypePayload,
				},
				{
					Name:  "web",
					Title: "Web",
					Path:  "./web",
					Type:  appdef.AppTypeGoLang,
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := PR(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "pr.yaml"))
		require.NoError(t, err)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		content := string(file)

		t.Log("Drift")
		{
			assert.Contains(t, content, "drift-detection:")
		}

		t.Log("Apps")
		{
			for _, app := range appDef.Apps {
				jobName := strings.ToLower(app.Name)
				assert.Contains(t, content, jobName+":", "workflow should contain job for app "+app.Name)

				switch app.Type {
				case appdef.AppTypeGoLang:
					assert.Contains(t, content, "Set up Go", "Go app should have Go setup")
				case appdef.AppTypePayload:
					assert.Contains(t, content, "Install pnpm", "JS app should have pnpm setup")
					assert.Contains(t, content, "Set up Node", "JS app should have Node setup")
				}
			}
		}

		t.Log("Commands")
		{
			// Get positions for each command in the canonical order.
			var positions []int
			for _, cmd := range appdef.Commands {
				pos := strings.Index(content, "name: "+strings.Title(cmd.String()))
				if pos != -1 {
					positions = append(positions, pos)
				}
			}

			// Verify positions are in ascending order.
			for i := 0; i < len(positions)-1; i++ {
				assert.Less(t, positions[i], positions[i+1],
					"commands should appear in order defined by appdef.Commands")
			}
		}
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:  "web",
					Title: "Web",
					Path:  "./web",
					Type:  appdef.AppTypeGoLang,
				},
			},
		}

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), appDef)

		err := PR(t.Context(), input)
		assert.Error(t, err)
	})

	t.Run("Migration Check For Payload With Postgres", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:  "cms",
					Title: "CMS",
					Path:  "./cms",
					Type:  appdef.AppTypePayload,
					Env: appdef.Environment{
						Production: appdef.EnvVar{
							"DATABASE_URL": {
								Source: appdef.EnvSourceResource,
								Value:  "db.connection_url",
							},
						},
					},
				},
			},
			Resources: []appdef.Resource{
				{
					Name: "db",
					Type: appdef.ResourceTypePostgres,
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := PR(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "pr.yaml"))
		require.NoError(t, err)

		content := string(file)

		err = validateGithubYaml(t, file, false)
		assert.NoError(t, err)

		t.Log("Migration Check")
		{
			assert.Contains(t, content, "migration-check-cms:", "should contain migration check job")
			assert.Contains(t, content, "Migration Check - CMS", "should contain migration check job name")
			assert.Contains(t, content, "Check for pending migrations", "should contain migration check step")
			assert.Contains(t, content, "pnpm migrate:create", "should run migration check command")
			assert.Contains(t, content, "db-add-ip", "should add runner IP to database")
			assert.Contains(t, content, "db-remove-ip", "should remove runner IP from database")
		}
	})

	t.Run("No Migration Check For Payload Without Postgres", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:  "cms",
					Title: "CMS",
					Path:  "./cms",
					Type:  appdef.AppTypePayload,
					Env: appdef.Environment{
						Production: appdef.EnvVar{
							"API_KEY": {
								Source: appdef.EnvSourceValue,
								Value:  "test-key",
							},
						},
					},
				},
			},
			Resources: []appdef.Resource{
				{
					Name: "storage",
					Type: appdef.ResourceTypeS3,
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := PR(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "pr.yaml"))
		require.NoError(t, err)

		content := string(file)

		assert.NotContains(t, content, "migration-check-cms:", "should not contain migration check job")
		assert.NotContains(t, content, "Check for pending migrations", "should not contain migration check step")
	})

	t.Run("No Migration Check For Non-Payload Apps", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name:  "api",
					Title: "API",
					Path:  "./api",
					Type:  appdef.AppTypeGoLang,
					Env: appdef.Environment{
						Production: appdef.EnvVar{
							"DATABASE_URL": {
								Source: appdef.EnvSourceResource,
								Value:  "db.connection_url",
							},
						},
					},
				},
			},
			Resources: []appdef.Resource{
				{
					Name: "db",
					Type: appdef.ResourceTypePostgres,
				},
			},
		}

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := PR(t.Context(), input)
		require.NoError(t, err)

		file, err := afero.ReadFile(input.FS, filepath.Join(workflowsPath, "pr.yaml"))
		require.NoError(t, err)

		content := string(file)

		assert.NotContains(t, content, "migration-check-api:", "should not contain migration check for non-Payload app")
		assert.NotContains(t, content, "Check for pending migrations", "should not contain migration check step")
	})
}
