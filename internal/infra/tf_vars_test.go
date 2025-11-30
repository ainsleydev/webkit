//go:build !race

package infra

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/mocks"
	"github.com/ainsleydev/webkit/pkg/env"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func setupTfVars(t *testing.T, appDef *appdef.Definition) *Terraform {
	t.Helper()

	ctrl := gomock.NewController(t)
	mockClient := mocks.NewGHClient(ctrl)

	// Default behavior: return empty string (no SHA tags).
	mockClient.EXPECT().
		GetLatestSHATag(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return("", nil).
		AnyTimes()

	return &Terraform{
		appDef:   appDef,
		fs:       afero.NewMemMapFs(),
		ghClient: mockClient,
	}
}

func TestTFVarsFromDefinition(t *testing.T) {
	t.Run("Nil Definition", func(t *testing.T) {
		tf := setupTfVars(t, nil)
		_, err := tf.tfVarsFromDefinition(context.Background(), env.Development)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "definition cannot be nil")
	})

	t.Run("Empty Definition", func(t *testing.T) {
		input := &appdef.Definition{
			Project:   appdef.Project{Name: "project"},
			Apps:      []appdef.App{},
			Resources: []appdef.Resource{},
		}

		tf := setupTfVars(t, input)
		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
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
					Title:    "Database",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"key": "value",
					},
					Backup: &appdef.ResourceBackupConfig{},
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

		tf := setupTfVars(t, input)
		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
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

		tf := setupTfVars(t, input)
		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
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

		tf := setupTfVars(t, input)
		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
		assert.NoError(t, err)

		t.Log("Metadata")
		{
			assert.Equal(t, "complex-project", got.ProjectName)
			assert.Equal(t, "production", got.Environment)
			assert.Equal(t, "owner", got.GithubConfig.Owner)
			assert.Equal(t, "complex-project", got.GithubConfig.Repo)
		}

		t.Log("SSH")
		{
			assert.ElementsMatch(t, []string{"Ainsley - Mac Studio"}, got.DigitalOceanSSHKeys)
			assert.ElementsMatch(t, []string{"hello@ainsley.dev"}, got.HetznerSSHKeys)
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

	t.Run("App with Domains", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name: "domain-project",
				Repo: appdef.GitHubRepo{
					Owner: "owner",
					Name:  "domain-project",
				},
			},
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "apps/web",
					Infra: appdef.Infra{
						Type:     "app",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config: map[string]any{
							"replicas": 2,
						},
					},
					Domains: []appdef.Domain{
						{
							Name:     "example.com",
							Type:     appdef.DomainTypePrimary,
							Zone:     "example.com",
							Wildcard: false,
						},
						{
							Name:     "www.example.com",
							Type:     appdef.DomainTypeAlias,
							Zone:     "example.com",
							Wildcard: false,
						},
						{
							Name:     "*.staging.example.com",
							Type:     appdef.DomainTypePrimary,
							Zone:     "staging.example.com",
							Wildcard: true,
						},
					},
				},
			},
		}

		tf := setupTfVars(t, input)
		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
		assert.NoError(t, err)

		t.Log("App with domains")
		{
			require.Len(t, got.Apps, 1)

			app := got.Apps[0]
			assert.Equal(t, "web", app.Name)

			require.Len(t, app.Domains, 3)
			assert.Equal(t, "example.com", app.Domains[0].Name)
			assert.Equal(t, "PRIMARY", app.Domains[0].Type)
			assert.Equal(t, "example.com", app.Domains[0].Zone)
			assert.Equal(t, false, app.Domains[0].Wildcard)

			assert.Equal(t, "www.example.com", app.Domains[1].Name)
			assert.Equal(t, "ALIAS", app.Domains[1].Type)
			assert.Equal(t, "example.com", app.Domains[1].Zone)
			assert.Equal(t, false, app.Domains[1].Wildcard)

			assert.Equal(t, "*.staging.example.com", app.Domains[2].Name)
			assert.Equal(t, "PRIMARY", app.Domains[2].Type)
			assert.Equal(t, "staging.example.com", app.Domains[2].Zone)
			assert.Equal(t, true, app.Domains[2].Wildcard)
		}
	})

	t.Run("Status page domain - no apps", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name: "no-apps-project",
				Repo: appdef.GitHubRepo{
					Owner: "owner",
					Name:  "no-apps-project",
				},
			},
			Apps: []appdef.App{},
		}

		tf := setupTfVars(t, input)
		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
		assert.NoError(t, err)

		t.Log("Status page domain should be nil when no apps")
		{
			assert.Nil(t, got.StatusPageDomain)
		}
	})

	t.Run("Status page domain - app with no domains", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name: "no-domains-project",
				Repo: appdef.GitHubRepo{
					Owner: "owner",
					Name:  "no-domains-project",
				},
			},
			Apps: []appdef.App{
				{
					Name: "api",
					Type: appdef.AppTypeGoLang,
					Path: "apps/api",
					Infra: appdef.Infra{
						Type:     "vm",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config:   map[string]any{},
					},
					Domains: []appdef.Domain{},
				},
			},
		}

		tf := setupTfVars(t, input)
		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
		assert.NoError(t, err)

		t.Log("Status page domain should be nil when app has no domains")
		{
			assert.Nil(t, got.StatusPageDomain)
		}
	})

	t.Run("Status page domain - explicit configuration", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name: "explicit-domain-project",
				Repo: appdef.GitHubRepo{
					Owner: "owner",
					Name:  "explicit-domain-project",
				},
			},
			Monitoring: appdef.Monitoring{
				StatusPage: appdef.StatusPage{
					Domain: "status.custom.com",
				},
			},
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "apps/cms",
					Infra: appdef.Infra{
						Type:     "app",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config:   map[string]any{},
					},
					Domains: []appdef.Domain{
						{Name: "cms.player2clubs.com", Type: appdef.DomainTypePrimary},
					},
				},
			},
		}

		tf := setupTfVars(t, input)
		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
		assert.NoError(t, err)

		t.Log("Status page domain should use explicit configuration")
		{
			require.NotNil(t, got.StatusPageDomain)
			assert.Equal(t, "status.custom.com", *got.StatusPageDomain)
		}
	})

	t.Run("Status page slug - not set", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name: "no-slug-project",
				Repo: appdef.GitHubRepo{
					Owner: "owner",
					Name:  "no-slug-project",
				},
			},
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "apps/web",
					Infra: appdef.Infra{
						Type:     "app",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config:   map[string]any{},
					},
				},
			},
		}

		tf := setupTfVars(t, input)
		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
		assert.NoError(t, err)

		t.Log("Status page slug should be nil when not set")
		{
			assert.Nil(t, got.StatusPageSlug)
		}
	})

	t.Run("Status page slug - explicit configuration", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name: "explicit-slug-project",
				Repo: appdef.GitHubRepo{
					Owner: "owner",
					Name:  "explicit-slug-project",
				},
			},
			Monitoring: appdef.Monitoring{
				StatusPage: appdef.StatusPage{
					Domain: "status.custom.com",
					Slug:   "my-custom-slug",
				},
			},
			Apps: []appdef.App{
				{
					Name: "cms",
					Type: appdef.AppTypePayload,
					Path: "apps/cms",
					Infra: appdef.Infra{
						Type:     "app",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config:   map[string]any{},
					},
					Domains: []appdef.Domain{
						{Name: "cms.example.com", Type: appdef.DomainTypePrimary},
					},
				},
			},
		}

		tf := setupTfVars(t, input)
		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
		assert.NoError(t, err)

		t.Log("Status page slug should use explicit configuration")
		{
			require.NotNil(t, got.StatusPageSlug)
			assert.Equal(t, "my-custom-slug", *got.StatusPageSlug)
		}

		t.Log("Status page domain should also be set")
		{
			require.NotNil(t, got.StatusPageDomain)
			assert.Equal(t, "status.custom.com", *got.StatusPageDomain)
		}
	})

	t.Run("Status page theme - not set", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name: "no-theme-project",
				Repo: appdef.GitHubRepo{
					Owner: "owner",
					Name:  "no-theme-project",
				},
			},
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "apps/web",
					Infra: appdef.Infra{
						Type:     "app",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config:   map[string]any{},
					},
				},
			},
		}

		tf := setupTfVars(t, input)
		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
		assert.NoError(t, err)

		t.Log("Status page theme should be nil when not set")
		{
			assert.Nil(t, got.StatusPageTheme)
		}
	})

	t.Run("Status page theme - explicit configuration", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name: "dark-theme-project",
				Repo: appdef.GitHubRepo{
					Owner: "owner",
					Name:  "dark-theme-project",
				},
			},
			Monitoring: appdef.Monitoring{
				StatusPage: appdef.StatusPage{
					Theme: "dark",
				},
			},
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "apps/web",
					Infra: appdef.Infra{
						Type:     "app",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config:   map[string]any{},
					},
				},
			},
		}

		tf := setupTfVars(t, input)
		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
		assert.NoError(t, err)

		t.Log("Status page theme should use explicit configuration")
		{
			require.NotNil(t, got.StatusPageTheme)
			assert.Equal(t, "dark", *got.StatusPageTheme)
		}
	})

	t.Run("Branding", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name:  "brand-project",
				Title: "Brand Project",
				Repo: appdef.GitHubRepo{
					Owner: "owner",
					Name:  "brand-project",
				},
				Brand: appdef.Brand{
					IconURL:       "https://example.com/favicon.ico",
					PrimaryColour: "#FF5733",
				},
			},
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Path: "apps/web",
					Infra: appdef.Infra{
						Type:     "app",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config:   map[string]any{},
					},
				},
			},
		}

		tf := setupTfVars(t, input)
		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
		assert.NoError(t, err)

		t.Log("Brand fields should be populated when set")
		{
			require.NotNil(t, got.BrandIconURL)
			assert.Equal(t, "https://example.com/favicon.ico", *got.BrandIconURL)
			require.NotNil(t, got.BrandPrimaryColor)
			assert.Equal(t, "#FF5733", *got.BrandPrimaryColor)
		}

		t.Log("Brand fields should be nil when empty")
		{
			input.Project.Brand = appdef.Brand{}
			tf = setupTfVars(t, input)
			got, err = tf.tfVarsFromDefinition(context.Background(), env.Production)
			assert.NoError(t, err)
			assert.Nil(t, got.BrandIconURL)
			assert.Nil(t, got.BrandPrimaryColor)
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

		tf := setupTfVars(t, input)
		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
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

		t.Log("Resources with JSON-encoded arrays")
		{
			require.Len(t, got.Resources, 2)

			db := got.Resources[0]
			assert.Equal(t, "db", db.Name)
			// Arrays should be JSON-encoded as strings for Terraform's jsondecode()
			assert.Equal(t, `["185.16.161.205","159.65.87.97"]`, db.Config["allowed_ips_addr"])
			assert.Equal(t, "17", db.Config["engine_version"])

			store := got.Resources[1]
			assert.Equal(t, "store", store.Name)
			assert.Equal(t, "private", store.Config["acl"])
		}
	})
}

