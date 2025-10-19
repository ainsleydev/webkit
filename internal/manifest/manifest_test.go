package manifest

import (
	"bytes"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/printer"
)

func TestCleanup(t *testing.T) {
	t.Parallel()

	t.Run("Scaffold Files Skipped", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		var buf bytes.Buffer
		console := printer.New(&buf)

		err := afero.WriteFile(fs, "user.env", []byte("content"), 0o644)
		require.NoError(t, err)

		old := &Manifest{
			Files: map[string]FileEntry{
				"user.env": {
					ScaffoldMode: true,
				},
			},
		}

		mani := &Manifest{Files: map[string]FileEntry{}}

		err = Cleanup(fs, old, mani, console)
		assert.NoError(t, err)

		exists, err := afero.Exists(fs, "user.env")
		require.NoError(t, err)
		assert.True(t, exists, "scaffold mode files should not be removed")
	})

	t.Run("File Exists In New Manifest - Skip", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		var buf bytes.Buffer
		console := printer.New(&buf)

		err := afero.WriteFile(fs, "keep.txt", []byte("content"), 0o644)
		require.NoError(t, err)

		old := &Manifest{
			Files: map[string]FileEntry{
				"keep.txt": {ScaffoldMode: false},
			},
		}

		mani := &Manifest{
			Files: map[string]FileEntry{
				"keep.txt": {ScaffoldMode: false},
			},
		}

		err = Cleanup(fs, old, mani, console)
		assert.NoError(t, err)

		exists, err := afero.Exists(fs, "keep.txt")
		require.NoError(t, err)
		assert.True(t, exists, "files in new manifest should not be removed")
	})

	t.Run("Remove Error", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		var buf bytes.Buffer
		console := printer.New(&buf)
		err := afero.WriteFile(fs, "orphaned.txt", []byte("old"), 0o644)
		require.NoError(t, err)
		old := &Manifest{
			Files: map[string]FileEntry{
				"orphaned.txt": {ScaffoldMode: false},
			},
		}
		mani := &Manifest{Files: map[string]FileEntry{}}

		readOnlyFs := afero.NewReadOnlyFs(fs)

		err = Cleanup(readOnlyFs, old, mani, console)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "removing orphaned.txt")
	})

	t.Run("Remove Orphaned Files", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		var buf bytes.Buffer
		console := printer.New(&buf)

		err := afero.WriteFile(fs, "orphaned.txt", []byte("old"), 0o644)
		require.NoError(t, err)

		old := &Manifest{
			Files: map[string]FileEntry{
				"orphaned.txt": {ScaffoldMode: false},
			},
		}

		mani := &Manifest{Files: map[string]FileEntry{}}

		err = Cleanup(fs, old, mani, console)
		assert.NoError(t, err)

		exists, err := afero.Exists(fs, "orphaned.txt")
		require.NoError(t, err)
		assert.False(t, exists, "orphaned files should be removed")

		output := buf.String()
		assert.Contains(t, output, "Removing orphaned: orphaned.txt")
		assert.Contains(t, output, "Removed: orphaned.txt")
	})
}

func TestHashContent(t *testing.T) {
	t.Parallel()

	t.Run("Empty Content", func(t *testing.T) {
		t.Parallel()

		got := HashContent([]byte{})
		assert.Equal(t, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", got)
	})

	t.Run("Consistent Hashing", func(t *testing.T) {
		t.Parallel()

		data := []byte("test content")
		hash1 := HashContent(data)
		hash2 := HashContent(data)

		assert.Equal(t, hash1, hash2)
	})

	t.Run("Different Content Different Hash", func(t *testing.T) {
		t.Parallel()

		hash1 := HashContent([]byte("content1"))
		hash2 := HashContent([]byte("content2"))
		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("Known Hash Value", func(t *testing.T) {
		t.Parallel()

		data := []byte("hello world")
		got := HashContent(data)
		want := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
		assert.Equal(t, want, got)
	})
}
