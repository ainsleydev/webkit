package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// AppName is the name of the application
	AppName = "webkit"

	// DirName is the config directory name under user's home
	DirName = ".config"
)

// Read reads a file from the WebKit config directory.
func Read(filename string) ([]byte, error) {
	path, err := Path(filename)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(path)
}

// Write writes a file to the WebKit config directory.
// Creates the config directory if it doesn't exist.
func Write(filename string, data []byte, perm os.FileMode) error {
	if err := ensureDir(); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	path, err := Path(filename)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, perm)
}

// Dir returns the WebKit configuration directory path.
// Returns: ~/.config/webkit
func Dir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(home, DirName, AppName), nil
}

// Path returns the full path to a file in the WebKit config directory.
// Example: config.Path("age.key") returns ~/.config/webkit/age.key
func Path(filename string) (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, filename), nil
}

// ensureDir creates the WebKit config directory if it doesn't exist.
func ensureDir() error {
	dir, err := Dir()
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0755)
}
