package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// OS implements the Storage interface for the local file system.
type OS struct {
	BasePath string
}

// NewOSStorage creates a new OS instance.
func NewOSStorage(basePath string) *OS {
	return &OS{BasePath: basePath}
}

func (s *OS) Upload(_ context.Context, path string, content io.Reader) error {
	fullPath := filepath.Join(s.BasePath, path)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", fullPath, err)
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	return err
}

func (s *OS) Delete(_ context.Context, path string) error {
	return os.Remove(filepath.Join(s.BasePath, path))
}

func (s *OS) List(_ context.Context, prefix string) ([]string, error) {
	var files []string
	err := filepath.Walk(filepath.Join(s.BasePath, prefix), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(s.BasePath, path)
			files = append(files, relPath)
		}
		return nil
	})
	return files, err
}

func (s *OS) Download(_ context.Context, path string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(s.BasePath, path))
}

func (s *OS) Exists(_ context.Context, path string) (bool, error) {
	_, err := os.Stat(filepath.Join(s.BasePath, path))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *OS) Stat(_ context.Context, path string) (*FileInfo, error) {
	info, err := os.Stat(filepath.Join(s.BasePath, path))
	if err != nil {
		return nil, err
	}
	return &FileInfo{
		Size:         info.Size(),
		LastModified: info.ModTime(),
		IsDir:        info.IsDir(),
		ContentType:  "", // Not supported for local files
	}, nil
}
