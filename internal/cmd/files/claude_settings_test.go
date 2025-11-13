package files

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestClaudeSettings(t *testing.T) {
	t.Parallel()

	appDef := &appdef.Definition{
		Project: appdef.Project{Name: "my-website"},
		Apps: []appdef.App{
			{
				Name: "web",
				Type: appdef.AppTypeGoLang,
				Path: "web",
			},
		},
	}

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), appDef)

		err := ClaudeSettings(t.Context(), input)
		assert.NoError(t, err)

		for path := range claudeSettingsTemplates {
			file, err := afero.ReadFile(input.FS, path)
			assert.NoError(t, err)
			assert.NotEmpty(t, file)
		}
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), &appdef.Definition{})

		got := ClaudeSettings(t.Context(), input)
		assert.Error(t, got)
	})
}
