package files

import (
	"context"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/templates"
)

// TurboJSON adds a boilerplate turbo.json file at the root of
// the project.
func TurboJSON(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

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

	return input.Generator().Template("./turbo.json",
		templates.MustLoadTemplate("turbo.json"),
		nil,
		scaffold.WithTracking(manifest.SourceProject()),
		scaffold.WithoutNotice(),
	)
}
