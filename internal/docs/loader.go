package docs

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
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

// HasAppType checks if the definition contains an app of the specified type.
func HasAppType(def *appdef.Definition, appType appdef.AppType) bool {
	if def == nil {
		return false
	}

	for _, app := range def.Apps {
		if app.Type == appType {
			return true
		}
	}

	return false
}

// GetAppsByType returns all apps of the specified type from the definition.
func GetAppsByType(def *appdef.Definition, appType appdef.AppType) []appdef.App {
	if def == nil {
		return nil
	}

	var apps []appdef.App
	for _, app := range def.Apps {
		if app.Type == appType {
			apps = append(apps, app)
		}
	}

	return apps
}
