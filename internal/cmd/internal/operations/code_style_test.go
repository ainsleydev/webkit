package operations

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/mocks"
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

		ctrl := gomock.NewController(t)
		fsMock := mocks.NewMockFS(ctrl)
		fsMock.EXPECT().
			MkdirAll(gomock.Any(), gomock.Any()).
			Return(fmt.Errorf("mkdir error"))

		input := cmdtools.CommandInput{
			FS:          fsMock,
			AppDefCache: &appdef.Definition{},
		}

		got := CreateCodeStyleFiles(t.Context(), input)
		assert.Error(t, got)
	})
}
