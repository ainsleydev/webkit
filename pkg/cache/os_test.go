package cache

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOSCache(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tempDir := t.TempDir()
		store, err := NewOSCache(tempDir)
		require.NoError(t, err)
		assert.NotNil(t, store)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()

		// Test with a file instead of a directory
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "not_a_dir")
		err := os.WriteFile(filePath, []byte("not a directory"), 0644)
		require.NoError(t, err)

		_, err = NewOSCache(filePath)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not a directory")
	})
}

func TestOSCache_Ping(t *testing.T) {
	store, err := NewOSCache(t.TempDir())
	require.NoError(t, err)
	got := store.Ping(context.Background())
	assert.NoError(t, got)
}

func TestOSCache_Set(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		key   string
		value string
		opts  Options
	}{
		"Simple": {
			key:   "simple_key",
			value: "simple_value",
			opts:  Options{},
		},
		"With Expiration": {
			key:   "expiring_key",
			value: "expiring_value",
			opts:  Options{Expiration: time.Second * 2},
		},
		"With Tags": {
			key:   "tagged_key",
			value: "tagged_value",
			opts:  Options{Tags: []string{"tag1", "tag2"}},
		},
		"Nested Directory": {
			key:   "nested/directory/key",
			value: "nested_value",
			opts:  Options{},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			tempDir := t.TempDir()
			store, err := NewOSCache(tempDir)
			require.NoError(t, err)

			store.Set(ctx, test.key, test.value, test.opts)

			// Verify the file was created
			_, err = os.Stat(filepath.Join(tempDir, test.key))
			assert.NoError(t, err, "File should exist after Set operation")

			// Verify the value was stored correctly
			var got string
			err = store.Get(ctx, test.key, &got)
			assert.NoError(t, err)
			assert.Equal(t, test.value, got)
		})
	}

	t.Run("Make Dir Error", func(t *testing.T) {
		var buf bytes.Buffer
		slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))

		tempDir := t.TempDir()
		store, err := NewOSCache(tempDir)
		require.NoError(t, err)

		// Make the directory read-only
		require.NoError(t, os.Chmod(tempDir, 0555))

		// Attempt to set a value, which should trigger a directory creation error
		store.Set(context.Background(), "nested/key", "value", Options{})

		assert.Contains(t, buf.String(), "Error creating directory")
	})
}

func TestOSCache_SetError(t *testing.T) {
	var buf bytes.Buffer
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))

	tempDir := t.TempDir()
	store, err := NewOSCache(tempDir)
	require.NoError(t, err)

	// Make the directory read-only
	require.NoError(t, os.Chmod(tempDir, 0555))

	// Attempt to set a value, which should trigger a save error
	store.Set(context.Background(), "key", "value", Options{})

	assert.Contains(t, buf.String(), "Error writing to file")
}

func TestOSCache_Get(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		setup      func(*OSCache, context.Context)
		key        string
		decodeInto func() any
		want       func(got any, err error)
	}{
		"OK String": {
			setup: func(store *OSCache, ctx context.Context) {
				store.Set(ctx, "string_key", "value", Options{})
			},
			key: "string_key",
			decodeInto: func() any {
				return new(string)
			},
			want: func(got any, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "value", *got.(*string))
			},
		},
		"OK Byte Slice": {
			setup: func(store *OSCache, ctx context.Context) {
				store.Set(ctx, "byte_key", []byte("value"), Options{})
			},
			key: "byte_key",
			decodeInto: func() any {
				return new([]byte)
			},
			want: func(got any, err error) {
				assert.NoError(t, err)
				assert.Equal(t, []byte("value"), *got.(*[]byte))
			},
		},
		"Invalid Type": {
			setup: func(store *OSCache, ctx context.Context) {
				store.Set(ctx, "invalid_key", "value", Options{})
			},
			key: "invalid_key",
			decodeInto: func() any {
				return new(int)
			},
			want: func(got any, err error) {
				assert.Error(t, err)
				assert.Equal(t, 0, *got.(*int))
			},
		},
		"Non Existent Key": {
			setup: func(store *OSCache, ctx context.Context) {},
			key:   "non_existent_key",
			decodeInto: func() any {
				return new(string)
			},
			want: func(got any, err error) {
				assert.Error(t, err)
				assert.Equal(t, "", *got.(*string))
			},
		},
		"Expired Key": {
			setup: func(store *OSCache, ctx context.Context) {
				store.Set(ctx, "expired_key", "expired_value", Options{Expiration: time.Nanosecond})
				time.Sleep(time.Millisecond) // Ensure expiration
			},
			key: "expired_key",
			decodeInto: func() any {
				return new(string)
			},
			want: func(got any, err error) {
				assert.Error(t, err)
				assert.Equal(t, "", *got.(*string))
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tempDir := t.TempDir()
			store, err := NewOSCache(tempDir)
			require.NoError(t, err)

			ctx := context.Background()
			test.setup(store, ctx)

			v := test.decodeInto()
			err = store.Get(ctx, test.key, v)

			test.want(v, err)
		})
	}
}

