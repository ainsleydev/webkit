package files

import (
	"context"

	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
)

// Manifest generates a blank file manifest file if it doesn't
// exist yet.
func Manifest(_ context.Context, input cmdtools.CommandInput) error {
	exists, err := afero.Exists(input.FS, manifest.Path)
	if err != nil {
		return err
	}

	// If the manifest exists, it shouldn't be overwritten.
	if exists {
		return nil
	}

	return input.Manifest.Save(input.FS)
}
