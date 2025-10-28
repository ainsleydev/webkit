package docs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const (
	// genDocsDir is the directory for generated documentation files.
	genDocsDir = "internal/gen/docs"

	// customDocsDir is the directory for custom documentation content.
	customDocsDir = "docs"

	// agentsFilename is the name of the custom agents file.
	agentsFilename = "AGENTS.md"
)

// LoadGenFile loads a generated documentation file from internal/gen/docs/.
func LoadGenFile(fs afero.Fs, filename string) (string, error) {
	path := filepath.Join(genDocsDir, filename)

	exists, err := afero.Exists(fs, path)
	if err != nil {
		return "", errors.Wrap(err, "checking file existence")
	}

	if !exists {
		return "", nil
	}

	content, err := afero.ReadFile(fs, path)
	if err != nil {
		return "", errors.Wrap(err, "reading generated file")
	}

	return string(content), nil
}

// MustLoadGenFile loads a generated documentation file and exits if it fails.
func MustLoadGenFile(fs afero.Fs, filename string) string {
	content, err := LoadGenFile(fs, filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading generated file %s: %v\n", filename, err)
		os.Exit(1)
	}

	if content == "" {
		fmt.Fprintf(os.Stderr, "Error: generated file %s does not exist or is empty\n", filename)
		fmt.Fprintf(os.Stderr, "Please run 'go run cmd/docs/main.go' to generate documentation files\n")
		os.Exit(1)
	}

	return content
}

// LoadCustomContent loads custom documentation content from docs/AGENTS.md.
func LoadCustomContent(fs afero.Fs) (string, error) {
	path := filepath.Join(customDocsDir, agentsFilename)

	exists, err := afero.Exists(fs, path)
	if err != nil {
		return "", errors.Wrap(err, "checking file existence")
	}

	if !exists {
		return "", nil
	}

	content, err := afero.ReadFile(fs, path)
	if err != nil {
		return "", errors.Wrap(err, "reading custom content")
	}

	return string(content), nil
}
