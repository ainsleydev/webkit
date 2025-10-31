//go:build !race

package infra

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestTFVarsFromDefinition(t *testing.T) {
	t.Run("Nil Definition", func(t *testing.T) {
		_, err := tfVarsFromDefinition(env.Development, nil)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "definition cannot be nil")
	})

	t.Run("Empty Definition", func(t *testing.T) {
		input := &appdef.Definition{
			Project:   appdef.Project{Name: "project"},
			Apps:      []appdef.App{},
			Resources: []appdef.Resource{},
		}

		got, err := tfVarsFromDefinition(env.Production, input)
		assert.NoError(t, err)

		t.Log("Metadata")
		{
			assert.Equal(t, "project", got.ProjectName)
			assert.Equal(t, "production", got.Environment)
		}
	})

	t.Run("Single Resource", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name: "project",
				Repo: appdef.GitHubRepo{
					Owner: "ainsley-dev",
					Name:  "project",
				},
			},
			Resources: []appdef.Resource{
				{
					Name:     "db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"key": "value",
					},
					Backup: appdef.ResourceBackupConfig{},
				},
			},
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Production: map[string]appdef.EnvValue{
						"VALUE_KEY":  {Value: "bar", Source: appdef.EnvSourceValue},
						"SECRET_KEY": {Value: "s3cr3t", Source: appdef.EnvSourceSOPS},
					},
				},
			},
		}

		got, err := tfVarsFromDefinition(env.Production, input)
		assert.NoError(t, err)

		t.Log("Metadata")
		{
			assert.Equal(t, "project", got.ProjectName)
			assert.Equal(t, "production", got.Environment)
			assert.Equal(t, "ainsley-dev", got.GithubConfig.Owner)
			assert.Equal(t, "project", got.GithubConfig.Repo)
		}

		t.Log("Resource")
		{
			require.Len(t, got.Resources, 1)

			resource := got.Resources[0]
			assert.Equal(t, resource.Name, "db")
			assert.Equal(t, resource.PlatformType, appdef.ResourceTypePostgres.String())
			assert.Equal(t, resource.PlatformProvider, appdef.ResourceProviderDigitalOcean.String())
			assert.Equal(t, resource.Config, map[string]any{"key": "value"})
		}
	})

	t.Run("Single App", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name: "single-app-project",
				Repo: appdef.GitHubRepo{
					Owner: "owner",
					Name:  "single-app-project",
				},
			},
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "apps/cms",
					Infra: appdef.Infra{
						Type:     "docker",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config: map[string]any{
							"replicas": 2,
						},
					},
					Env: appdef.Environment{
						Production: map[string]appdef.EnvValue{
							"VALUE_KEY": {Value: "nested", Source: appdef.EnvSourceValue},
						},
					},
				},
			},
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Production: map[string]appdef.EnvValue{
						"VALUE_KEY":  {Value: "parent", Source: appdef.EnvSourceValue},
						"SECRET_KEY": {Value: "s3cr3t", Source: appdef.EnvSourceSOPS},
					},
				},
			},
		}

		got, err := tfVarsFromDefinition(env.Production, input)
		assert.NoError(t, err)

		t.Log("Metadata")
		{
			assert.Equal(t, "single-app-project", got.ProjectName)
			assert.Equal(t, "production", got.Environment)
			assert.Equal(t, "owner", got.GithubConfig.Owner)
			assert.Equal(t, "single-app-project", got.GithubConfig.Repo)
		}

		t.Log("Apps")
		{
			require.Len(t, got.Apps, 1)

			app := got.Apps[0]
			assert.Equal(t, "cms", app.Name)
			assert.Equal(t, "docker", app.PlatformType)
			assert.Equal(t, appdef.ResourceProviderDigitalOcean.String(), app.PlatformProvider)
			assert.Equal(t, map[string]any{"replicas": 2}, app.Config)

			require.Len(t, app.Environment, 2)
			assert.ElementsMatch(t, app.Environment, []tfEnvVar{
				{Key: "VALUE_KEY", Value: "nested", Source: "value", Scope: "GENERAL"},
				{Key: "SECRET_KEY", Value: "s3cr3t", Source: "sops", Scope: "SECRET"},
			})
		}
	})

	t.Run("Multiple Apps and Resources", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name: "complex-project",
				Repo: appdef.GitHubRepo{
					Owner: "owner",
					Name:  "complex-project",
				},
			},
			Apps: []appdef.App{
				{
					Name: "frontend",
					Type: appdef.AppTypePayload,
					Path: "apps/frontend",
					Infra: appdef.Infra{
						Type:     "docker",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config: map[string]any{
							"replicas": 2,
						},
					},
					Env: appdef.Environment{
						Production: map[string]appdef.EnvValue{
							"FRONTEND_KEY": {Value: "frontend_value", Source: appdef.EnvSourceValue},
						},
					},
				},
				{
					Name: "backend",
					Type: appdef.AppTypePayload,
					Path: "apps/backend",
					Infra: appdef.Infra{
						Type:     "docker",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config: map[string]any{
							"replicas": 3,
						},
					},
					Env: appdef.Environment{
						Production: map[string]appdef.EnvValue{
							"BACKEND_KEY": {Value: "backend_value", Source: appdef.EnvSourceValue},
						},
						// Prove that it doesn't get appended.
						Dev: map[string]appdef.EnvValue{
							"DEV_KEY": {Value: "dev", Source: appdef.EnvSourceValue},
						},
					},
				},
			},
			Resources: []appdef.Resource{
				{
					Name:     "db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"version": "18",
					},
				},
				{
					Name:     "storage",
					Type:     appdef.ResourceTypeS3,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"acl": "private",
					},
				},
			},
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Production: map[string]appdef.EnvValue{
						"COMMON_KEY": {Value: "common_value", Source: appdef.EnvSourceValue},
						"SECRET_KEY": {Value: "s3cr3t", Source: appdef.EnvSourceSOPS},
					},
				},
			},
		}

		got, err := tfVarsFromDefinition(env.Production, input)
		assert.NoError(t, err)

		t.Log("Metadata")
		{
			assert.Equal(t, "complex-project", got.ProjectName)
			assert.Equal(t, "production", got.Environment)
			assert.Equal(t, "owner", got.GithubConfig.Owner)
			assert.Equal(t, "complex-project", got.GithubConfig.Repo)
		}

		t.Log("Apps")
		{
			require.Len(t, got.Apps, 2)

			frontend := got.Apps[0]
			assert.Equal(t, "frontend", frontend.Name)
			assert.Equal(t, "docker", frontend.PlatformType)
			assert.Equal(t, appdef.ResourceProviderDigitalOcean.String(), frontend.PlatformProvider)
			assert.Equal(t, map[string]any{"replicas": 2}, frontend.Config)
			assert.ElementsMatch(t, frontend.Environment, []tfEnvVar{
				{Key: "FRONTEND_KEY", Value: "frontend_value", Source: "value", Scope: "GENERAL"},
				{Key: "COMMON_KEY", Value: "common_value", Source: "value", Scope: "GENERAL"},
				{Key: "SECRET_KEY", Value: "s3cr3t", Source: "sops", Scope: "SECRET"},
			})

			backend := got.Apps[1]
			assert.Equal(t, "backend", backend.Name)
			assert.Equal(t, "docker", backend.PlatformType)
			assert.Equal(t, appdef.ResourceProviderDigitalOcean.String(), backend.PlatformProvider)
			assert.Equal(t, map[string]any{"replicas": 3}, backend.Config)
			assert.ElementsMatch(t, backend.Environment, []tfEnvVar{
				{Key: "BACKEND_KEY", Value: "backend_value", Source: "value", Scope: "GENERAL"},
				{Key: "COMMON_KEY", Value: "common_value", Source: "value", Scope: "GENERAL"},
				{Key: "SECRET_KEY", Value: "s3cr3t", Source: "sops", Scope: "SECRET"},
			})
		}

		t.Log("Resources")
		{
			require.Len(t, got.Resources, 2)

			db := got.Resources[0]
			assert.Equal(t, "db", db.Name)
			assert.Equal(t, appdef.ResourceTypePostgres.String(), db.PlatformType)
			assert.Equal(t, appdef.ResourceProviderDigitalOcean.String(), db.PlatformProvider)
			assert.Equal(t, map[string]any{"version": "18"}, db.Config)

			cache := got.Resources[1]
			assert.Equal(t, "storage", cache.Name)
			assert.Equal(t, appdef.ResourceTypeS3.String(), cache.PlatformType)
			assert.Equal(t, appdef.ResourceProviderDigitalOcean.String(), cache.PlatformProvider)
			assert.Equal(t, map[string]any{"acl": "private"}, cache.Config)
		}
	})

	t.Run("Mixed null and empty configs with arrays", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name: "mixed-config-project",
				Repo: appdef.GitHubRepo{
					Owner: "owner",
					Name:  "mixed-config-project",
				},
			},
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "cms",
					Infra: appdef.Infra{
						Type:     "vm",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config:   nil, // Null config
					},
				},
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "web",
					Infra: appdef.Infra{
						Type:     "app",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config:   map[string]any{}, // Empty config
					},
				},
			},
			Resources: []appdef.Resource{
				{
					Name:     "db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"allowed_ips_addr": []any{"185.16.161.205", "159.65.87.97"},
						"engine_version":   "17",
					},
				},
				{
					Name:     "store",
					Type:     appdef.ResourceTypeS3,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"acl": "private",
					},
				},
			},
		}

		got, err := tfVarsFromDefinition(env.Production, input)
		assert.NoError(t, err)

		t.Log("Apps with consistent config types")
		{
			require.Len(t, got.Apps, 2)

			cms := got.Apps[0]
			assert.Equal(t, "cms", cms.Name)
			assert.Equal(t, map[string]any{}, cms.Config) // Nil should become {}

			web := got.Apps[1]
			assert.Equal(t, "web", web.Name)
			assert.Equal(t, map[string]any{}, web.Config) // Empty should stay {}
		}

		t.Log("Resources with normalized arrays")
		{
			require.Len(t, got.Resources, 2)

			db := got.Resources[0]
			assert.Equal(t, "db", db.Name)
			// Arrays should be properly typed
			assert.Equal(t, []string{"185.16.161.205", "159.65.87.97"}, db.Config["allowed_ips_addr"])
			assert.Equal(t, "17", db.Config["engine_version"])

			store := got.Resources[1]
			assert.Equal(t, "store", store.Name)
			assert.Equal(t, "private", store.Config["acl"])
		}
	})
}