func TestEncodeConfigForTerraform(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input map[string]any
		want  map[string]any
	}{
		"Nil config returns empty map": {
			input: nil,
			want:  map[string]any{},
		},
		"Empty config": {
			input: map[string]any{},
			want:  map[string]any{},
		},
		"String array gets JSON encoded": {
			input: map[string]any{
				"allowed_ips_addr": []any{"185.16.161.205", "159.65.87.97"},
				"engine_version":   "17",
			},
			want: map[string]any{
				"allowed_ips_addr": `["185.16.161.205","159.65.87.97"]`,
				"engine_version":   "17",
			},
		},
		"Typed string slice gets JSON encoded": {
			input: map[string]any{
				"allowed_droplet_ips": []string{"192.168.1.1", "192.168.1.2"},
				"size":                "db-s-1vcpu-1gb",
			},
			want: map[string]any{
				"allowed_droplet_ips": `["192.168.1.1","192.168.1.2"]`,
				"size":                "db-s-1vcpu-1gb",
			},
		},
		"Primitives pass through unchanged": {
			input: map[string]any{
				"string": "value",
				"number": 42,
				"bool":   true,
				"null":   nil,
			},
			want: map[string]any{
				"string": "value",
				"number": 42,
				"bool":   true,
				"null":   nil,
			},
		},
		"Mixed primitives and arrays": {
			input: map[string]any{
				"acl":   "private",
				"ports": []int{8080, 443, 3000},
			},
			want: map[string]any{
				"acl":   "private",
				"ports": `[8080,443,3000]`,
			},
		},
		"Real-world Postgres config": {
			input: map[string]any{
				"allowed_ips_addr":    []string{"185.16.161.205", "159.65.87.97"},
				"allowed_droplet_ips": []string{"droplet-123"},
				"engine_version":      "17",
				"size":                "db-s-1vcpu-1gb",
			},
			want: map[string]any{
				"allowed_ips_addr":    `["185.16.161.205","159.65.87.97"]`,
				"allowed_droplet_ips": `["droplet-123"]`,
				"engine_version":      "17",
				"size":                "db-s-1vcpu-1gb",
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := encodeConfigForTerraform(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestEncodeConfigValue(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input any
		want  any
	}{
		"Nil value": {
			input: nil,
			want:  nil,
		},
		"String value": {
			input: "test",
			want:  "test",
		},
		"Number value": {
			input: 42,
			want:  42,
		},
		"Bool value": {
			input: true,
			want:  true,
		},
		"Interface slice with strings": {
			input: []any{"a", "b", "c"},
			want:  `["a","b","c"]`,
		},
		"Typed string slice": {
			input: []string{"x", "y", "z"},
			want:  `["x","y","z"]`,
		},
		"Int slice": {
			input: []int{1, 2, 3},
			want:  `[1,2,3]`,
		},
		"Float slice": {
			input: []float64{1.5, 2.5},
			want:  `[1.5,2.5]`,
		},
		"Bool slice": {
			input: []bool{true, false},
			want:  `[true,false]`,
		},
		"Empty slice": {
			input: []string{},
			want:  `[]`,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := encodeConfigValue(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestTerraform_TFVarsFromDefinition_ImageTag(t *testing.T) {
	t.Run("Container app gets image tag from client", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
				Repo: appdef.GitHubRepo{
					Owner: "test-owner",
					Name:  "test-repo",
				},
			},
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Infra: appdef.Infra{
						Type:     "container",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config:   map[string]any{},
					},
				},
			},
		}

		// Create Terraform with mock client that returns a specific tag.
		ctrl := gomock.NewController(t)
		mockClient := mocks.NewGHClient(ctrl)
		mockClient.EXPECT().
			GetLatestSHATag(gomock.Any(), "test-owner", "test-repo", "web").
			Return("sha-test123", nil)

		tf := &Terraform{
			appDef:   input,
			fs:       afero.NewMemMapFs(),
			ghClient: mockClient,
		}

		// Ensure GITHUB_SHA is not set.
		t.Setenv("GITHUB_SHA", "")

		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
		require.NoError(t, err)

		require.Len(t, got.Apps, 1)
		assert.Equal(t, "sha-test123", got.Apps[0].ImageTag)
	})

	t.Run("Non-container app does not get image tag", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
				Repo: appdef.GitHubRepo{
					Owner: "test-owner",
					Name:  "test-repo",
				},
			},
			Apps: []appdef.App{
				{
					Name: "api",
					Type: appdef.AppTypeGoLang,
					Infra: appdef.Infra{
						Type:     "vm",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config:   map[string]any{},
					},
				},
			},
		}

		tf := setupTfVars(t, input)
		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
		require.NoError(t, err)

		require.Len(t, got.Apps, 1)
		assert.Equal(t, "", got.Apps[0].ImageTag)
	})

	t.Run("Uses GITHUB_SHA when set", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-project",
				Repo: appdef.GitHubRepo{
					Owner: "test-owner",
					Name:  "test-repo",
				},
			},
			Apps: []appdef.App{
				{
					Name: "web",
					Type: appdef.AppTypeSvelteKit,
					Infra: appdef.Infra{
						Type:     "container",
						Provider: appdef.ResourceProviderDigitalOcean,
						Config:   map[string]any{},
					},
				},
			},
		}

		tf := setupTfVars(t, input)

		// Set GITHUB_SHA environment variable.
		t.Setenv("GITHUB_SHA", "ci-sha-123")

		got, err := tf.tfVarsFromDefinition(context.Background(), env.Production)
		require.NoError(t, err)

		require.Len(t, got.Apps, 1)
		assert.Equal(t, "sha-ci-sha-123", got.Apps[0].ImageTag)
	})
}

