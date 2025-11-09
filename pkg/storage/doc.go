// Package storage provides a unified interface for file operations across different backends.
//
// It supports multiple storage providers including S3 and local filesystems,
// allowing applications to switch between storage backends without code changes.
// All operations are context-aware and support standard file operations like
// upload, download, delete, list, and stat.
package storage
