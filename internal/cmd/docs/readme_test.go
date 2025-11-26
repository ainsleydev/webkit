package docs

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/state/outputs"
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

func TestGetPeekapingEndpoint(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input *outputs.WebkitOutputs
		want  string
	}{
		"Nil outputs": {
			input: nil,
			want:  "https://uptime.ainsley.dev",
		},
		"Empty endpoint": {
			input: &outputs.WebkitOutputs{
				Peekaping: outputs.Peekaping{
					Endpoint: "",
				},
			},
			want: "https://uptime.ainsley.dev",
		},
		"Custom endpoint": {
			input: &outputs.WebkitOutputs{
				Peekaping: outputs.Peekaping{
					Endpoint: "https://peekaping.example.com",
				},
			},
			want: "https://peekaping.example.com",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := getPeekapingEndpoint(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGetLogoWidth(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input *ReadmeContent
		want  int
	}{
		"Nil logo": {
			input: &ReadmeContent{},
			want:  0,
		},
		"Logo with width": {
			input: &ReadmeContent{
				Meta: ReadmeFrontMatter{
					Logo: &LogoConfig{
						Width: 400,
					},
				},
			},
			want: 400,
		},
		"Logo with width and height": {
			input: &ReadmeContent{
				Meta: ReadmeFrontMatter{
					Logo: &LogoConfig{
						Width:  300,
						Height: 150,
					},
				},
			},
			want: 300,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := getLogoWidth(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGetLogoHeight(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input *ReadmeContent
		want  int
	}{
		"Nil logo": {
			input: &ReadmeContent{},
			want:  0,
		},
		"Logo with height": {
			input: &ReadmeContent{
				Meta: ReadmeFrontMatter{
					Logo: &LogoConfig{
						Height: 200,
					},
				},
			},
			want: 200,
		},
		"Logo with width and height": {
			input: &ReadmeContent{
				Meta: ReadmeFrontMatter{
					Logo: &LogoConfig{
						Width:  300,
						Height: 150,
					},
				},
			},
			want: 150,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := getLogoHeight(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGetDashboardURL(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input *outputs.WebkitOutputs
		want  string
	}{
		"Nil outputs": {
			input: nil,
			want:  "",
		},
		"Empty project tag": {
			input: &outputs.WebkitOutputs{
				Peekaping: outputs.Peekaping{
					Endpoint:   "https://peekaping.example.com",
					ProjectTag: "",
				},
			},
			want: "https://peekaping.example.com/monitors",
		},
		"Valid project tag": {
			input: &outputs.WebkitOutputs{
				Peekaping: outputs.Peekaping{
					Endpoint:   "https://peekaping.example.com",
					ProjectTag: "test-tag-123",
				},
			},
			want: "https://peekaping.example.com/monitors?tags=test-tag-123",
		},
		"Empty endpoint with project tag": {
			input: &outputs.WebkitOutputs{
				Peekaping: outputs.Peekaping{
					Endpoint:   "",
					ProjectTag: "test-tag-456",
				},
			},
			want: "https://uptime.ainsley.dev/monitors?tags=test-tag-456",
		},
		"Empty endpoint without project tag": {
			input: &outputs.WebkitOutputs{
				Peekaping: outputs.Peekaping{
					Endpoint:   "",
					ProjectTag: "",
				},
			},
			want: "https://uptime.ainsley.dev/monitors",
		},
		"Real world example": {
			input: &outputs.WebkitOutputs{
				Peekaping: outputs.Peekaping{
					Endpoint:   "https://uptime.ainsley.dev",
					ProjectTag: "08ba3cee-0afb-4d51-815e-daca3f2172f2",
				},
			},
			want: "https://uptime.ainsley.dev/monitors?tags=08ba3cee-0afb-4d51-815e-daca3f2172f2",
		},
		"Project tag with spaces": {
			input: &outputs.WebkitOutputs{
				Peekaping: outputs.Peekaping{
					Endpoint:   "https://uptime.ainsley.dev",
					ProjectTag: "tag with spaces",
				},
			},
			want: "https://uptime.ainsley.dev/monitors?tags=tag+with+spaces",
		},
		"Project tag with special characters": {
			input: &outputs.WebkitOutputs{
				Peekaping: outputs.Peekaping{
					Endpoint:   "https://uptime.ainsley.dev",
					ProjectTag: "tag&special=chars",
				},
			},
			want: "https://uptime.ainsley.dev/monitors?tags=tag%26special%3Dchars",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := getDashboardURL(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
