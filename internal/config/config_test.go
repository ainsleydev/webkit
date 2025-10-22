package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDir(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	dir, err := Dir()
	require.NoError(t, err)

	expected := filepath.Join(tmpHome, ".config", "webkit")
	assert.Equal(t, expected, dir)
}

func TestPath(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	t.Run("Single File", func(t *testing.T) {
		path, err := Path("age.key")
		require.NoError(t, err)

		expected := filepath.Join(tmpHome, ".config", "webkit", "age.key")
		assert.Equal(t, expected, path)
	})

	t.Run("Nested File", func(t *testing.T) {
		path, err := Path("secrets/prod.yaml")
		require.NoError(t, err)

		expected := filepath.Join(tmpHome, ".config", "webkit", "secrets", "prod.yaml")
		assert.Equal(t, expected, path)
	})
}

func TestReadWrite(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	t.Run("Write And Read", func(t *testing.T) {
		content := []byte("test content")

		err := Write("test.txt", content, 0o600)
		require.NoError(t, err)

		read, err := Read("test.txt")
		require.NoError(t, err)
		assert.Equal(t, content, read)
	})

	t.Run("Creates Config Dir Automatically", func(t *testing.T) {
		// Config dir shouldn't exist yet for this test
		newHome := t.TempDir()
		t.Setenv("HOME", newHome)

		err := Write("auto-create.txt", []byte("test"), 0o600)
		require.NoError(t, err)

		// Dir should exist now
		dir, err := Dir()
		require.NoError(t, err)

		info, err := os.Stat(dir)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
	})

	t.Run("Overwrites Existing File", func(t *testing.T) {
		err := Write("overwrite.txt", []byte("first"), 0o600)
		require.NoError(t, err)

		err = Write("overwrite.txt", []byte("second"), 0o600)
		require.NoError(t, err)

		read, err := Read("overwrite.txt")
		require.NoError(t, err)
		assert.Equal(t, []byte("second"), read)
	})

	t.Run("Read Non-Existent File", func(t *testing.T) {
		_, err := Read("does-not-exist.txt")
		assert.Error(t, err)
	})

	t.Run("File Permissions", func(t *testing.T) {
		err := Write("perms.txt", []byte("test"), 0o600)
		require.NoError(t, err)

		path, err := Path("perms.txt")
		require.NoError(t, err)

		info, err := os.Stat(path)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	})
}

func TestEnsureDir(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)

	// Directory shouldn't exist yet
	dir, err := Dir()
	require.NoError(t, err)
	_, err = os.Stat(dir)
	assert.True(t, os.IsNotExist(err))

	// Create it
	err = ensureDir()
	require.NoError(t, err)

	// Should exist now
	info, err := os.Stat(dir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())

	// Should be idempotent
	err = ensureDir()
	assert.NoError(t, err)
}
