package files

import (
	"context"
	"path/filepath"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

// DockerIgnore scaffolds .dockerignore files for every app that's defined
// in the app manifest.
func DockerIgnore(_ context.Context, input cmdtools.CommandInput) error {
	for _, app := range input.AppDef().Apps {
		err := input.Generator().Template(filepath.Join(app.Path, ".dockerignore"),
			templates.MustLoadTemplate(".dockerignore"),
			scaffold.WithTracking(manifest.SourceProject()),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
