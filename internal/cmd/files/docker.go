package files

import (
	"context"
	"path/filepath"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

// DockerIgnore scaffolds .dockerignore files for every app that's defined
// in the app manifest.
func DockerIgnore(_ context.Context, input cmdtools.CommandInput) error {
	gen := scaffold.New(input.FS, input.Manifest)

	for _, app := range input.AppDef().Apps {
		err := gen.Template(filepath.Join(app.Path, ".dockerignore"),
			templates.MustLoadTemplate(".dockerignore"),
			scaffold.WithTracking("project:root", true),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
