package storage

import (
	"bytes"
	"context"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInMemory(t *testing.T) {
	t.Parallel()

	store := NewInMemory()
	assert.NotNil(t, store)
	assert.NotNil(t, store.data)
	assert.Equal(t, 0, len(store.data))
}

func TestInMemory_Upload(t *testing.T) {
	t.Parallel()

	t.Run("Successful upload", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		content := bytes.NewReader([]byte("test content"))
		err := store.Upload(context.Background(), "test.txt", content)

		require.NoError(t, err)
		assert.Equal(t, []byte("test content"), store.data["test.txt"].data)
		assert.False(t, store.data["test.txt"].lastModified.IsZero())
	})

	t.Run("Upload overwrites existing file", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		content1 := bytes.NewReader([]byte("first content"))
		err := store.Upload(context.Background(), "test.txt", content1)
		require.NoError(t, err)

		firstModified := store.data["test.txt"].lastModified

		time.Sleep(time.Millisecond)

		content2 := bytes.NewReader([]byte("second content"))
		err = store.Upload(context.Background(), "test.txt", content2)
		require.NoError(t, err)

		assert.Equal(t, []byte("second content"), store.data["test.txt"].data)
		assert.True(t, store.data["test.txt"].lastModified.After(firstModified))
	})

	t.Run("Upload with cancelled context", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		content := bytes.NewReader([]byte("test content"))
		err := store.Upload(ctx, "test.txt", content)

		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})

	t.Run("Upload with read error", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		reader := &errorReader{err: io.ErrUnexpectedEOF}
		err := store.Upload(context.Background(), "test.txt", reader)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reading content")
	})
}

func TestInMemory_Delete(t *testing.T) {
	t.Parallel()

	t.Run("Successful delete", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		store.data["test.txt"] = &fileEntry{
			data:         []byte("test content"),
			lastModified: time.Now(),
		}

		err := store.Delete(context.Background(), "test.txt")
		require.NoError(t, err)
		assert.NotContains(t, store.data, "test.txt")
	})

	t.Run("Delete non-existent file", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		err := store.Delete(context.Background(), "non-existent.txt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file not found")
	})

	t.Run("Delete with cancelled context", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := store.Delete(ctx, "test.txt")
		assert.Error(t, err)
		assert.Equal(t, context.Canceled, err)
	})
}

func TestInMemory_List(t *testing.T) {
	t.Parallel()

	t.Run("List with prefix", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		now := time.Now()
		store.data["files/test1.txt"] = &fileEntry{data: []byte("content1"), lastModified: now}
		store.data["files/test2.txt"] = &fileEntry{data: []byte("content2"), lastModified: now}
		store.data["other/test3.txt"] = &fileEntry{data: []byte("content3"), lastModified: now}

		keys, err := store.List(context.Background(), "files/")
		require.NoError(t, err)
		assert.Len(t, keys, 2)
		assert.Contains(t, keys, "files/test1.txt")
		assert.Contains(t, keys, "files/test2.txt")
		assert.NotContains(t, keys, "other/test3.txt")
	})

	t.Run("List all files", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		now := time.Now()
		store.data["test1.txt"] = &fileEntry{data: []byte("content1"), lastModified: now}
		store.data["test2.txt"] = &fileEntry{data: []byte("content2"), lastModified: now}

		keys, err := store.List(context.Background(), "")
		require.NoError(t, err)
		assert.Len(t, keys, 2)
	})

	t.Run("List empty store", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		keys, err := store.List(context.Background(), "")
		require.NoError(t, err)
		assert.Len(t, keys, 0)
	})

	t.Run("List with cancelled context", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		keys, err := store.List(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, keys)
		assert.Equal(t, context.Canceled, err)
	})
}

func TestInMemory_Download(t *testing.T) {
	t.Parallel()

	t.Run("Successful download", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		store.data["test.txt"] = &fileEntry{
			data:         []byte("test content"),
			lastModified: time.Now(),
		}

		reader, err := store.Download(context.Background(), "test.txt")
		require.NoError(t, err)
		defer reader.Close()

		content, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, []byte("test content"), content)
	})

	t.Run("Download non-existent file", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		reader, err := store.Download(context.Background(), "non-existent.txt")

		assert.Error(t, err)
		assert.Nil(t, reader)
		assert.Contains(t, err.Error(), "file not found")
	})

	t.Run("Download with cancelled context", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		reader, err := store.Download(ctx, "test.txt")
		assert.Error(t, err)
		assert.Nil(t, reader)
		assert.Equal(t, context.Canceled, err)
	})
}

func TestInMemory_Exists(t *testing.T) {
	t.Parallel()

	t.Run("File exists", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		store.data["test.txt"] = &fileEntry{
			data:         []byte("test content"),
			lastModified: time.Now(),
		}

		exists, err := store.Exists(context.Background(), "test.txt")
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("File does not exist", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		exists, err := store.Exists(context.Background(), "non-existent.txt")
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Exists with cancelled context", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		exists, err := store.Exists(ctx, "test.txt")
		assert.Error(t, err)
		assert.False(t, exists)
		assert.Equal(t, context.Canceled, err)
	})
}

