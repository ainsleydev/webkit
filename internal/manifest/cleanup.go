package manifest

import (
	"fmt"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/printer"
)

// Cleanup removes files that are no longer needed
func Cleanup(fs afero.Fs, oldManifest, newManifest *Manifest, console *printer.Console) error {
	for path, entry := range oldManifest.Files {
		// Don't try and hash stuff that's managed by the
		// user and not WebKit.
		if entry.ScaffoldMode {
			continue
		}

		// File existed before but doesn't exist now = orphaned
		if _, exists := newManifest.Files[path]; !exists {
			console.Warn(fmt.Sprintf("Removing orphaned: %s", path))

			if err := fs.Remove(path); err != nil {
				return fmt.Errorf("removing %s: %w", path, err)
			}

			console.Success(fmt.Sprintf("Removed: %s", path))
		}
	}

	return nil
}

// DetectDrift checks if files have been manually modified
func DetectDrift(fs afero.Fs, manifest *Manifest) ([]string, error) {
	var drifted []string

	for path, entry := range manifest.Files {
		// File might be deleted or moved.
		data, err := afero.ReadFile(fs, path)
		if err != nil {
			continue
		}

		// Don't try and hash stuff that's managed by the
		// user and not WebKit.
		if entry.ScaffoldMode {
			continue
		}

		currentHash := HashContent(data)
		if currentHash != entry.Hash {
			drifted = append(drifted, path)
		}
	}

	return drifted, nil
}