func TestOSCache_Delete(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		setup       func(*OSCache, context.Context)
		deleteKey   string
		expectError bool
	}{
		"OK": {
			setup: func(store *OSCache, ctx context.Context) {
				store.Set(ctx, "key1", "value1", Options{})
			},
			deleteKey:   "key1",
			expectError: false,
		},
		"Error": {
			setup:       func(store *OSCache, ctx context.Context) {},
			deleteKey:   "non_existent_key",
			expectError: true,
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			store, err := NewOSCache(t.TempDir())
			require.NoError(t, err)

			// Setup the test case
			tc.setup(store, ctx)

			// Perform the delete operation
			err = store.Delete(ctx, tc.deleteKey)
			assert.Equal(t, tc.expectError, err != nil)

			// Verify that the key no longer exists in the cache
			var v string
			err = store.Get(ctx, tc.deleteKey, &v)
			assert.Error(t, err, "Key should not exist after deletion")
		})
	}
}

func TestOSCache_Flush(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	store, err := NewOSCache(tempDir)
	require.NoError(t, err)

	ctx := context.Background()
	store.Set(ctx, "key1", "value1", Options{})
	store.Set(ctx, "key2", "value2", Options{})

	store.Flush(ctx)

	var value string
	err = store.Get(ctx, "key1", &value)
	assert.Error(t, err)
	err = store.Get(ctx, "key2", &value)
	assert.Error(t, err)
}

func TestOSCache_Invalidate(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		initialData    map[string]osCacheEntry
		invalidateTags []string
		expectedKeys   []string
	}{
		"SimpleInvalidation": {
			initialData: map[string]osCacheEntry{
				"key1": {Tags: []string{"tag1", "tag2"}},
				"key2": {Tags: []string{"tag1", "tag3"}},
				"key3": {Tags: []string{"tag3"}},
			},
			invalidateTags: []string{"tag1", "tag2"},
			expectedKeys:   []string{"key3"},
		},
		"NoMatchingTags": {
			initialData: map[string]osCacheEntry{
				"key1": {Tags: []string{"tag1"}},
				"key2": {Tags: []string{"tag2"}},
			},
			invalidateTags: []string{"tag3"},
			expectedKeys:   []string{"key1", "key2"},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tempDir := t.TempDir()
			store, err := NewOSCache(tempDir)
			require.NoError(t, err)

			ctx := context.Background()
			for k, v := range test.initialData {
				store.Set(ctx, k, "value", Options{Tags: v.Tags})
			}

			store.Invalidate(ctx, test.invalidateTags)

			for _, key := range test.expectedKeys {
				var value string
				err = store.Get(ctx, key, &value)
				assert.NoError(t, err)
			}

			for k := range test.initialData {
				if !contains(test.expectedKeys, k) {
					var value string
					err = store.Get(ctx, k, &value)
					assert.Error(t, err)
				}
			}
		})
	}
}

func TestOSCache_Persistence(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	// Create and populate the cache
	store1, err := NewOSCache(tempDir)
	require.NoError(t, err)

	ctx := context.Background()
	store1.Set(ctx, "key1", "value1", Options{})
	store1.Set(ctx, "key2", "value2", Options{})

	err = store1.Close()
	require.NoError(t, err)

	// Create a new cache instance and verify the data persists
	store2, err := NewOSCache(tempDir)
	require.NoError(t, err)

	var value1, value2 string
	err = store2.Get(ctx, "key1", &value1)
	require.NoError(t, err)
	assert.Equal(t, "value1", value1)

	err = store2.Get(ctx, "key2", &value2)
	require.NoError(t, err)
	assert.Equal(t, "value2", value2)
}
