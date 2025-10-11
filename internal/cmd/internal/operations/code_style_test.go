package operations

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/testutil"
)

func TestCreateCodeStyleFiles(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		got := CreateCodeStyleFiles(t.Context(), cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: &appdef.Definition{},
		})
		assert.NoError(t, got)

		for path, _ := range codeStyleTemplates {
			file, err := afero.ReadFile(fs, path)
			assert.NoError(t, err)
			assert.NotEmpty(t, file)
		}
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		got := CreateCodeStyleFiles(t.Context(), cmdtools.CommandInput{
			FS:          &testutil.AferoErrCreateFs{Fs: afero.NewMemMapFs()},
			AppDefCache: &appdef.Definition{},
		})
		assert.Error(t, got)
	})
}
