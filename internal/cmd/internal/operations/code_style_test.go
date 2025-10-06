package operations

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
)

type errCreateFs struct {
	afero.Fs
}

func (e *errCreateFs) Create(_ string) (afero.File, error) {
	return nil, fmt.Errorf("create error")
}

func (e *errCreateFs) OpenFile(_ string, _ int, _ os.FileMode) (afero.File, error) {
	return nil, fmt.Errorf("openfile error")
}

func TestCreateCodeStyleFiles(t *testing.T) {
	t.Parallel()

	t.Run("Creates Successfully", func(t *testing.T) {
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

	t.Run("Errors on Failure", func(t *testing.T) {
		t.Parallel()

		// Use an FS that fails to create files — generator should return an error.
		fs := &errCreateFs{Fs: afero.NewMemMapFs()}

		got := CreateCodeStyleFiles(t.Context(), cmdtools.CommandInput{
			FS:          fs,
			AppDefCache: &appdef.Definition{},
		})
		assert.Error(t, got)
	})
}