func TestInMemory_Stat(t *testing.T) {
	t.Parallel()

	t.Run("Stat existing file", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		now := time.Now()
		store.data["test.txt"] = &fileEntry{
			data:         []byte("test content"),
			lastModified: now,
		}

		info, err := store.Stat(context.Background(), "test.txt")
		require.NoError(t, err)
		assert.Equal(t, int64(12), info.Size)
		assert.False(t, info.IsDir)
		assert.Equal(t, "", info.ContentType)
		assert.Equal(t, now, info.LastModified)
	})

	t.Run("Stat non-existent file", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		info, err := store.Stat(context.Background(), "non-existent.txt")

		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Contains(t, err.Error(), "file not found")
	})

	t.Run("Stat with cancelled context", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		info, err := store.Stat(ctx, "test.txt")
		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, context.Canceled, err)
	})
}

func TestInMemory_ThreadSafety(t *testing.T) {
	t.Parallel()

	t.Run("Concurrent uploads", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		var wg sync.WaitGroup
		iterations := 100

		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				content := bytes.NewReader([]byte("content"))
				_ = store.Upload(context.Background(), "test.txt", content)
			}()
		}

		wg.Wait()
		assert.NotNil(t, store.data["test.txt"])
	})

	t.Run("Concurrent reads and writes", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		store.data["test.txt"] = &fileEntry{
			data:         []byte("initial content"),
			lastModified: time.Now(),
		}

		var wg sync.WaitGroup
		iterations := 50

		for i := 0; i < iterations; i++ {
			wg.Add(2)

			go func() {
				defer wg.Done()
				content := bytes.NewReader([]byte("new content"))
				_ = store.Upload(context.Background(), "test.txt", content)
			}()

			go func() {
				defer wg.Done()
				_, _ = store.Download(context.Background(), "test.txt")
			}()
		}

		wg.Wait()
	})

	t.Run("Concurrent list operations", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		now := time.Now()
		store.data["test1.txt"] = &fileEntry{data: []byte("content1"), lastModified: now}
		store.data["test2.txt"] = &fileEntry{data: []byte("content2"), lastModified: now}

		var wg sync.WaitGroup
		iterations := 50

		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, _ = store.List(context.Background(), "")
			}()
		}

		wg.Wait()
	})
}

func TestInMemory_IntegrationScenario(t *testing.T) {
	t.Parallel()

	t.Run("Full workflow", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		ctx := context.Background()

		exists, err := store.Exists(ctx, "test.txt")
		require.NoError(t, err)
		assert.False(t, exists)

		content := bytes.NewReader([]byte("test content"))
		err = store.Upload(ctx, "test.txt", content)
		require.NoError(t, err)

		exists, err = store.Exists(ctx, "test.txt")
		require.NoError(t, err)
		assert.True(t, exists)

		info, err := store.Stat(ctx, "test.txt")
		require.NoError(t, err)
		assert.Equal(t, int64(12), info.Size)

		reader, err := store.Download(ctx, "test.txt")
		require.NoError(t, err)
		downloadedContent, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, []byte("test content"), downloadedContent)
		reader.Close()

		keys, err := store.List(ctx, "")
		require.NoError(t, err)
		assert.Len(t, keys, 1)
		assert.Contains(t, keys, "test.txt")

		err = store.Delete(ctx, "test.txt")
		require.NoError(t, err)

		exists, err = store.Exists(ctx, "test.txt")
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestInMemory_ContextTimeout(t *testing.T) {
	t.Parallel()

	t.Run("Upload with timeout", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		time.Sleep(2 * time.Nanosecond)
		content := bytes.NewReader([]byte("test content"))
		err := store.Upload(ctx, "test.txt", content)
		assert.Error(t, err)
	})

	t.Run("Download with timeout", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
		defer cancel()

		time.Sleep(2 * time.Nanosecond)
		_, err := store.Download(ctx, "test.txt")
		assert.Error(t, err)
	})
}

func TestInMemory_Reset(t *testing.T) {
	t.Parallel()

	t.Run("Reset clears all data", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		ctx := context.Background()

		content := bytes.NewReader([]byte("test content"))
		err := store.Upload(ctx, "test.txt", content)
		require.NoError(t, err)

		exists, err := store.Exists(ctx, "test.txt")
		require.NoError(t, err)
		assert.True(t, exists)

		store.Reset()

		exists, err = store.Exists(ctx, "test.txt")
		require.NoError(t, err)
		assert.False(t, exists)

		keys, err := store.List(ctx, "")
		require.NoError(t, err)
		assert.Len(t, keys, 0)
	})

	t.Run("Reset allows reuse", func(t *testing.T) {
		t.Parallel()

		store := NewInMemory()
		ctx := context.Background()

		content1 := bytes.NewReader([]byte("first"))
		err := store.Upload(ctx, "file1.txt", content1)
		require.NoError(t, err)

		store.Reset()

		content2 := bytes.NewReader([]byte("second"))
		err = store.Upload(ctx, "file2.txt", content2)
		require.NoError(t, err)

		keys, err := store.List(ctx, "")
		require.NoError(t, err)
		assert.Len(t, keys, 1)
		assert.Contains(t, keys, "file2.txt")
		assert.NotContains(t, keys, "file1.txt")
	})
}

// errorReader is a helper type that always returns an error when read.
type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}
