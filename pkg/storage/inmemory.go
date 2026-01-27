// Package storage provides in-memory storage implementations for testing.
package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"
)

// InMemory implements the Provider interface using an in-memory map for testing purposes.
// All operations are thread-safe using a read-write mutex. Writes are fully isolated
// and reads see a consistent snapshot. All operations respect context cancellation.
type InMemory struct {
	mu   sync.RWMutex
	data map[string]*fileEntry
}

// fileEntry holds file data and metadata.
type fileEntry struct {
	data         []byte
	lastModified time.Time
}

// NewInMemory creates a new in-memory storage provider.
func NewInMemory() *InMemory {
	return &InMemory{
		data: make(map[string]*fileEntry),
	}
}

// Upload stores the content from the provided io.Reader at the specified path.
func (m *InMemory) Upload(ctx context.Context, path string, content io.Reader) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	data, err := io.ReadAll(content)
	if err != nil {
		return fmt.Errorf("reading content: %w", err)
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[path] = &fileEntry{
		data:         data,
		lastModified: time.Now().UTC(),
	}
	return nil
}

// Delete removes the file at the specified path.
// Returns an error if the file doesn't exist.
func (m *InMemory) Delete(ctx context.Context, path string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.data[path]; !exists {
		return fmt.Errorf("file not found: %s", path)
	}

	delete(m.data, path)
	return nil
}

// List returns a slice of strings representing the keys that match the specified prefix.
func (m *InMemory) List(ctx context.Context, prefix string) ([]string, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var keys []string
	for key := range m.data {
		if strings.HasPrefix(key, prefix) {
			keys = append(keys, key)
		}
	}

	return keys, nil
}

// Download retrieves the content of the file at the specified path.
func (m *InMemory) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, exists := m.data[path]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	return io.NopCloser(bytes.NewReader(entry.data)), nil
}

// Exists checks if a file exists at the specified path.
func (m *InMemory) Exists(ctx context.Context, path string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.data[path]
	return exists, nil
}

// Stat retrieves metadata about the file at the specified path.
func (m *InMemory) Stat(ctx context.Context, path string) (*FileInfo, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, exists := m.data[path]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	return &FileInfo{
		Size:         int64(len(entry.data)),
		LastModified: entry.lastModified,
		IsDir:        false,
		ContentType:  "",
	}, nil
}

// Reset clears all stored data. This method is only available on the in-memory
// implementation and is primarily useful for resetting state between tests without
// creating new instances. Callers must ensure no concurrent operations are in progress,
// as this method acquires a write lock and will block until all reads complete.
func (m *InMemory) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = make(map[string]*fileEntry)
}
