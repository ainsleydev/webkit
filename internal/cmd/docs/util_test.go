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

func TestParseContentWithFrontMatter(t *testing.T) {
	t.Parallel()

	t.Run("File does not exist", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		var meta readmeFrontMatter
		got, err := parseContentWithFrontMatter(fs, "docs/nonexistent.md", &meta)
		require.NoError(t, err)
		assert.Equal(t, "", got)
		assert.Nil(t, meta.Logo)
	})

	t.Run("File without front matter", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll(customDocsDir, 0o755)
		require.NoError(t, err)

		content := "Plain markdown content"
		err = afero.WriteFile(fs,
			filepath.Join(customDocsDir, "test.md"),
			[]byte(content),
			0o644,
		)
		require.NoError(t, err)

		var meta readmeFrontMatter
		got, err := parseContentWithFrontMatter(fs, filepath.Join(customDocsDir, "test.md"), &meta)
		require.NoError(t, err)
		assert.Equal(t, content, got)
		assert.Nil(t, meta.Logo)
	})

	t.Run("File with YAML front matter", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll(customDocsDir, 0o755)
		require.NoError(t, err)

		content := `---
logo:
  width: 500
  height: 250
---
Content after front matter`
		err = afero.WriteFile(fs,
			filepath.Join(customDocsDir, "test.md"),
			[]byte(content),
			0o644,
		)
		require.NoError(t, err)

		var meta readmeFrontMatter
		got, err := parseContentWithFrontMatter(fs, filepath.Join(customDocsDir, "test.md"), &meta)
		require.NoError(t, err)
		assert.Equal(t, "Content after front matter", got)
		require.NotNil(t, meta.Logo)
		assert.Equal(t, 500, meta.Logo.Width)
		assert.Equal(t, 250, meta.Logo.Height)
	})

	t.Run("Different struct type", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll(customDocsDir, 0o755)
		require.NoError(t, err)

		type CustomMeta struct {
			Title string `yaml:"title"`
			Draft bool   `yaml:"draft"`
		}

		content := `---
title: Test Document
draft: true
---
Custom content`
		err = afero.WriteFile(fs,
			filepath.Join(customDocsDir, "custom.md"),
			[]byte(content),
			0o644,
		)
		require.NoError(t, err)

		var meta CustomMeta
		got, err := parseContentWithFrontMatter(fs, filepath.Join(customDocsDir, "custom.md"), &meta)
		require.NoError(t, err)
		assert.Equal(t, "Custom content", got)
		assert.Equal(t, "Test Document", meta.Title)
		assert.True(t, meta.Draft)
	})

	t.Run("FS error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mock := mocks.NewMockFS(ctrl)
		mock.EXPECT().
			Open(gomock.Any()).
			Return(nil, errors.New("disk error"))

		var meta readmeFrontMatter
		_, err := parseContentWithFrontMatter(mock, "docs/test.md", &meta)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "reading test.md")
	})
}

func TestLoadReadmeContent(t *testing.T) {
	t.Parallel()

	t.Run("File does not exist", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		got, err := loadReadmeContent(fs)
		require.NoError(t, err)
		assert.Equal(t, "", got.Content)
		assert.Nil(t, got.Meta.Logo)
	})

	t.Run("Content without front matter", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll(customDocsDir, 0o755)
		require.NoError(t, err)

		content := "This is my custom README content."
		err = afero.WriteFile(fs,
			filepath.Join(customDocsDir, "README.md"),
			[]byte(content),
			0o644,
		)
		require.NoError(t, err)

		got, err := loadReadmeContent(fs)
		require.NoError(t, err)
		assert.Equal(t, content, got.Content)
		assert.Nil(t, got.Meta.Logo)
	})

	t.Run("Content with front matter logo width only", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll(customDocsDir, 0o755)
		require.NoError(t, err)

		content := `---
logo:
  width: 400
---
This is my custom README content.`
		err = afero.WriteFile(fs,
			filepath.Join(customDocsDir, "README.md"),
			[]byte(content),
			0o644,
		)
		require.NoError(t, err)

		got, err := loadReadmeContent(fs)
		require.NoError(t, err)
		assert.Equal(t, "This is my custom README content.", got.Content)
		require.NotNil(t, got.Meta.Logo)
		assert.Equal(t, 400, got.Meta.Logo.Width)
		assert.Equal(t, 0, got.Meta.Logo.Height)
	})

	t.Run("Content with front matter width and height", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		err := fs.MkdirAll(customDocsDir, 0o755)
		require.NoError(t, err)

		content := `---
logo:
  width: 300
  height: 150
---
This is my custom README content with both dimensions.`
		err = afero.WriteFile(fs,
			filepath.Join(customDocsDir, "README.md"),
			[]byte(content),
			0o644,
		)
		require.NoError(t, err)

		got, err := loadReadmeContent(fs)
		require.NoError(t, err)
		assert.Equal(t, "This is my custom README content with both dimensions.", got.Content)
		require.NotNil(t, got.Meta.Logo)
		assert.Equal(t, 300, got.Meta.Logo.Width)
		assert.Equal(t, 150, got.Meta.Logo.Height)
	})

	t.Run("FS error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mock := mocks.NewMockFS(ctrl)
		mock.EXPECT().
			Open(gomock.Any()).
			Return(nil, errors.New("read error"))

		_, err := loadReadmeContent(mock)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "reading README.md")
	})
}
