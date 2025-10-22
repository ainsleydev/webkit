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

		file, err := afero.ReadFile(input.FS, outputPath)
		require.NoError(t, err)
		assert.NotEmpty(t, file)
		assert.Contains(t, string(file), "# Agent Guidelines")
		assert.Contains(t, string(file), "## WebKit")
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

		file, err := afero.ReadFile(input.FS, outputPath)
		require.NoError(t, err)
		assert.NotEmpty(t, file)
		assert.Contains(t, string(file), "# Agent Guidelines")
		assert.Contains(t, string(file), customContent)
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

		file, err := afero.ReadFile(input.FS, outputPath)
		require.NoError(t, err)
		assert.NotEmpty(t, file)
		assert.Contains(t, string(file), "# Agent Guidelines")
		assert.Contains(t, string(file), "## App Name: test-app")
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

	tt := map[string]struct {
		setupFS      func(t *testing.T) afero.Fs
		wantContains string
		wantErr      bool
	}{
		"No custom content": {
			setupFS: func(t *testing.T) afero.Fs {
				return afero.NewMemMapFs()
			},
			wantContains: "",
			wantErr:      false,
		},
		"Static markdown file": {
			setupFS: func(t *testing.T) afero.Fs {
				fs := afero.NewMemMapFs()
				err := afero.WriteFile(fs, customContentPath, []byte("# Custom Content"), 0644)
				require.NoError(t, err)
				return fs
			},
			wantContains: "# Custom Content",
			wantErr:      false,
		},
		"Template file with app name": {
			setupFS: func(t *testing.T) afero.Fs {
				fs := afero.NewMemMapFs()
				err := afero.WriteFile(fs, customContentPathTmpl, []byte("App: {{.Definition.Name}}"), 0644)
				require.NoError(t, err)
				return fs
			},
			wantContains: "App: test-app",
			wantErr:      false,
		},
		"Template file takes precedence": {
			setupFS: func(t *testing.T) afero.Fs {
				fs := afero.NewMemMapFs()
				err := afero.WriteFile(fs, customContentPathTmpl, []byte("Template content"), 0644)
				require.NoError(t, err)
				err = afero.WriteFile(fs, customContentPath, []byte("Static content"), 0644)
				require.NoError(t, err)
				return fs
			},
			wantContains: "Template content",
			wantErr:      false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs := test.setupFS(t)
			appDef := &appdef.Definition{Name: "test-app"}

			got, err := loadCustomContent(fs, appDef)

			assert.Equal(t, test.wantErr, err != nil)
			if test.wantContains != "" {
				assert.Contains(t, got, test.wantContains)
			} else {
				assert.Empty(t, got)
			}
		})
	}
}