func TestGenerateMonitors(t *testing.T) {
	t.Run("No Apps Or Resources", func(t *testing.T) {
		input := &appdef.Definition{
			Project:   appdef.Project{Name: "empty", Title: "Empty Project"},
			Apps:      []appdef.App{},
			Resources: []appdef.Resource{},
		}

		tf := setupTfVars(t, input)
		monitors := tf.generateMonitors(env.Production)
		// Codebase backup monitor is always generated.
		require.Len(t, monitors, 1)
		assert.Equal(t, "Backup - Codebase", monitors[0].Name)
		assert.Equal(t, "push", monitors[0].Type)
		assert.Equal(t, appdef.MonitorIntervalBackup, monitors[0].Interval)
	})

	t.Run("Single App With Monitoring Enabled", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{Name: "test", Title: "Test Project"},
			Apps: []appdef.App{
				{
					Name:  "web",
					Title: "Web",
					Domains: []appdef.Domain{
						{Name: "example.com", Type: appdef.DomainTypePrimary},
					},
					Infra: appdef.Infra{
						Config: map[string]any{"health_check_path": "/health"},
					},
					Monitoring: ptr.BoolPtr(true),
				},
			},
		}

		tf := setupTfVars(t, input)
		monitors := tf.generateMonitors(env.Production)
		require.Len(t, monitors, 3) // HTTP + DNS + Codebase Backup

		// HTTP monitor.
		assert.Equal(t, "HTTP - example.com", monitors[0].Name)
		assert.Equal(t, "http", monitors[0].Type)
		assert.Equal(t, "https://example.com", monitors[0].URL)
		assert.Equal(t, "GET", monitors[0].Method)
		assert.Equal(t, appdef.MonitorIntervalHTTP, monitors[0].Interval)

		// DNS monitor.
		assert.Equal(t, "DNS - example.com", monitors[1].Name)
		assert.Equal(t, "dns", monitors[1].Type)
		assert.Equal(t, "example.com", monitors[1].Domain)
		assert.Equal(t, appdef.MonitorIntervalDNS, monitors[1].Interval)

		// Codebase backup monitor.
		assert.Equal(t, "Backup - Codebase", monitors[2].Name)
		assert.Equal(t, "push", monitors[2].Type)
		assert.Equal(t, appdef.MonitorIntervalBackup, monitors[2].Interval)
	})

	t.Run("App With Monitoring Disabled", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{Name: "test", Title: "Test Project"},
			Apps: []appdef.App{
				{
					Name: "web",
					Domains: []appdef.Domain{
						{Name: "example.com", Type: appdef.DomainTypePrimary},
					},
					Monitoring: ptr.BoolPtr(false),
				},
			},
		}

		tf := setupTfVars(t, input)
		monitors := tf.generateMonitors(env.Production)
		// Even with app monitoring disabled, codebase backup monitor is still generated.
		require.Len(t, monitors, 1)
		assert.Equal(t, "Backup - Codebase", monitors[0].Name)
		assert.Equal(t, "push", monitors[0].Type)
		assert.Equal(t, appdef.MonitorIntervalBackup, monitors[0].Interval)
	})

	t.Run("Multiple Apps Multiple Domains", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{Name: "test", Title: "Test Project"},
			Apps: []appdef.App{
				{
					Name:  "web",
					Title: "Web",
					Domains: []appdef.Domain{
						{Name: "example.com", Type: appdef.DomainTypePrimary},
						{Name: "www.example.com", Type: appdef.DomainTypeAlias},
					},
					Infra:      appdef.Infra{},
					Monitoring: ptr.BoolPtr(true),
				},
				{
					Name:  "api",
					Title: "API",
					Domains: []appdef.Domain{
						{Name: "api.example.com", Type: appdef.DomainTypePrimary},
					},
					Infra:      appdef.Infra{},
					Monitoring: ptr.BoolPtr(true),
				},
			},
		}

		tf := setupTfVars(t, input)
		monitors := tf.generateMonitors(env.Production)
		require.Len(t, monitors, 7) // 3 domains × 2 types (HTTP + DNS) + Codebase Backup

		assert.Equal(t, "HTTP - example.com", monitors[0].Name)
		assert.Equal(t, "http", monitors[0].Type)
		assert.Equal(t, appdef.MonitorIntervalHTTP, monitors[0].Interval)
		assert.Equal(t, "DNS - example.com", monitors[1].Name)
		assert.Equal(t, "dns", monitors[1].Type)
		assert.Equal(t, appdef.MonitorIntervalDNS, monitors[1].Interval)
		assert.Equal(t, "HTTP - www.example.com", monitors[2].Name)
		assert.Equal(t, "http", monitors[2].Type)
		assert.Equal(t, appdef.MonitorIntervalHTTP, monitors[2].Interval)
		assert.Equal(t, "DNS - www.example.com", monitors[3].Name)
		assert.Equal(t, "dns", monitors[3].Type)
		assert.Equal(t, appdef.MonitorIntervalDNS, monitors[3].Interval)
		assert.Equal(t, "HTTP - api.example.com", monitors[4].Name)
		assert.Equal(t, "http", monitors[4].Type)
		assert.Equal(t, appdef.MonitorIntervalHTTP, monitors[4].Interval)
		assert.Equal(t, "DNS - api.example.com", monitors[5].Name)
		assert.Equal(t, "dns", monitors[5].Type)
		assert.Equal(t, appdef.MonitorIntervalDNS, monitors[5].Interval)
		assert.Equal(t, "Backup - Codebase", monitors[6].Name)
		assert.Equal(t, "push", monitors[6].Type)
		assert.Equal(t, appdef.MonitorIntervalBackup, monitors[6].Interval)
	})

	t.Run("Codebase Backup Monitor Always Generated", func(t *testing.T) {
		input := &appdef.Definition{
			Project:   appdef.Project{Name: "test", Title: "Test Project"},
			Apps:      []appdef.App{},
			Resources: []appdef.Resource{},
		}

		tf := setupTfVars(t, input)
		monitors := tf.generateMonitors(env.Production)

		// Codebase backup monitor is always generated, regardless of apps or resources.
		require.Len(t, monitors, 1)
		assert.Equal(t, "Backup - Codebase", monitors[0].Name)
		assert.Equal(t, "push", monitors[0].Type)
		assert.Empty(t, monitors[0].URL, "Push monitors should have empty URL")
		assert.Empty(t, monitors[0].Method, "Push monitors should have empty method")
		assert.Equal(t, appdef.MonitorIntervalBackup, monitors[0].Interval)
	})

	t.Run("Mixed Apps And Resources", func(t *testing.T) {
		input := &appdef.Definition{
			Project: appdef.Project{Name: "test", Title: "Test Project"},
			Apps: []appdef.App{
				{
					Name:  "web",
					Title: "Web",
					Domains: []appdef.Domain{
						{Name: "example.com", Type: appdef.DomainTypePrimary},
					},
					Infra:      appdef.Infra{},
					Monitoring: ptr.BoolPtr(true),
				},
			},
			Resources: []appdef.Resource{
				{
					Name:       "db",
					Title:      "Database",
					Type:       appdef.ResourceTypePostgres,
					Monitoring: ptr.BoolPtr(true),
					Backup:     &appdef.ResourceBackupConfig{Enabled: true},
				},
			},
		}

		tf := setupTfVars(t, input)
		monitors := tf.generateMonitors(env.Production)
		require.Len(t, monitors, 4) // HTTP + DNS + Database Backup + Codebase Backup

		// HTTP monitor.
		assert.Equal(t, "HTTP - example.com", monitors[0].Name)
		assert.Equal(t, "http", monitors[0].Type)
		assert.Equal(t, appdef.MonitorIntervalHTTP, monitors[0].Interval)

		// DNS monitor.
		assert.Equal(t, "DNS - example.com", monitors[1].Name)
		assert.Equal(t, "dns", monitors[1].Type)
		assert.Equal(t, appdef.MonitorIntervalDNS, monitors[1].Interval)

		// Database backup monitor.
		assert.Equal(t, "Backup - Database", monitors[2].Name)
		assert.Equal(t, "push", monitors[2].Type)
		assert.Equal(t, appdef.MonitorIntervalBackup, monitors[2].Interval)
		assert.Equal(t, "PROD_DB_BACKUP_PING_URL", monitors[2].VariableName, "VariableName should use resource Name, not Title")

		// Codebase backup monitor.
		assert.Equal(t, "Backup - Codebase", monitors[3].Name)
		assert.Equal(t, "push", monitors[3].Type)
		assert.Equal(t, appdef.MonitorIntervalBackup, monitors[3].Interval)
		assert.Equal(t, "PROD_CODEBASE_BACKUP_PING_URL", monitors[3].VariableName)
	})
}
