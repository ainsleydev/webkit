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
	"github.com/ainsleydev/webkit/internal/pkgjson"
)

const (
	// payloadOwner is the GitHub organisation that owns the Payload repository.
	payloadOwner = "payloadcms"
	// payloadRepo is the GitHub repository name for Payload CMS.
	payloadRepo = "payload"
)

var BumpCmd = &cli.Command{
	Name:  "bump",
	Usage: "Bump Payload CMS and associated dependencies to the latest version",
	Description: "Fetches the latest stable Payload CMS release from GitHub and updates all " +
		"Payload-related dependencies in package.json files for Payload apps. " +
		"This includes payload, @payloadcms/* packages, AND all dependencies that Payload itself uses " +
		"(e.g., lexical, @lexical/headless, etc.) to ensure version compatibility.",
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

// Bump updates all Payload CMS and associated dependencies to the latest version.
// It fetches Payload's package.json from GitHub to determine which dependencies to update,
// ensuring compatibility with the target Payload version.
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

	// Create GitHub client.
	ghClient := ghapi.New("")

	// Determine target version.
	var targetVersion string
	if v := input.Command.String("version"); v != "" {
		targetVersion = v
		printer.Printf("Using specified version: %s\n", targetVersion)
	} else {
		printer.Println("Fetching latest Payload CMS version from GitHub...")
		version, err := ghClient.GetLatestRelease(ctx, payloadOwner, payloadRepo)
		if err != nil {
			return errors.Wrap(err, "fetching latest Payload version")
		}
		targetVersion = version
		printer.Success(fmt.Sprintf("Latest version: %s", targetVersion))
	}

	// Fetch Payload's dependencies from GitHub to know what to bump.
	printer.Println("Fetching Payload's dependencies...")
	payloadDeps, err := FetchPayloadDependencies(ctx, ghClient, targetVersion)
	if err != nil {
		return errors.Wrap(err, "fetching Payload dependencies")
	}
	printer.Success(fmt.Sprintf("Found %d dependencies in Payload %s", len(payloadDeps.AllDeps), targetVersion))

	printer.LineBreak()

	isDryRun := input.Command.Bool("dry-run")
	if isDryRun {
		printer.Println("🔍 DRY RUN - No changes will be made")
		printer.LineBreak()
	}

	// Process each Payload app.
	var hasChanges bool
	for _, app := range payloadApps {
		changed, err := bumpAppDependencies(ctx, input, app, targetVersion, payloadDeps, isDryRun)
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

// bumpAppDependencies updates Payload and associated dependencies for a single app.
// Returns true if any changes were made.
func bumpAppDependencies(
	_ context.Context,
	input cmdtools.CommandInput,
	app appdef.App,
	version string,
	payloadDeps *PayloadDependencies,
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
	pkg, err := pkgjson.Read(input.FS, pkgPath)
	if err != nil {
		return false, err
	}

	// Create a matcher that matches:
	// 1. payload and @payloadcms/* packages (always update to target version)
	// 2. Any dependency that Payload itself uses (update to Payload's version)
	matcher := func(name string) bool {
		// Always match payload and @payloadcms/* packages.
		if name == "payload" {
			return true
		}
		if len(name) > len("@payloadcms/") && name[0:len("@payloadcms/")] == "@payloadcms/" {
			return true
		}
		// Match if this dependency is in Payload's package.json.
		_, inPayload := payloadDeps.AllDeps[name]
		return inPayload
	}

	// Check if this package has any matchable dependencies.
	if !pkgjson.HasAnyDependency(pkg, matcher) {
		printer.Printf("• %s - no Payload-related dependencies found\n", app.Name)
		return false, nil
	}

	// Create a version formatter that:
	// - Uses the target version for payload and @payloadcms/* packages
	// - Uses Payload's version for other dependencies
	// - Respects exact vs caret formatting based on dependency type
	versionFormatter := func(name, _ string) string {
		// For payload and @payloadcms/* packages, use the target version.
		isPayloadPackage := name == "payload" || (len(name) > len("@payloadcms/") && name[0:len("@payloadcms/")] == "@payloadcms/")
		if isPayloadPackage {
			useExact := pkgjson.IsDevDependency(pkg, name)
			return pkgjson.FormatVersion(version, useExact)
		}

		// For other dependencies, use Payload's version.
		if payloadVer, ok := payloadDeps.AllDeps[name]; ok {
			// Preserve exact versions from Payload (they know what they're doing).
			return payloadVer
		}

		// This shouldn't happen since matcher should prevent this.
		return ""
	}

	// Update dependencies.
	result := pkgjson.UpdateDependencies(pkg, matcher, versionFormatter)

	if len(result.Updated) == 0 {
		printer.Printf("✓ %s - already up to date\n", app.Name)
		return false, nil
	}

	// Display changes.
	printer.Printf("📦 %s (%s)\n", app.Name, pkgPath)
	for _, dep := range result.Updated {
		oldVer := result.OldVersions[dep]
		newVer := versionFormatter(dep, oldVer)
		if dryRun {
			printer.Printf("   %s: %s → %s\n", dep, oldVer, newVer)
		} else {
			printer.Printf("   ✓ %s: %s → %s\n", dep, oldVer, newVer)
		}
	}

	// Write updated package.json if not a dry run.
	if !dryRun {
		if err := pkgjson.Write(input.FS, pkgPath, pkg); err != nil {
			return false, err
		}
	}

	return true, nil
}
