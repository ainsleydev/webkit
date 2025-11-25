package docs

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestReadme(t *testing.T) {
	t.Parallel()

	t.Run("With no custom content", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{
			WebkitVersion: "v0.1.0",
			Project: appdef.Project{
				Name:        "test-project",
				Title:       "Test Project",
				Description: "A test project for README generation.",
				Repo: appdef.GitHubRepo{
					Owner: "testuser",
					Name:  "test-repo",
				},
			},
		}
		input := setup(t, fs, appDef)

		err := Readme(context.Background(), input)
		require.NoError(t, err)

		got, err := afero.ReadFile(fs, "README.md")
		require.NoError(t, err)
		assert.NotEmpty(t, got)
		assert.Contains(t, string(got), "# Test Project")
		assert.Contains(t, string(got), "A test project for README generation.")
		assert.Contains(t, string(got), "Built with [WebKit v0.1.0]")
	})

	t.Run("With custom content from docs/README.md", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		customContent := "## Custom Section\n\nThis is custom content for the README."

		err := fs.MkdirAll("docs", 0o755)
		require.NoError(t, err)

		err = afero.WriteFile(fs, "docs/README.md", []byte(customContent), 0o644)
		require.NoError(t, err)

		appDef := &appdef.Definition{
			WebkitVersion: "v0.1.0",
			Project: appdef.Project{
				Name:        "test-project",
				Title:       "Test Project",
				Description: "A test project.",
				Repo: appdef.GitHubRepo{
					Owner: "testuser",
					Name:  "test-repo",
				},
			},
		}
		input := setup(t, fs, appDef)

		err = Readme(context.Background(), input)
		require.NoError(t, err)

		got, err := afero.ReadFile(fs, "README.md")
		require.NoError(t, err)
		assert.NotEmpty(t, got)
		assert.Contains(t, string(got), customContent)
	})

	t.Run("With logo in resources folder", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		err := fs.MkdirAll("resources", 0o755)
		require.NoError(t, err)

		err = afero.WriteFile(fs, "resources/logo.png", []byte("fake-logo"), 0o644)
		require.NoError(t, err)

		appDef := &appdef.Definition{
			WebkitVersion: "v0.1.0",
			Project: appdef.Project{
				Name:        "test-project",
				Title:       "Test Project",
				Description: "A test project.",
				Repo: appdef.GitHubRepo{
					Owner: "testuser",
					Name:  "test-repo",
				},
			},
		}
		input := setup(t, fs, appDef)

		err = Readme(context.Background(), input)
		require.NoError(t, err)

		got, err := afero.ReadFile(fs, "README.md")
		require.NoError(t, err)
		assert.NotEmpty(t, got)
		assert.Contains(t, string(got), "./resources/logo.png")
		assert.NotContains(t, string(got), webkitSymbolURL)
	})

	t.Run("With apps and resources", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		appDef := &appdef.Definition{
			WebkitVersion: "v0.1.0",
			Project: appdef.Project{
				Name:        "test-project",
				Title:       "Test Project",
				Description: "A test project.",
				Repo: appdef.GitHubRepo{
					Owner: "testuser",
					Name:  "test-repo",
				},
			},
			Apps: []appdef.App{
				{
					Name:        "web",
					Title:       "Web App",
					Type:        appdef.AppTypeSvelteKit,
					Description: "Web application.",
					Path:        "web",
					Build: appdef.Build{
						Port: 3000,
					},
					Infra: appdef.Infra{
						Provider: appdef.ResourceProviderDigitalOcean,
						Type:     "vm",
					},
					Domains: []appdef.Domain{
						{
							Name: "example.com",
							Type: appdef.DomainTypePrimary,
						},
					},
				},
			},
			Resources: []appdef.Resource{
				{
					Name:        "db",
					Type:        appdef.ResourceTypePostgres,
					Description: "PostgreSQL database for application data.",
					Provider:    appdef.ResourceProviderDigitalOcean,
					Config: map[string]any{
						"size": "db-s-1vcpu-1gb",
					},
					Backup: appdef.ResourceBackupConfig{
						Enabled: true,
					},
				},
			},
		}
		input := setup(t, fs, appDef)

		err := Readme(context.Background(), input)
		require.NoError(t, err)

		got, err := afero.ReadFile(fs, "README.md")
		require.NoError(t, err)
		assert.NotEmpty(t, got)
		assert.Contains(t, string(got), "Web App")
		assert.Contains(t, string(got), "svelte-kit")
		assert.Contains(t, string(got), "## Resources")
		assert.Contains(t, string(got), "PostgreSQL database for application data.")
		assert.Contains(t, string(got), "example.com")
	})

	t.Run("With status badge from outputs", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		// Create outputs.json with monitoring data
		outputsJSON := `{
			"peekaping": {
				"endpoint": "https://peekaping.example.com",
				"project_tag": "test-tag-123"
			},
			"monitors": [
				{"id": "mon123", "name": "HTTP - example.com", "type": "http"},
				{"id": "mon456", "name": "DNS - example.com", "type": "dns"}
			],
			"slack": {"channel_name": "alerts", "channel_id": "C123"}
		}`
		err := fs.MkdirAll(".webkit", 0o755)
		require.NoError(t, err)
		err = afero.WriteFile(fs, ".webkit/outputs.json", []byte(outputsJSON), 0o644)
		require.NoError(t, err)

		appDef := &appdef.Definition{
			WebkitVersion: "v0.1.0",
			Project: appdef.Project{
				Name:        "test-project",
				Title:       "Test Project",
				Description: "A test project.",
				Repo: appdef.GitHubRepo{
					Owner: "testuser",
					Name:  "test-repo",
				},
			},
		}
		input := setup(t, fs, appDef)

		err = Readme(context.Background(), input)
		require.NoError(t, err)

		got, err := afero.ReadFile(fs, "README.md")
		require.NoError(t, err)

		// Verify Status section
		assert.Contains(t, string(got), "## Status")
		assert.Contains(t, string(got), "status page")
		assert.Contains(t, string(got), "uptime.ainsley.dev") // default status page URL
		assert.Contains(t, string(got), "dashboard")
		assert.Contains(t, string(got), "test-tag-123") // verify dashboard link contains project tag
		assert.Contains(t, string(got), "HTTP - example.com")
		assert.Contains(t, string(got), "DNS - example.com")
		assert.Contains(t, string(got), "mon123")
		assert.Contains(t, string(got), "mon456")
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			WebkitVersion: "v0.1.0",
			Project: appdef.Project{
				Name:  "test",
				Title: "Test",
				Repo: appdef.GitHubRepo{
					Owner: "test",
					Name:  "test",
				},
			},
		}
		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), appDef)

		got := Readme(t.Context(), input)
		assert.Error(t, got)
	})
}
