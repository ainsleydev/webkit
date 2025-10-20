package cache

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFile_Ping(t *testing.T) {
	tempDir := t.TempDir()
	store, err := NewFileCache(filepath.Join(tempDir, "cache.json"))
	require.NoError(t, err)

	err = store.Ping(context.Background())
	require.NoError(t, err)
}

func TestFile_SetAndGet(t *testing.T) {
	var (
		key   = "key"
		value = "value"
	)

	tt := map[string]struct {
		key       string
		opts      Options
		wantValue any
		wantErr   bool
	}{
		"SetAndGetWithExpiration": {
			key:       key,
			opts:      Options{Expiration: time.Second * 2},
			wantValue: value,
			wantErr:   false,
		},
		"Not found": {
			key:       "wrong",
			wantValue: "",
			wantErr:   true,
		},
		"Expired": {
			key:       key,
			opts:      Options{Expiration: time.Nanosecond},
			wantValue: "",
			wantErr:   true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			tempDir := t.TempDir()
			store, err := NewFileCache(filepath.Join(tempDir, "cache.json"))
			require.NoError(t, err)

			store.Set(context.Background(), test.key, value, test.opts)

			var got string
			err = store.Get(context.Background(), key, &got)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.wantValue, got)
		})
	}

	t.Run("Returns error if value is not a pointer", func(t *testing.T) {
		tempDir := t.TempDir()
		store, err := NewFileCache(filepath.Join(tempDir, "cache.json"))
		require.NoError(t, err)

		store.Set(context.Background(), "key", "value", Options{})

		var value string
		err = store.Get(context.Background(), "key", value)
		assert.Error(t, err)
	})

	t.Run("Works with slices", func(t *testing.T) {
		tempDir := t.TempDir()
		store, err := NewFileCache(filepath.Join(tempDir, "cache.json"))
		require.NoError(t, err)

		s := []string{"1", "2", "3"}
		store.Set(context.Background(), "key", s, Options{})

		var got []string
		err = store.Get(context.Background(), "key", &got)
		require.NoError(t, err)
		require.Equal(t, s, got)
	})
}

func TestFile_Delete(t *testing.T) {
	tempDir := t.TempDir()
	store, err := NewFileCache(filepath.Join(tempDir, "cache.json"))
	require.NoError(t, err)

	store.Set(context.Background(), "key1", "value1", Options{})

	err = store.Delete(context.Background(), "key1")
	require.NoError(t, err)

	err = store.Get(context.Background(), "key1", nil)
	assert.Error(t, err)
}

func TestFile_Flush(t *testing.T) {
	tempDir := t.TempDir()
	store, err := NewFileCache(filepath.Join(tempDir, "cache.json"))
	require.NoError(t, err)

	store.Set(context.Background(), "key1", "value1", Options{})

	store.Flush(context.Background())
	err = store.Get(context.Background(), "key1", nil)
	assert.Error(t, err)
}

func TestFile_Invalidate(t *testing.T) {
	tt := map[string]struct {
		initialData    map[string]fileCacheItem
		invalidateTags []string
		expectedKeys   []string
	}{
		"Simple": {
			initialData: map[string]fileCacheItem{
				"key1": {Value: "value1", Tags: []string{"tag1", "tag2"}},
				"key2": {Value: "value2", Tags: []string{"tag1", "tag3"}},
				"key3": {Value: "value3", Tags: []string{"tag3"}},
			},
			invalidateTags: []string{"tag1", "tag2"},
			expectedKeys:   []string{"key3"},
		},
		"Error Case": {
			initialData: map[string]fileCacheItem{
				"key1": {Value: "value1", Tags: []string{"tag1"}},
				"key2": {Value: "value2", Tags: []string{"tag1"}},
				"key3": {Value: "value3", Tags: []string{"tag2"}},
			},
			invalidateTags: []string{"tag1"},
			expectedKeys:   []string{"key3"},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tempDir := t.TempDir()
			store, err := NewFileCache(filepath.Join(tempDir, "cache.json"))
			require.NoError(t, err)

			for k, v := range test.initialData {
				store.Set(context.Background(), k, v.Value, Options{Tags: v.Tags})
			}

			store.Invalidate(context.Background(), test.invalidateTags)

			for _, key := range test.expectedKeys {
				var value string
				err := store.Get(context.Background(), key, &value)
				assert.NoError(t, err)
			}

			nonExpectedKeys := make([]string, 0)
			for k := range test.initialData {
				if !contains(test.expectedKeys, k) {
					nonExpectedKeys = append(nonExpectedKeys, k)
				}
			}

			for _, key := range nonExpectedKeys {
				var value string
				err := store.Get(context.Background(), key, &value)
				assert.Error(t, err)
			}
		})
	}
}

func TestFile_Close(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	store, err := NewFileCache(filepath.Join(tempDir, "cache.json"))
	require.NoError(t, err)

	err = store.Close()
	require.NoError(t, err)
}

func TestFile_Persistence(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "cache")

	// Create and populate the cache
	store1, err := NewFileCache(filePath)
	require.NoError(t, err)

	store1.Set(context.Background(), "key1", "value1", Options{})
	store1.Set(context.Background(), "key2", "value2", Options{})

	err = store1.Close()
	require.NoError(t, err)

	// Read the file contents for debugging
	fileContents, err := os.ReadFile(filePath)
	require.NoError(t, err)
	t.Logf("File contents after store1: %s", string(fileContents))

	// Create a new cache instance and verify the data persists
	store2, err := NewFileCache(filePath)
	require.NoError(t, err)

	// Read the file contents again for debugging
	fileContents, err = os.ReadFile(filePath)
	require.NoError(t, err)
	t.Logf("File contents after store2 init: %s", string(fileContents))

	var value1, value2 string
	err = store2.Get(context.Background(), "key1", &value1)
	if err != nil {
		t.Logf("Error getting key1: %v", err)
		t.Logf("store2 data: %+v", store2.data)
	}
	require.NoError(t, err)
	assert.Equal(t, "value1", value1)

	err = store2.Get(context.Background(), "key2", &value2)
	if err != nil {
		t.Logf("Error getting key2: %v", err)
		t.Logf("store2 data: %+v", store2.data)
	}
	require.NoError(t, err)
	assert.Equal(t, "value2", value2)
}

func TestFile_LoadError(t *testing.T) {
	// Test with a directory instead of a file
	tempDir := t.TempDir()
	_, err := NewFileCache(tempDir)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open cache file")
}

func TestFile_SaveError(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "cache.json")
	store, err := NewFileCache(filePath)
	require.NoError(t, err)

	// Make the file read-only
	require.NoError(t, os.Chmod(filePath, 0o444))

	// Attempt to set a value, which should trigger a save error
	store.Set(context.Background(), "key", "value", Options{})
}

func TestFile_GetMarshalError(t *testing.T) {
	tempDir := t.TempDir()
	store, err := NewFileCache(filepath.Join(tempDir, "cache.json"))
	require.NoError(t, err)

	// Set a value that can't be marshaled (e.g., a channel)
	store.Set(context.Background(), "key", make(chan int), Options{})

	var result chan int
	err = store.Get(context.Background(), "key", &result)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal cached value")
}
