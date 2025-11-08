package payload

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/ghapi"
)

const (
	// payloadOwner is the GitHub organisation that owns the Payload repository.
	payloadOwner = "payloadcms"
	// payloadRepo is the GitHub repository name for Payload CMS.
	payloadRepo = "payload"
)

var BumpCmd = &cli.Command{
	Name:  "bump",
	Usage: "Bump Payload CMS dependencies to the latest version",
	Description: "Fetches the latest stable Payload CMS release from GitHub and updates all " +
		"payload and @payloadcms/* dependencies in package.json files for Payload apps.",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "dry-run",
			Aliases: []string{"d"},
			Usage:   "Show what would be changed without making modifications",
		},
		&cli.StringFlag{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Bump to a specific version instead of fetching the latest",
		},
	},
	Action: cmdtools.Wrap(Bump),
}

// Bump updates all Payload CMS dependencies to the latest version across all Payload apps.
// It fetches the latest stable release from GitHub (or uses a specified version) and updates
// all package.json files for apps with type "payload".
func Bump(ctx context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()
	printer := input.Printer()

	// Find all Payload apps.
	payloadApps := findPayloadApps(appDef)
	if len(payloadApps) == 0 {
		printer.Println("No Payload CMS apps found in app.json")
		return nil
	}

	printer.Printf("Found %d Payload app(s)\n", len(payloadApps))

	// Determine target version.
	var targetVersion string
	if v := input.Command.String("version"); v != "" {
		targetVersion = v
		printer.Printf("Using specified version: %s\n", targetVersion)
	} else {
		printer.Println("Fetching latest Payload CMS version from GitHub...")
		client := ghapi.New("")
		version, err := client.GetLatestRelease(ctx, payloadOwner, payloadRepo)
		if err != nil {
			return errors.Wrap(err, "fetching latest Payload version")
		}
		targetVersion = version
		printer.Success(fmt.Sprintf("Latest version: %s", targetVersion))
	}

	printer.LineBreak()

	isDryRun := input.Command.Bool("dry-run")
	if isDryRun {
		printer.Println("🔍 DRY RUN - No changes will be made")
		printer.LineBreak()
	}

	// Process each Payload app.
	var hasChanges bool
	for _, app := range payloadApps {
		changed, err := bumpAppDependencies(ctx, input, app, targetVersion, isDryRun)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("bumping dependencies for app %s", app.Name))
		}
		if changed {
			hasChanges = true
		}
	}

	printer.LineBreak()

	if !hasChanges {
		printer.Println("All Payload dependencies are already up to date!")
		return nil
	}

	if isDryRun {
		printer.Println("Dry run complete. Run without --dry-run to apply changes.")
	} else {
		printer.Success("Successfully bumped Payload dependencies!")
		printer.Println("\n💡 Run 'pnpm install' to update your lockfile")
	}

	return nil
}

// findPayloadApps returns all apps with type "payload".
func findPayloadApps(appDef *appdef.Definition) []appdef.App {
	var apps []appdef.App
	for _, app := range appDef.Apps {
		if app.Type == appdef.AppTypePayload {
			apps = append(apps, app)
		}
	}
	return apps
}

// bumpAppDependencies updates Payload dependencies for a single app.
// Returns true if any changes were made.
func bumpAppDependencies(
	_ context.Context,
	input cmdtools.CommandInput,
	app appdef.App,
	version string,
	dryRun bool,
) (bool, error) {
	printer := input.Printer()
	pkgPath := filepath.Join(app.Path, "package.json")

	// Check if package.json exists.
	exists, err := afero.Exists(input.FS, pkgPath)
	if err != nil {
		return false, errors.Wrap(err, "checking if package.json exists")
	}
	if !exists {
		printer.Printf("⚠️  Skipping %s - package.json not found at %s\n", app.Name, pkgPath)
		return false, nil
	}

	// Read package.json.
	pkg, err := ReadPackageJSON(input.FS, pkgPath)
	if err != nil {
		return false, err
	}

	// Check if this package has any Payload dependencies.
	if !HasPayloadDependencies(pkg) {
		printer.Printf("• %s - no Payload dependencies found\n", app.Name)
		return false, nil
	}

	// Bump dependencies.
	result := BumpPayloadDependencies(pkg, version)
	result.Path = pkgPath

	if len(result.Bumped) == 0 {
		printer.Printf("✓ %s - already at version %s\n", app.Name, version)
		return false, nil
	}

	// Display changes.
	printer.Printf("📦 %s (%s)\n", app.Name, pkgPath)
	for _, dep := range result.Bumped {
		oldVer := result.OldVersions[dep]
		newVer := formatVersion(version, isDevDependency(pkg, dep))
		if dryRun {
			printer.Printf("   %s: %s → %s\n", dep, oldVer, newVer)
		} else {
			printer.Printf("   ✓ %s: %s → %s\n", dep, oldVer, newVer)
		}
	}

	// Write updated package.json if not a dry run.
	if !dryRun {
		if err := WritePackageJSON(input.FS, pkgPath, pkg); err != nil {
			return false, err
		}
	}

	return true, nil
}

// isDevDependency checks if a package is in devDependencies.
func isDevDependency(pkg *PackageJSON, name string) bool {
	_, ok := pkg.DevDependencies[name]
	return ok
}
