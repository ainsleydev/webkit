package files

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/appdef"
)

func TestCreateCodeStyleFiles(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewMemMapFs(), &appdef.Definition{})

		got := CodeStyle(t.Context(), input)
		assert.NoError(t, got)

		for path, _ := range codeStyleTemplates {
			file, err := afero.ReadFile(input.FS, path)
			assert.NoError(t, err)
			assert.NotEmpty(t, file)
		}
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		input := setup(t, afero.NewReadOnlyFs(afero.NewMemMapFs()), &appdef.Definition{})

		got := CodeStyle(t.Context(), input)
		assert.Error(t, got)
	})
}
