package docsutil

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

func TestTemplate_String(t *testing.T) {
	t.Parallel()

	got := CodeStyleTemplate.String()
	assert.Equal(t, "CODE_STYLE.md", got)
	assert.IsType(t, "", got)
}

func TestLoadGenFile(t *testing.T) {
	t.Parallel()

	t.Run("FS Error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mock := mocks.NewMockFS(ctrl)
		mock.EXPECT().
			Open(gomock.Any()).
			Return(nil, errors.New("open error"))

		_, err := LoadGenFile(mock, CodeStyleTemplate)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "open error")
	})

	t.Run("File does not exist", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		_, err := LoadGenFile(fs, "MISSING.md")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "doc template does not exist")
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll(genDocsDir, 0755)
		require.NoError(t, err)

		err = afero.WriteFile(fs,
			filepath.Join(genDocsDir, CodeStyleTemplate.String()),
			[]byte("# Code Style"),
			0644,
		)
		require.NoError(t, err)

		got, err := LoadGenFile(fs, CodeStyleTemplate)
		assert.NoError(t, err)
		assert.Equal(t, "# Code Style", got)
	})
}

func TestMustLoadGenFile(t *testing.T) {
	t.Parallel()

	t.Run("File does not exist panics", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		assert.Panics(t, func() {
			MustLoadGenFile(fs, "MISSING.md")
		})
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll(genDocsDir, 0755)
		require.NoError(t, err)

		err = afero.WriteFile(fs,
			filepath.Join(genDocsDir, CodeStyleTemplate.String()),
			[]byte("# Code Style"),
			0644,
		)
		require.NoError(t, err)

		got := MustLoadGenFile(fs, "CODE_STYLE.md")
		assert.Equal(t, "# Code Style", got)
	})
}

func TestLoadCustomContent(t *testing.T) {
	t.Parallel()

	t.Run("FS Error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mock := mocks.NewMockFS(ctrl)
		mock.EXPECT().
			Open(gomock.Any()).
			Return(nil, errors.New("read error"))

		_, err := LoadCustomContent(mock)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "read error")
	})

	t.Run("File does not exist", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()

		_, err := LoadCustomContent(fs)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "doc template does not exist")
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll(customDocsDir, 0755)
		require.NoError(t, err)

		err = afero.WriteFile(fs,
			filepath.Join(customDocsDir, agentsFilename),
			[]byte("## WebKit\n\nCustom content."),
			0644,
		)
		require.NoError(t, err)

		got, err := LoadCustomContent(fs)
		assert.NoError(t, err)
		assert.Equal(t, "## WebKit\n\nCustom content.", got)
	})

}
