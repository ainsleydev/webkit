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

type (
	// ReadmeFrontMatter contains front matter metadata for README templates.
	ReadmeFrontMatter struct {
		Logo *LogoConfig `yaml:"logo,omitempty" json:"logo,omitempty"`
	}

	// LogoConfig contains logo display configuration.
	LogoConfig struct {
		Width  int `yaml:"width,omitempty" json:"width,omitempty"`
		Height int `yaml:"height,omitempty" json:"height,omitempty"`
	}

	// ReadmeContent contains parsed front matter and content.
	ReadmeContent struct {
		Meta    ReadmeFrontMatter
		Content string
	}
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

// loadReadmeContent loads README content and parses front matter if present.
func loadReadmeContent(fs afero.Fs) (*ReadmeContent, error) {
	content, err := afero.ReadFile(fs, filepath.Join(customDocsDir, "README.md"))
	if errors.Is(err, afero.ErrFileNotFound) || errors.Is(err, os.ErrNotExist) {
		return &ReadmeContent{}, nil
	}
	if err != nil {
		return nil, errors.Wrap(err, "reading README.md")
	}

	var meta ReadmeFrontMatter
	rest, err := frontmatter.Parse(bytes.NewReader(content), &meta)
	if err != nil {
		return nil, errors.Wrap(err, "parsing front matter")
	}

	return &ReadmeContent{
		Meta:    meta,
		Content: string(rest),
	}, nil
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
