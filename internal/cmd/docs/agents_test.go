package docs

//
//import (
//	"testing"
//
//	"github.com/spf13/afero"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
//
//	"github.com/ainsleydev/webkit/internal/appdef"
//)
//
//func TestGenerate(t *testing.T) {
//	t.Parallel()
//
//	t.Run("With no custom content", func(t *testing.T) {
//		t.Parallel()
//
//		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{})
//
//		err := Agents(t.Context(), input)
//		require.NoError(t, err)
//
//		got, err := afero.ReadFile(input.FS, "AGENTS.md")
//		require.NoError(t, err)
//		assert.NotEmpty(t, got)
//		assert.Contains(t, string(got), "# Agent Guidelines")
//	})
//
//	t.Run("With custom content from docs/AGENTS.md", func(t *testing.T) {
//		t.Parallel()
//
//		fs := afero.NewMemMapFs()
//		customContent := "## Custom Project Rules\n\nThis is custom content for the project."
//
//		err := afero.WriteFile(fs, agentsPath, []byte(customContent), 0o644)
//		require.NoError(t, err)
//
//		input := setup(t, fs, &appdef.Definition{})
//
//		err = Agents(t.Context(), input)
//		require.NoError(t, err)
//
//		got, err := afero.ReadFile(input.FS, "AGENTS.md")
//		require.NoError(t, err)
//		assert.NotEmpty(t, got)
//		assert.Contains(t, string(got), "# Agent Guidelines")
//		assert.Contains(t, string(got), customContent)
//	})
//
//	t.Run("With custom template from docs/AGENTS.md.tmpl", func(t *testing.T) {
//		t.Parallel()
//
//		fs := afero.NewMemMapFs()
//		customTemplate := "## App Name: {{ .Definition.Project.Name }}\n\nThis is a template."
//
//		err := afero.WriteFile(fs, agentsPathTpl, []byte(customTemplate), 0o644)
//		require.NoError(t, err)
//
//		input := setup(t, fs, &appdef.Definition{
//			Project: appdef.Project{
//				Name: "test-app",
//			},
//		})
//
//		err = Agents(t.Context(), input)
//		require.NoError(t, err)
//
//		got, err := afero.ReadFile(input.FS, "AGENTS.md")
//		require.NoError(t, err)
//		assert.NotEmpty(t, got)
//		assert.Contains(t, string(got), "# Agent Guidelines")
//		assert.Contains(t, string(got), "## App Name: test-app")
//	})
//
//	t.Run("FS Failure", func(t *testing.T) {
//		t.Parallel()
//
//		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), &appdef.Definition{})
//
//		err := Agents(t.Context(), input)
//		assert.Error(t, err)
//	})
//}
