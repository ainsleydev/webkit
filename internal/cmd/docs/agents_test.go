package docs

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/fsext"
	"github.com/ainsleydev/webkit/internal/gen"
)

func TestAgents(t *testing.T) {
	t.Parallel()

	t.Run("With no custom content", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		input := setup(t, fs, &appdef.Definition{})

		err := Agents(context.Background(), input)
		require.NoError(t, err)

		got, err := afero.ReadFile(fs, "AGENTS.md")
		require.NoError(t, err)
		assert.NotEmpty(t, got)
		assert.Contains(t, string(got), "# Agent Guidelines")
	})

	t.Run("With custom content from docs/AGENTS.md", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		customContent := "## Custom Project Rules\n\nThis is custom content for the project."

		err := fs.MkdirAll("docs", 0o755)
		require.NoError(t, err)

		err = afero.WriteFile(fs, "docs/AGENTS.md", []byte(customContent), 0o644)
		require.NoError(t, err)

		input := setup(t, fs, &appdef.Definition{})

		err = Agents(context.Background(), input)
		require.NoError(t, err)

		got, err := afero.ReadFile(fs, "AGENTS.md")
		require.NoError(t, err)
		assert.NotEmpty(t, got)
		assert.Contains(t, string(got), "# Agent Guidelines")
		assert.Contains(t, string(got), customContent)
	})

	t.Run("With Payload App", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "payload-app",
					Type: appdef.AppTypePayload,
					Path: "cms",
				},
			},
		}

		input := setup(t, fs, appDef)

		err := Agents(context.Background(), input)
		require.NoError(t, err)

		appAgents, err := afero.ReadFile(fs, filepath.Join("cms", "AGENTS.md"))
		require.NoError(t, err)
		assert.NotEmpty(t, appAgents)
		assert.Contains(t, string(appAgents), fsext.MustReadFromEmbed(gen.Embed, "docs/PAYLOAD.md"))
	})

	t.Run("With SvelteKit App", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		appDef := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "svelte-app",
					Type: appdef.AppTypeSvelteKit,
					Path: "web",
				},
			},
		}

		input := setup(t, fs, appDef)

		err := Agents(context.Background(), input)
		require.NoError(t, err)

		appAgents, err := afero.ReadFile(fs, filepath.Join("web", "AGENTS.md"))
		require.NoError(t, err)
		assert.NotEmpty(t, appAgents)
		assert.Contains(t, string(appAgents), fsext.MustReadFromEmbed(gen.Embed, "docs/SVELTEKIT.md"))
	})
}
