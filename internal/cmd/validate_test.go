package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/pkg/util/ptr"
)

func TestValidate(t *testing.T) {
	t.Parallel()

	t.Run("Valid Definition", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll("/apps/test", 0o755)
		assert.NoError(t, err)

		def := &appdef.Definition{
			WebkitVersion: "1.0.0",
			Project: appdef.Project{
				Name:        "test-project",
				Title:       "Test Project",
				Description: "Test description",
				Repo:        appdef.GitHubRepo{Owner: "test", Name: "repo"},
			},
			Apps: []appdef.App{
				{
					Name:  "test-app",
					Title: "Test App",
					Type:  appdef.AppTypeGoLang,
					Path:  "/apps/test",
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "vm",
					},
					Domains: []appdef.Domain{{Name: "example.com"}},
				},
			},
		}

		input, buf := setupWithPrinter(t, fs, def)

		err = validate(t.Context(), input)
		assert.NoError(t, err)
		assert.Contains(t, buf.String(), "Validating app.json...")
		assert.Contains(t, buf.String(), "Validation passed! No errors found.")
	})

	t.Run("Invalid Definition - Domain With Protocol", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll("/apps/test", 0o755)
		assert.NoError(t, err)

		def := &appdef.Definition{
			WebkitVersion: "1.0.0",
			Project: appdef.Project{
				Name:        "test-project",
				Title:       "Test Project",
				Description: "Test description",
				Repo:        appdef.GitHubRepo{Owner: "test", Name: "repo"},
			},
			Apps: []appdef.App{
				{
					Name:  "test-app",
					Title: "Test App",
					Type:  appdef.AppTypeGoLang,
					Path:  "/apps/test",
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "vm",
					},
					Domains: []appdef.Domain{{Name: "https://example.com"}},
				},
			},
		}

		input, buf := setupWithPrinter(t, fs, def)

		err = validate(t.Context(), input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed with 1 error(s)")
		assert.Contains(t, buf.String(), "Validation failed with 1 error(s):")
		assert.Contains(t, buf.String(), "should not contain protocol prefix")
	})

	t.Run("Invalid Definition - Non-existent Path", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		def := &appdef.Definition{
			WebkitVersion: "1.0.0",
			Project: appdef.Project{
				Name:        "test-project",
				Title:       "Test Project",
				Description: "Test description",
				Repo:        appdef.GitHubRepo{Owner: "test", Name: "repo"},
			},
			Apps: []appdef.App{
				{
					Name:  "test-app",
					Title: "Test App",
					Type:  appdef.AppTypeGoLang,
					Path:  "/apps/nonexistent",
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "vm",
					},
					Domains: []appdef.Domain{{Name: "example.com"}},
				},
			},
		}

		input, buf := setupWithPrinter(t, fs, def)

		err := validate(t.Context(), input)
		assert.Error(t, err)
		assert.Contains(t, buf.String(), "Validation failed")
		assert.Contains(t, buf.String(), "does not exist")
	})

	t.Run("Invalid Definition - Multiple Errors", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		def := &appdef.Definition{
			WebkitVersion: "1.0.0",
			Project: appdef.Project{
				Name:        "test-project",
				Title:       "Test Project",
				Description: "Test description",
				Repo:        appdef.GitHubRepo{Owner: "test", Name: "repo"},
			},
			Apps: []appdef.App{
				{
					Name:  "test-app",
					Title: "Test App",
					Type:  appdef.AppTypeGoLang,
					Path:  "/apps/nonexistent",
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "vm",
					},
					Domains: []appdef.Domain{{Name: "https://example.com"}},
				},
			},
		}

		input, buf := setupWithPrinter(t, fs, def)

		err := validate(t.Context(), input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed with 2 error(s)")
		assert.Contains(t, buf.String(), "Validation failed with 2 error(s):")
		assert.Contains(t, buf.String(), "should not contain protocol prefix")
		assert.Contains(t, buf.String(), "does not exist")
	})

	t.Run("Invalid Definition - Terraform-managed VM Without Domains", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll("/apps/test", 0o755)
		assert.NoError(t, err)

		def := &appdef.Definition{
			WebkitVersion: "1.0.0",
			Project: appdef.Project{
				Name:        "test-project",
				Title:       "Test Project",
				Description: "Test description",
				Repo:        appdef.GitHubRepo{Owner: "test", Name: "repo"},
			},
			Apps: []appdef.App{
				{
					Name:             "test-app",
					Title:            "Test App",
					Type:             appdef.AppTypeGoLang,
					Path:             "/apps/test",
					TerraformManaged: ptr.BoolPtr(true),
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "vm",
					},
					Domains: []appdef.Domain{},
				},
			},
		}

		input, buf := setupWithPrinter(t, fs, def)

		err = validate(t.Context(), input)
		assert.Error(t, err)
		assert.Contains(t, buf.String(), "Validation failed")
		assert.Contains(t, buf.String(), "terraform-managed VM/app must have at least one domain configured")
	})

	t.Run("Invalid Definition - Invalid Env Resource Reference", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll("/apps/test", 0o755)
		assert.NoError(t, err)

		def := &appdef.Definition{
			WebkitVersion: "1.0.0",
			Project: appdef.Project{
				Name:        "test-project",
				Title:       "Test Project",
				Description: "Test description",
				Repo:        appdef.GitHubRepo{Owner: "test", Name: "repo"},
			},
			Apps: []appdef.App{
				{
					Name:  "test-app",
					Title: "Test App",
					Type:  appdef.AppTypeGoLang,
					Path:  "/apps/test",
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "vm",
					},
					Domains: []appdef.Domain{{Name: "example.com"}},
					Env: appdef.Environment{
						Production: appdef.EnvVar{
							"DATABASE_URL": appdef.EnvValue{
								Source: appdef.EnvSourceResource,
								Value:  "nonexistent.connection_url",
							},
						},
					},
				},
			},
			Resources: []appdef.Resource{
				{
					Name:     "db",
					Type:     appdef.ResourceTypePostgres,
					Provider: appdef.ResourceProviderDigitalOcean,
				},
			},
		}

		input, buf := setupWithPrinter(t, fs, def)

		err = validate(t.Context(), input)
		assert.Error(t, err)
		assert.Contains(t, buf.String(), "Validation failed")
		assert.Contains(t, buf.String(), "references non-existent resource")
	})
}
