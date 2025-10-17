package files

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
)

func TestCreateCodeStyleFiles(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		input := cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: &appdef.Definition{},
		}

		got := CreateCodeStyleFiles(t.Context(), input)
		assert.NoError(t, got)

		for path, _ := range codeStyleTemplates {
			file, err := afero.ReadFile(fs, path)
			assert.NoError(t, err)
			assert.NotEmpty(t, file)
		}
	})

	t.Run("FS Failure", func(t *testing.T) {
		t.Parallel()

		input := cmdtools.CommandInput{
			FS:          afero.NewReadOnlyFs(afero.NewMemMapFs()),
			AppDefCache: &appdef.Definition{},
		}

		got := CreateCodeStyleFiles(t.Context(), input)
		assert.Error(t, got)
	})
}
