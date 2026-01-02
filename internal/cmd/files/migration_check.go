package files

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/scaffold"
	"github.com/ainsleydev/webkit/internal/state/manifest"
	"github.com/ainsleydev/webkit/internal/templates"
)

// MigrationCheckScript scaffolds a dependency check script for Payload CMS apps.
// This script ensures dependencies are up-to-date before running migrations.
func MigrationCheckScript(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	// Filter to only Payload apps.
	payloadApps := appDef.GetAppsByType(appdef.AppTypePayload)

	for _, app := range payloadApps {
		scriptPath := filepath.Join(app.Path, "scripts", "check-deps.cjs")

		// Scaffold the check script from template.
		err := input.Generator().CopyFromEmbed(
			templates.Embed,
			"scripts/check-deps.cjs",
			scriptPath,
			scaffold.WithTracking(manifest.SourceApp(app.Name)),
			scaffold.WithScaffoldMode(),
		)
		if err != nil {
			return errors.Wrap(err, "creating migration check script")
		}
	}

	return nil
}
