package files

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestCodeStyle(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{})

		got := CodeStyle(t.Context(), input)
		assert.NoError(t, got)

		for path := range codeStyleTemplates {
			file, err := afero.ReadFile(input.FS, path)
			assert.NoError(t, err)
			assert.NotEmpty(t, file)
		}
	})

	t.Run("With Go", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "api",
					Type: appdef.AppTypeGoLang,
				},
			},
		})

		got := CodeStyle(t.Context(), input)
		assert.NoError(t, got)

		file, err := afero.ReadFile(input.FS, ".golangci.yaml")
		assert.NoError(t, err)
		assert.NotEmpty(t, file)
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), &appdef.Definition{})

		got := CodeStyle(t.Context(), input)
		assert.Error(t, got)
	})
}
