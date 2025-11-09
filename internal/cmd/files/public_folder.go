package files

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/fsext"
	"github.com/ainsleydev/webkit/internal/scaffold"
)

// PublicFolder creates a public folder with .gitkeep for Payload CMS apps.
// This ensures that Next.js builds do not fail due to missing public directory.
// The .gitkeep file is not tracked in the manifest as it's meant to be temporary
// and can be safely removed once actual files are added to the public folder.
func PublicFolder(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	// Filter to only Payload apps.
	payloadApps := appDef.GetAppsByType(appdef.AppTypePayload)

	for _, app := range payloadApps {
		publicPath := filepath.Join(app.Path, "public")

		// Check if folder has any files (excluding .gitkeep).
		hasFiles, err := folderHasFiles(input.FS, publicPath)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("checking public folder contents for app %s", app.Name))
		}

		// Skip if folder already has files.
		if hasFiles {
			input.Printer().Println(fmt.Sprintf("• skipping %s - public folder has files", app.Name))
			continue
		}

		// Create .gitkeep file to ensure empty folder is tracked by git.
		gitkeepPath := filepath.Join(publicPath, ".gitkeep")
		err = input.Generator().Bytes(gitkeepPath, []byte{}, scaffold.WithoutNotice())
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("creating .gitkeep for app %s", app.Name))
		}

		input.Printer().Println(fmt.Sprintf("• created %s/public/.gitkeep", app.Name))
	}

	return nil
}

// folderHasFiles checks if a folder exists and contains any files other than .gitkeep.
func folderHasFiles(fs afero.Fs, path string) (bool, error) {
	exists := fsext.DirExists(fs, path)

	// Folder doesn't exist, so no files.
	if !exists {
		return false, nil
	}

	// Read folder contents.
	entries, err := afero.ReadDir(fs, path)
	if err != nil {
		return false, err
	}

	// Check if there are any files other than .gitkeep.
	for _, entry := range entries {
		if entry.Name() != ".gitkeep" {
			return true, nil
		}
	}

	return false, nil
}
