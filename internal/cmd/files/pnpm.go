package files

import (
	"context"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/state/manifest"
)

// PnpmWorkspace scaffolds the pnpm-workspace.yaml file with
// all apps and utilities that use NPM/pnpm.
func PnpmWorkspace(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	if len(appDef.Apps) == 0 && len(appDef.Utilities) == 0 {
		return nil
	}

	var packages []string
	for _, app := range appDef.Apps {
		if app.ShouldUseNPM() {
			packages = append(packages, app.Path)
		}
	}
	for _, util := range appDef.Utilities {
		if util.ShouldUseNPM() {
			packages = append(packages, util.Path)
		}
	}

	if len(packages) == 0 {
		return nil
	}

	return input.Generator().YAML("pnpm-workspace.yaml", map[string]any{
		"packages": packages,
	}, scaffold.WithTracking(manifest.SourceProject()))
}
