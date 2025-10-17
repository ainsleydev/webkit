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
	gen := scaffold.New(input.FS)

	for _, app := range input.AppDef().Apps {
		path := filepath.Join(app.Path, ".dockerignore")
		tpl := templates.MustLoadTemplate(".dockerignore")

		err := gen.Template(path, tpl, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
