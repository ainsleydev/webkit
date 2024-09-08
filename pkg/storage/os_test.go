package storage

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupOSStorage(t *testing.T) *OS {
	t.Helper()
	return NewOSStorage(t.TempDir())
}

func TestOS_Upload(t *testing.T) {
	t.Parallel()

	t.Run("Successful Upload", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)

		content := bytes.NewBufferString("test content")
		err := s.Upload(context.Background(), "test.txt", content)
		require.NoError(t, err)

		data, err := os.ReadFile(filepath.Join(s.BasePath, "test.txt"))
		require.NoError(t, err)
		assert.Equal(t, "test content", string(data))
	})

	t.Run("Upload to Sub Dir", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)

		content := bytes.NewBufferString("subdir content")
		err := s.Upload(context.Background(), "subdir/test.txt", content)
		require.NoError(t, err)

		data, err := os.ReadFile(filepath.Join(s.BasePath, "subdir/test.txt"))
		require.NoError(t, err)
		assert.Equal(t, "subdir content", string(data))
	})

	t.Run("Cannot Create Dir", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)

		// Create a file with the same name as the directory we want to create
		filePathConflict := filepath.Join(s.BasePath, "conflict")
		err := os.WriteFile(filePathConflict, []byte("conflict"), 0644)
		require.NoError(t, err)

		content := bytes.NewBufferString("conflict content")
		err = s.Upload(context.Background(), "conflict/test.txt", content)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a directory")
	})

	t.Run("Cannot Create File", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)

		// Make BasePath read-only to cause an error
		err := os.Chmod(s.BasePath, 0555)
		require.NoError(t, err)

		content := bytes.NewBufferString("error content")
		err = s.Upload(context.Background(), "error.txt", content)
		assert.Error(t, err)
	})
}

func TestOS_Delete(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)

		// Create a file to delete
		err := os.WriteFile(filepath.Join(s.BasePath, "delete_me.txt"), []byte("to be deleted"), 0644)
		require.NoError(t, err)

		err = s.Delete(context.Background(), "delete_me.txt")
		require.NoError(t, err)

		// Verify file no longer exists
		_, err = os.Stat(filepath.Join(s.BasePath, "delete_me.txt"))
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("Non-Existent File", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)

		err := s.Delete(context.Background(), "non_existent.txt")
		assert.Error(t, err)
	})
}

func TestOS_List(t *testing.T) {
	t.Parallel()

	// Create test files and directories
	write := func(s *OS) {
		require.NoError(t, os.MkdirAll(filepath.Join(s.BasePath, "dir1"), 0755))
		require.NoError(t, os.WriteFile(filepath.Join(s.BasePath, "file1.txt"), []byte("file1"), 0644))
		require.NoError(t, os.WriteFile(filepath.Join(s.BasePath, "dir1/file2.txt"), []byte("file2"), 0644))
	}

	t.Run("List", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)
		write(s)
		files, err := s.List(context.Background(), "")
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{"file1.txt", filepath.Join("dir1", "file2.txt")}, files)
	})

	t.Run("List With Prefix", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)
		write(s)
		files, err := s.List(context.Background(), "dir1")
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{filepath.Join("dir1", "file2.txt")}, files)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)
		files, err := s.List(context.Background(), "non_existent")
		assert.Error(t, err)
		assert.Nil(t, files)
	})
}

func TestOS_Download(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)

		// Create a file to download
		content := []byte("test content")
		require.NoError(t, os.WriteFile(filepath.Join(s.BasePath, "download.txt"), content, 0644))

		reader, err := s.Download(context.Background(), "download.txt")
		require.NoError(t, err)
		defer reader.Close()

		downloadedContent, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, content, downloadedContent)
	})

	t.Run("Non-Existent File", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)
		_, err := s.Download(context.Background(), "non_existent.txt")
		assert.Error(t, err)
	})
}

func TestOS_Exists(t *testing.T) {
	t.Parallel()

	t.Run("Exists", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)
		require.NoError(t, os.WriteFile(filepath.Join(s.BasePath, "exists.txt"), []byte("exists"), 0644))

		exists, err := s.Exists(context.Background(), "exists.txt")
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Non-Existent File", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)

		exists, err := s.Exists(context.Background(), "non_existent.txt")
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)

		// Create a directory with no read permissions
		noReadDir := filepath.Join(s.BasePath, "no_read")
		require.NoError(t, os.Mkdir(noReadDir, 0000))
		defer func(path string) {
			err := os.RemoveAll(path)
			require.NoError(t, err)
		}(noReadDir)

		_, err := s.Exists(context.Background(), filepath.Join("no_read", "file.txt"))
		assert.Error(t, err)
	})
}

func TestOS_Stat(t *testing.T) {
	t.Parallel()

	t.Run("Stat File", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)

		content := []byte("test content")
		filePath := filepath.Join(s.BasePath, "stat.txt")
		require.NoError(t, os.WriteFile(filePath, content, 0644))

		info, err := s.Stat(context.Background(), "stat.txt")
		require.NoError(t, err)
		assert.Equal(t, int64(len(content)), info.Size)
		assert.False(t, info.IsDir)
		assert.Empty(t, info.ContentType)
		assert.WithinDuration(t, time.Now(), info.LastModified, 2*time.Second)
	})

	t.Run("Stat Dir", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)

		dirPath := filepath.Join(s.BasePath, "statdir")
		require.NoError(t, os.Mkdir(dirPath, 0755))

		info, err := s.Stat(context.Background(), "statdir")
		require.NoError(t, err)
		assert.True(t, info.IsDir)
	})

	t.Run("Non-Existent File", func(t *testing.T) {
		t.Parallel()
		s := setupOSStorage(t)
		_, err := s.Stat(context.Background(), "non_existent.txt")
		assert.Error(t, err)
	})
}
