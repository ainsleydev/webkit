package docs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const (
	// customDocsDir is the directory for custom documentation content.
	customDocsDir = "docs"
)

// loadCustomContent loads custom a documentation file content from the file name.
func loadCustomContent(fs afero.Fs, fileName string) (string, error) {
	return readFile(fs, filepath.Join(customDocsDir, fileName))
}

// mustLoadCustomContent returns an empty string if it didn't exist.
func mustLoadCustomContent(fs afero.Fs, fileName string) string {
	got, err := readFile(fs, filepath.Join(customDocsDir, fileName))
	if err != nil {
		return ""
	}
	return got
}

func readFile(fs afero.Fs, path string) (string, error) {
	content, err := afero.ReadFile(fs, path)
	if errors.Is(err, afero.ErrFileNotFound) || errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("doc template does not exist: %s", path)
	} else if err != nil {
		return "", errors.Wrap(err, "reading file")
	}
	return string(content), nil
}
