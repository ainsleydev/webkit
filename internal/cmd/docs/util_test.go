package docs

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/mocks"
)

func TestLoadCustomContent(t *testing.T) {
	t.Parallel()

	t.Run("FS Error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mock := mocks.NewMockFS(ctrl)
		mock.EXPECT().
			Open(gomock.Any()).
			Return(nil, errors.New("read error"))

		_, err := loadCustomContent(mock, "AGENTS.md")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "read error")
	})

	t.Run("File does not exist", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		_, err := loadCustomContent(fs, "AGENTS.md")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "doc template does not exist")
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll(customDocsDir, 0o755)
		require.NoError(t, err)

		err = afero.WriteFile(fs,
			filepath.Join(customDocsDir, "AGENTS.md"),
			[]byte("## WebKit\n\nCustom content."),
			0o644,
		)
		require.NoError(t, err)

		got, err := loadCustomContent(fs, "AGENTS.md")
		assert.NoError(t, err)
		assert.Equal(t, "## WebKit\n\nCustom content.", got)
	})
}
