package storage

import (
	"context"
	"io"
	"time"
)

//go:generate mockgen -package=storagefakes -destination=storagefakes/fake.go . Provider

// Provider defines the interface for file operations on different backends
// (e.g., S3, local filesystem).
type Provider interface {
	// Upload stores the content from the provided io.Reader at the specified path.
	// It returns an error if the upload fails.
	Upload(ctx context.Context, path string, content io.Reader) error

	// Delete removes the file at the specified path. It returns an error if the
	// deletion fails or if the file doesn't exist.
	Delete(ctx context.Context, path string) error

	// List returns a slice of strings representing the names of files and
	// directories under the specified prefix. It returns an error if the listing
	// operation fails.
	List(ctx context.Context, prefix string) ([]string, error)

	// Download retrieves the content of the file at the specified path. It returns
	// an io.ReadCloser for reading the file's content and an error if the download
	// fails. The caller is responsible for closing the returned io.ReadCloser.
	Download(ctx context.Context, path string) (io.ReadCloser, error)

	// Exists checks if a file or directory exists at the specified path. It
	// returns true if the path exists, false if it doesn't, and an error if the
	// check fails.
	Exists(ctx context.Context, path string) (bool, error)

	// Stat retrieves metadata about the file or directory at the specified path.
	// It returns a FileInfo struct containing the metadata and an error if the
	// operation fails.
	Stat(ctx context.Context, path string) (*FileInfo, error)
}

// FileInfo holds metadata about a file, such as size and last modified time.
type FileInfo struct {
	// Size represents the size of the file in bytes.
	Size int64

	// LastModified is the timestamp of when the file was last modified.
	LastModified time.Time

	// IsDir indicates whether the path represents a directory.
	IsDir bool

	// ContentType is the MIME type of the file.
	ContentType string
}
