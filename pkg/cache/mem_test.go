package cache

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMem_Ping(t *testing.T) {
	t.Parallel()

	store := NewInMemory(time.Second * 5)
	err := store.Ping(context.Background())
	require.NoError(t, err)
}

func TestMem_SetAndGet(t *testing.T) {
	t.Parallel()

	var (
		key   = "key"
		value = "value"
	)

	tests := map[string]struct {
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

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			store := NewInMemory(time.Second * 5)
			store.Set(context.Background(), test.key, value, test.opts)

			var got string
			err := store.Get(context.Background(), key, &got)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.wantValue, got)
		})
	}

	t.Run("Returns error if value is not a pointer", func(t *testing.T) {
		t.Parallel()
		store := NewInMemory(time.Second * 5)
		store.Set(context.Background(), "key", "value", Options{})

		var value string
		err := store.Get(context.Background(), "key", value)
		assert.Error(t, err)
	})

	t.Run("Works with slices", func(t *testing.T) {
		t.Parallel()
		store := NewInMemory(time.Second * 5)
		s := []string{"1", "2", "3"}
		store.Set(context.Background(), "key", s, Options{})

		var got []string
		err := store.Get(context.Background(), "key", &got)
		require.NoError(t, err)
		require.Equal(t, s, got)
	})
}

func TestMem_Delete(t *testing.T) {
	t.Parallel()

	store := NewInMemory(time.Second * 5)
	store.Set(context.Background(), "key1", "value1", Options{})

	err := store.Delete(context.Background(), "key1")
	require.NoError(t, err)

	err = store.Get(context.Background(), "key1", nil)
	assert.Error(t, err)
}

func TestMem_Flush(t *testing.T) {
	t.Parallel()

	store := NewInMemory(time.Second * 5)
	store.Set(context.Background(), "key1", "value1", Options{})

	store.Flush(context.Background())
	err := store.Get(context.Background(), "key1", nil)
	assert.Error(t, err)
}

func TestMem_Invalidate(t *testing.T) {
	t.Parallel()

	data := map[string]inMemCacheItem{
		"key1": {value: "value1", tags: []string{"tag1", "tag2"}},
		"key2": {value: "value2", tags: []string{"tag1", "tag3"}},
		"key3": {value: "value3", tags: []string{"tag3"}},
	}

	store := &MemCache{
		cache: data,
		mutex: &sync.RWMutex{},
	}

	store.Invalidate(context.Background(), []string{"tag1", "tag2"})

	require.Len(t, store.cache, 1)
	assert.Equal(t, "value3", store.cache["key3"].value)
}

func TestMem_Close(t *testing.T) {
	t.Parallel()

	store := NewInMemory(time.Second * 5)
	err := store.Close()
	require.NoError(t, err)
}
