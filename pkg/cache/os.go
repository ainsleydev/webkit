package cache

import (
	"context"
	"encoding/json"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type (
	// OSCache represents a file system-based cache where each key is stored as
	// a separate file.
	// It uses a separate index file to store expiration times and tags named
	// index.json which is stored in the base directory.
	OSCache struct {
		basePath  string
		indexPath string
		mtx       sync.RWMutex
		index     map[string]*cacheEntry
	}
	// cacheEntry represents a single item in the cache, including its expiration
	// time and associated tags.
	cacheEntry struct {
		Expiration time.Time `json:"expiration"`
		Tags       []string  `json:"tags"`
	}
)

const osIndexFileName = "index.json"

// NewOSCache creates and initializes a new OSCache instance.
// It creates the base directory if it doesn't exist and loads the index.
func NewOSCache(basePath string) (*OSCache, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}

	cache := &OSCache{
		basePath:  basePath,
		indexPath: filepath.Join(basePath, osIndexFileName),
		index:     make(map[string]*cacheEntry),
	}

	if err := cache.loadIndex(); err != nil {
		return nil, err
	}

	return cache, nil
}

func (o *OSCache) Ping(_ context.Context) error {
	return nil // Always available
}

func (o *OSCache) Get(ctx context.Context, key string, v any) error {
	o.mtx.RLock()
	entry, exists := o.index[key]
	o.mtx.RUnlock()

	if !exists {
		return errors.New("key not found")
	}

	if !entry.Expiration.IsZero() && time.Now().After(entry.Expiration) {
		err := o.Delete(ctx, key)
		if err != nil {
			slog.Error("Error deleting expired key: " + err.Error())
		}
		return errors.New("key expired")
	}

	filePath := o.getFilePath(key)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func (o *OSCache) Set(_ context.Context, key string, value any, options Options) {
	filePath := o.getFilePath(key)

	var data []byte
	var err error

	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		// For other types, we'll still need to marshal to JSON
		data, err = json.Marshal(value)
		if err != nil {
			slog.Error("Error marshaling value: " + err.Error())
			return
		}
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		slog.Error("Error writing to file: " + err.Error())
		return
	}

	var expiration time.Time
	if options.Expiration > 0 {
		expiration = time.Now().Add(options.Expiration)
	}

	o.mtx.Lock()
	o.index[key] = &cacheEntry{
		Expiration: expiration,
		Tags:       options.Tags,
	}
	o.mtx.Unlock()

	if err = o.saveIndex(); err != nil {
		slog.Error("Error saving index: " + err.Error())
	}
}

func (o *OSCache) Delete(_ context.Context, key string) error {
	o.mtx.Lock()
	delete(o.index, key)
	o.mtx.Unlock()

	if err := o.saveIndex(); err != nil {
		return err
	}

	return os.Remove(o.getFilePath(key))
}

func (o *OSCache) Invalidate(_ context.Context, tags []string) {
	var keysToDelete []string

	o.mtx.RLock()
	for key, entry := range o.index {
		for _, tag := range tags {
			if contains(entry.Tags, tag) {
				keysToDelete = append(keysToDelete, key)
				break
			}
		}
	}
	o.mtx.RUnlock()

	o.mtx.Lock()
	for _, key := range keysToDelete {
		delete(o.index, key)
		_ = os.Remove(o.getFilePath(key))
	}
	o.mtx.Unlock()

	if err := o.saveIndex(); err != nil {
		slog.Error("Error saving index after invalidation: " + err.Error())
	}
}

func (o *OSCache) Flush(_ context.Context) {
	o.mtx.Lock()
	o.index = make(map[string]*cacheEntry)
	o.mtx.Unlock()

	err := filepath.Walk(o.basePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && path != o.indexPath {
			return os.Remove(path)
		}
		return nil
	})

	if err != nil {
		slog.Error("Error flushing cache: " + err.Error())
	}

	if err = o.saveIndex(); err != nil {
		slog.Error("Error saving index after flush: " + err.Error())
	}
}

func (o *OSCache) Close() error {
	return o.saveIndex()
}

func (o *OSCache) loadIndex() error {
	o.mtx.Lock()
	defer o.mtx.Unlock()

	data, err := os.ReadFile(o.indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No index file yet, start with empty index
		}
		return err
	}

	return json.Unmarshal(data, &o.index)
}

func (o *OSCache) saveIndex() error {
	o.mtx.RLock()
	indexCopy := make(map[string]*cacheEntry)
	for k, v := range o.index {
		indexCopy[k] = v
	}
	o.mtx.RUnlock()

	data, err := json.MarshalIndent(indexCopy, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(o.indexPath, data, 0644)
}

func (o *OSCache) getFilePath(key string) string {
	return filepath.Join(o.basePath, key)
}