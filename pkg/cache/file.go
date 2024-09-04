package cache

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-json"
)

type (
	// FileCache represents a file-based cache store that persists data
	// to a JSON file on disk.
	FileCache struct {
		filePath string
		data     map[string]fileCacheItem
		mtx      *sync.RWMutex
	}
	// fileCacheItem represents a single item in the cache, including its value,
	// expiration time, and associated tags.
	fileCacheItem struct {
		Value      any      `json:"value"`
		Expiration int64    `json:"expiration"`
		Tags       []string `json:"tags"`
	}
)

// NewFileCache creates a new FileCache instance, initializing it with the given
// file path. If the file doesn't exist, it will be created.
func NewFileCache(filePath string) (*FileCache, error) {
	fc := &FileCache{
		filePath: filePath,
		data:     make(map[string]fileCacheItem),
		mtx:      &sync.RWMutex{},
	}
	// Add .json extension if not present
	if !strings.HasSuffix(filePath, ".json") {
		filePath += ".json"
	}
	// Load cache data from file
	if err := fc.load(); err != nil {
		return nil, err
	}
	return fc, nil
}

func (f *FileCache) load() error {
	file, err := os.OpenFile(f.filePath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return errors.New("failed to open cache file: " + err.Error())
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return errors.New("failed to get file stat: " + err.Error())
	}

	if stat.Size() == 0 {
		// File is empty or newly created, initialize with an empty JSON object
		_, err = file.Write([]byte("{}"))
		if err != nil {
			return errors.New("failed to initialize cache file: " + err.Error())
		}
		// Reset file pointer to the beginning
		file.Seek(0, 0) // nolint: errcheck
	}

	return json.NewDecoder(file).Decode(&f.data)
}

func (f *FileCache) save() error {
	file, err := os.Create(f.filePath)
	if err != nil {
		return fmt.Errorf("failed to create cache file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	return encoder.Encode(f.data)
}

func (f *FileCache) Ping(_ context.Context) error {
	return nil // File-based cache is always available
}

func (f *FileCache) Get(_ context.Context, key string, v any) error {
	f.mtx.RLock()
	defer f.mtx.RUnlock()

	item, ok := f.data[key]
	if !ok {
		return errors.New("key not found")
	}

	if item.Expiration != 0 && item.Expiration < time.Now().UnixNano() {
		delete(f.data, key)
		return errors.New("key expired: " + key)
	}

	b, err := json.Marshal(item.Value)
	if err != nil {
		return errors.New("failed to marshal cached value: " + err.Error())
	}

	return json.Unmarshal(b, v)
}

func (f *FileCache) Set(_ context.Context, key string, value any, options Options) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	var expiration int64
	if options.Expiration > 0 {
		expiration = time.Now().Add(options.Expiration).UnixNano()
	}

	f.data[key] = fileCacheItem{
		Value:      value,
		Expiration: expiration,
		Tags:       options.Tags,
	}

	if err := f.save(); err != nil {
		slog.Error("Error saving key: " + err.Error())
	}
}

func (f *FileCache) Delete(ctx context.Context, key string) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	delete(f.data, key)
	return f.save()
}

func (f *FileCache) Invalidate(_ context.Context, tags []string) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	for key, item := range f.data {
		for _, tag := range tags {
			if contains(item.Tags, tag) {
				delete(f.data, key)
				break
			}
		}
	}

	if err := f.save(); err != nil {
		slog.Error("Error saving cache to file after invalidation: " + err.Error())
	}
}

func (f *FileCache) Flush(_ context.Context) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	f.data = make(map[string]fileCacheItem)
	if err := f.save(); err != nil {
		slog.Error("Error saving cache to file after flush: " + err.Error())
	}
}

func (f *FileCache) Close() error {
	return f.save()
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
