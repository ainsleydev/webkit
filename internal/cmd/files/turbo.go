package files

import (
	"context"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

func CreateTurboJson(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()
	gen := scaffold.New(input.FS, input.Manifest)

	if len(appDef.Apps) == 0 {
		return nil
	}

	var packages []string
	for _, app := range input.AppDef().Apps {
		if app.ShouldUseNPM() {
			packages = append(packages, app.Path)
		}
	}

	if len(packages) == 0 {
		return nil
	}

	return gen.Template("./turbo.json",
		templates.MustLoadTemplate("turbo.json"),
		nil,
		scaffold.WithTracking(manifest.SourceProject()),
		scaffold.WithoutNotice(),
	)
}
