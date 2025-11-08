package files

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
)

// PublicFolder creates a public folder with .gitkeep for Payload CMS apps.
// This ensures that Next.js builds do not fail due to missing public directory.
func PublicFolder(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	// Filter to only Payload apps.
	payloadApps := appDef.GetAppsByType(appdef.AppTypePayload)

	for _, app := range payloadApps {
		publicPath := filepath.Join(app.Path, "public")

		// Check if public folder already exists.
		exists, err := afero.DirExists(input.FS, publicPath)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("checking if public folder exists for app %s", app.Name))
		}

		// Skip if folder already exists.
		if exists {
			input.Printer().Println(fmt.Sprintf("• skipping %s - public folder already exists", app.Name))
			continue
		}

		// Create .gitkeep file in public folder.
		gitkeepPath := filepath.Join(publicPath, ".gitkeep")
		err = input.Generator().Bytes(gitkeepPath, []byte{},
			scaffold.WithTracking(manifest.SourceApp(app.Name)),
			scaffold.WithoutNotice(),
		)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("creating public folder for app %s", app.Name))
		}
	}

	return nil
}
