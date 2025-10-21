package files

import (
	"context"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
)

// PnpmWorkspace scaffolds the pnpm-workspace.yaml file with
// all apps that use NPM/pnpm.
func PnpmWorkspace(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	if len(appDef.Apps) == 0 {
		return nil
	}

	var packages []string
	for _, app := range appDef.Apps {
		if app.ShouldUseNPM() {
			packages = append(packages, app.Path)
		}
	}

	if len(packages) == 0 {
		return nil
	}

	return input.Generator().YAML("pnpm-workspace.yaml", map[string]any{
		"packages": packages,
	}, scaffold.WithTracking(manifest.SourceProject()))
}
