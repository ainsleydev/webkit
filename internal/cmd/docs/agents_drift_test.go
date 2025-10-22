package docs

import (
	"context"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/manifest"
)

// TestAgentsDeterminism verifies that AGENTS.md generation is deterministic
// and doesn't cause false drift warnings.
func TestAgentsDeterminism(t *testing.T) {
	t.Parallel()

	t.Run("AGENTS.md generation is deterministic", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-app",
			},
		}

		// Generate AGENTS.md first time
		fs1 := afero.NewMemMapFs()
		input1 := setup(t, fs1, appDef)
		err := Agents(t.Context(), input1)
		require.NoError(t, err)

		content1, err := afero.ReadFile(fs1, "AGENTS.md")
		require.NoError(t, err)
		hash1 := manifest.HashContent(content1)

		// Generate AGENTS.md second time with same app definition
		fs2 := afero.NewMemMapFs()
		input2 := setup(t, fs2, appDef)
		err = Agents(t.Context(), input2)
		require.NoError(t, err)

		content2, err := afero.ReadFile(fs2, "AGENTS.md")
		require.NoError(t, err)
		hash2 := manifest.HashContent(content2)

		// Hashes should match for deterministic generation
		assert.Equal(t, hash1, hash2, "AGENTS.md generation should be deterministic")
		assert.Equal(t, string(content1), string(content2), "Content should be identical")
	})

	t.Run("AGENTS.md with custom content is deterministic", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-app",
			},
		}

		customContent := "## Custom Rules\n\nThese are project-specific rules."

		// Generate first time
		fs1 := afero.NewMemMapFs()
		require.NoError(t, afero.WriteFile(fs1, agentsPath, []byte(customContent), 0o644))
		input1 := setup(t, fs1, appDef)
		err := Agents(t.Context(), input1)
		require.NoError(t, err)

		content1, err := afero.ReadFile(fs1, "AGENTS.md")
		require.NoError(t, err)
		hash1 := manifest.HashContent(content1)

		// Generate second time
		fs2 := afero.NewMemMapFs()
		require.NoError(t, afero.WriteFile(fs2, agentsPath, []byte(customContent), 0o644))
		input2 := setup(t, fs2, appDef)
		err = Agents(t.Context(), input2)
		require.NoError(t, err)

		content2, err := afero.ReadFile(fs2, "AGENTS.md")
		require.NoError(t, err)
		hash2 := manifest.HashContent(content2)

		// Hashes should match
		assert.Equal(t, hash1, hash2, "AGENTS.md with custom content should be deterministic")
		assert.Equal(t, string(content1), string(content2), "Content should be identical")
	})

	t.Run("AGENTS.md with template is deterministic", func(t *testing.T) {
		t.Parallel()

		appDef := &appdef.Definition{
			Project: appdef.Project{
				Name: "test-app",
			},
		}

		customTemplate := "## App: {{ .Definition.Project.Name }}\n\nRules for this app."

		// Generate first time
		fs1 := afero.NewMemMapFs()
		require.NoError(t, afero.WriteFile(fs1, agentsPathTpl, []byte(customTemplate), 0o644))
		input1 := setup(t, fs1, appDef)
		err := Agents(context.Background(), input1)
		require.NoError(t, err)

		content1, err := afero.ReadFile(fs1, "AGENTS.md")
		require.NoError(t, err)
		hash1 := manifest.HashContent(content1)

		// Generate second time
		fs2 := afero.NewMemMapFs()
		require.NoError(t, afero.WriteFile(fs2, agentsPathTpl, []byte(customTemplate), 0o644))
		input2 := setup(t, fs2, appDef)
		err = Agents(context.Background(), input2)
		require.NoError(t, err)

		content2, err := afero.ReadFile(fs2, "AGENTS.md")
		require.NoError(t, err)
		hash2 := manifest.HashContent(content2)

		// Hashes should match
		assert.Equal(t, hash1, hash2, "AGENTS.md with template should be deterministic")
		assert.Equal(t, string(content1), string(content2), "Content should be identical")
	})
}
