package docs

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{
			Name: "test-app",
		})

		err := Generate(t.Context(), input)
		require.NoError(t, err)

		got, err := afero.ReadFile(input.FS, outputPath)
		require.NoError(t, err)
		assert.NotEmpty(t, got)
		assert.Contains(t, string(got), "# Agent Guidelines")
		assert.Contains(t, string(got), "## WebKit")
	})

	t.Run("With custom content from docs/AGENTS.md", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		customContent := "## Custom Project Rules\n\nThis is custom content for the project."

		err := afero.WriteFile(fs, customContentPath, []byte(customContent), 0644)
		require.NoError(t, err)

		input := setup(t, fs, &appdef.Definition{
			Name: "test-app",
		})

		err = Generate(t.Context(), input)
		require.NoError(t, err)

		got, err := afero.ReadFile(input.FS, outputPath)
		require.NoError(t, err)
		assert.NotEmpty(t, got)
		assert.Contains(t, string(got), "# Agent Guidelines")
		assert.Contains(t, string(got), customContent)
	})

	t.Run("With custom template from docs/AGENTS.md.tmpl", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		customTemplate := "## App Name: {{.Definition.Name}}\n\nThis is a template."

		err := afero.WriteFile(fs, customContentPathTmpl, []byte(customTemplate), 0644)
		require.NoError(t, err)

		input := setup(t, fs, &appdef.Definition{
			Name: "test-app",
		})

		err = Generate(t.Context(), input)
		require.NoError(t, err)

		got, err := afero.ReadFile(input.FS, outputPath)
		require.NoError(t, err)
		assert.NotEmpty(t, got)
		assert.Contains(t, string(got), "# Agent Guidelines")
		assert.Contains(t, string(got), "## App Name: test-app")
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), &appdef.Definition{})

		err := Generate(t.Context(), input)
		assert.Error(t, err)
	})
}

func TestLoadCustomContent(t *testing.T) {
	t.Parallel()

	t.Run("No custom content", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		appDef := &appdef.Definition{Name: "test-app"}

		got, err := loadCustomContent(fs, appDef)

		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("Static markdown file", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, customContentPath, []byte("# Custom Content"), 0644)
		require.NoError(t, err)

		appDef := &appdef.Definition{Name: "test-app"}

		got, err := loadCustomContent(fs, appDef)

		require.NoError(t, err)
		assert.Contains(t, got, "# Custom Content")
	})

	t.Run("Template file with app name", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, customContentPathTmpl, []byte("App: {{.Definition.Name}}"), 0644)
		require.NoError(t, err)

		appDef := &appdef.Definition{Name: "test-app"}

		got, err := loadCustomContent(fs, appDef)

		require.NoError(t, err)
		assert.Contains(t, got, "App: test-app")
	})

	t.Run("Template file takes precedence", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := afero.WriteFile(fs, customContentPathTmpl, []byte("Template content"), 0644)
		require.NoError(t, err)
		err = afero.WriteFile(fs, customContentPath, []byte("Static content"), 0644)
		require.NoError(t, err)

		appDef := &appdef.Definition{Name: "test-app"}

		got, err := loadCustomContent(fs, appDef)

		require.NoError(t, err)
		assert.Contains(t, got, "Template content")
	})
}
