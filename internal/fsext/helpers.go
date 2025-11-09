package fsext

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

// Exists checks if a path exists, ignoring any errors.
//
// Returns true if the path exists, false otherwise.
func Exists(fs afero.Fs, path string) bool {
	exists, _ := afero.Exists(fs, path) //nolint:errcheck
	return exists
}

// DirExists checks if a directory exists, ignoring any errors.
//
// Returns true if the directory exists, false otherwise.
func DirExists(fs afero.Fs, path string) bool {
	exists, _ := afero.DirExists(fs, path) //nolint:errcheck
	return exists
}

// EnsureDir ensures the parent directory of the given file path exists.
//
// Returns an error if the directory cannot be created.
func EnsureDir(fs afero.Fs, filePath string) error {
	return fs.MkdirAll(filepath.Dir(filePath), os.ModePerm)
}
