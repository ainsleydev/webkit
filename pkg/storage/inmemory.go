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
// All operations are thread-safe and respect context cancellation.
type InMemory struct {
	mu   sync.RWMutex
	data map[string][]byte
}

// NewInMemory creates a new in-memory storage provider.
func NewInMemory() *InMemory {
	return &InMemory{
		data: make(map[string][]byte),
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

	m.data[path] = data
	return nil
}

// Delete removes the file at the specified path.
// Returns nil if the file doesn't exist (matching S3 behaviour).
func (m *InMemory) Delete(ctx context.Context, path string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

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

	data, exists := m.data[path]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	return io.NopCloser(bytes.NewReader(data)), nil
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

	data, exists := m.data[path]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", path)
	}

	return &FileInfo{
		Size:         int64(len(data)),
		LastModified: time.Now(),
		IsDir:        false,
		ContentType:  "",
	}, nil
}
