package fsext

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyAllEmbed(t *testing.T) {
	t.Parallel()

	dest, err := os.MkdirTemp("", "copy-all")
	require.NoError(t, err)

	err = CopyAllEmbed(testFS, dest)
	assert.NoError(t, err)

	var files []string
	_ = filepath.Walk(dest, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			files = append(files, filepath.Base(path))
		}
		return nil
	})

	assert.Contains(t, files, "one.txt")
	assert.Contains(t, files, "nested.txt")
}

func TestCopyFromEmbed(t *testing.T) {
	t.Parallel()

	t.Run("Valid Copy", func(t *testing.T) {
		t.Parallel()

		dest, err := os.MkdirTemp("", "copy-ok")
		require.NoError(t, err)

		err = CopyFromEmbed(testFS, "testdata", dest)
		assert.NoError(t, err)

		var files []string
		_ = filepath.Walk(dest, func(path string, info os.FileInfo, err error) error {
			if err == nil && !info.IsDir() {
				files = append(files, filepath.Base(path))
			}
			return nil
		})

		assert.Contains(t, files, "one.txt")
		assert.Contains(t, files, "nested.txt")
	})

	t.Run("Missing Source", func(t *testing.T) {
		t.Parallel()

		dest, err := os.MkdirTemp("", "copy-missing")
		require.NoError(t, err)

		err = CopyFromEmbed(testFS, "does-not-exist", dest)
		assert.Error(t, err)
	})

	t.Run("Permission Error", func(t *testing.T) {
		t.Parallel()

		dir := t.TempDir()
		noWrite := filepath.Join(dir, "readonly")
		err := os.Mkdir(noWrite, 0o500)
		require.NoError(t, err)

		err = CopyFromEmbed(testFS, "testdata", noWrite)
		assert.Error(t, err)
	})

	t.Run("Bad Relative Path", func(t *testing.T) {
		t.Parallel()

		dest, err := os.MkdirTemp("", "copy-rel")
		require.NoError(t, err)

		err = CopyFromEmbed(testFS, "testdata", dest)
		assert.NoError(t, err)
	})

}
