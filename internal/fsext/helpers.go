package fsext

import "github.com/spf13/afero"

// Exists checks if a path exists, ignoring any errors.
//
// Returns true if the path exists, false otherwise.
func Exists(fs afero.Fs, path string) bool {
	exists, _ := afero.Exists(fs, path)
	return exists
}

// DirExists checks if a directory exists, ignoring any errors.
//
// Returns true if the directory exists, false otherwise.
func DirExists(fs afero.Fs, path string) bool {
	exists, _ := afero.DirExists(fs, path)
	return exists
}
