package docs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/frontmatter"
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

// parseContentWithFrontMatter parses a markdown file with optional YAML front matter.
// The meta parameter should be a pointer to the struct where front matter will be unmarshalled.
// Returns the content without front matter, or empty string if file doesn't exist.
func parseContentWithFrontMatter(fs afero.Fs, filePath string, meta any) (string, error) {
	content, err := afero.ReadFile(fs, filePath)
	if errors.Is(err, afero.ErrFileNotFound) || errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", errors.Wrapf(err, "reading %s", filepath.Base(filePath))
	}

	rest, err := frontmatter.Parse(bytes.NewReader(content), meta)
	if err != nil {
		return "", errors.Wrap(err, "parsing front matter")
	}

	return string(rest), nil
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
