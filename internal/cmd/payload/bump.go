package payload

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/ghapi"
	"github.com/ainsleydev/webkit/internal/pkgjson"
)

// ============================================================================
// Constants
// ============================================================================

const (
	// payloadOwner is the GitHub organisation that owns the Payload repository.
	payloadOwner = "payloadcms"
	// payloadRepo is the GitHub repository name for Payload CMS.
	payloadRepo = "payload"
	// payloadTemplateURL is the URL to Payload's blank template package.json.
	payloadTemplateURL = "https://raw.githubusercontent.com/payloadcms/payload/refs/heads/main/templates/blank/package.json"
)

// ============================================================================
// Types
// ============================================================================

// payloadDependencies contains all dependencies from Payload's package.json.
// This is used to determine which dependencies should be bumped when updating.
type payloadDependencies struct {
	Dependencies    map[string]string
	DevDependencies map[string]string
	AllDeps         map[string]string // Combined for easier lookup
}

// ============================================================================
// Command
// ============================================================================

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
		&cli.BoolFlag{
			Name:  "no-install",
			Usage: "Skip running pnpm install after bumping dependencies",
		},
		&cli.BoolFlag{
			Name:  "no-migrate",
			Usage: "Skip running pnpm migrate:create after installing dependencies",
		},
	},
	Action: cmdtools.Wrap(Bump),
}

// ============================================================================
// Main bump function
// ============================================================================

// Bump updates all Payload CMS and associated dependencies to the latest version.
// It fetches Payload's template package.json to determine which dependencies to update,
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

	// Create GitHub client with the token that comes
	// out of the box with Terraform.
	ghClient := ghapi.New(os.Getenv("GITHUB_TOKEN_CLASSIC"))

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

	// Fetch Payload's dependencies from the template to know what to bump.
	printer.Println("Fetching Payload's template dependencies...")
	payloadDeps, err := fetchPayloadDependencies(ctx)
	if err != nil {
		return errors.Wrap(err, "fetching Payload dependencies")
	}
	printer.Success(fmt.Sprintf("Found %d dependencies in Payload template", len(payloadDeps.AllDeps)))

	printer.LineBreak()

	isDryRun := input.Command.Bool("dry-run")
	if isDryRun {
		printer.Println("🔍 DRY RUN - No changes will be made")
		printer.LineBreak()
	}

	// Process each Payload app and track which ones changed.
	var changedApps []appdef.App
	for _, app := range payloadApps {
		changed, err := bumpAppDependencies(ctx, input, app, targetVersion, payloadDeps, isDryRun)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("bumping dependencies for app %s", app.Name))
		}
		if changed {
			changedApps = append(changedApps, app)
		}
	}

	printer.LineBreak()

	if len(changedApps) == 0 {
		printer.Println("All Payload dependencies are already up to date!")
		return nil
	}

	if isDryRun {
		printer.Println("Dry run complete. Run without --dry-run to apply changes.")
		return nil
	}

	printer.Success("Successfully bumped Payload dependencies!")

	// Run pnpm install and migrate for each changed app
	for _, app := range changedApps {
		// Run pnpm install unless --no-install is specified
		if !input.Command.Bool("no-install") {
			if err := runPnpmInstall(ctx, input, app); err != nil {
				return errors.Wrap(err, fmt.Sprintf("running pnpm install for %s", app.Name))
			}
		}

		// Run pnpm migrate:create unless --no-migrate is specified
		if !input.Command.Bool("no-migrate") {
			if err := runPnpmMigrate(ctx, input, app); err != nil {
				return errors.Wrap(err, fmt.Sprintf("running pnpm migrate:create for %s", app.Name))
			}
		}
	}

	if input.Command.Bool("no-install") {
		printer.Println("\n💡 Run 'pnpm install' in each app directory to update your lockfile")
	}

	return nil
}

// ============================================================================
// Helper functions
// ============================================================================

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

// fetchPayloadDependencies fetches Payload CMS's template package.json from GitHub
// and extracts all its dependencies (excluding workspace:* versions).
//
// This uses the blank template's package.json which represents the actual dependencies
// a new Payload project would have. For payload and @payloadcms/* packages that use
// workspace:*, we'll fetch the actual version from GitHub releases.
func fetchPayloadDependencies(ctx context.Context) (*payloadDependencies, error) {
	// Fetch template package.json from GitHub (always use main branch).
	pkg, err := pkgjson.FetchFromRemote(ctx, payloadTemplateURL)
	if err != nil {
		return nil, errors.Wrap(err, "fetching template package.json")
	}

	// Combine all dependencies, filtering out workspace:* versions.
	allDeps := make(map[string]string)
	for name, version := range pkg.Dependencies {
		if version != "workspace:*" {
			allDeps[name] = version
		}
	}
	for name, version := range pkg.DevDependencies {
		if version != "workspace:*" {
			allDeps[name] = version
		}
	}

	return &payloadDependencies{
		Dependencies:    pkg.Dependencies,
		DevDependencies: pkg.DevDependencies,
		AllDeps:         allDeps,
	}, nil
}

// bumpAppDependencies updates Payload and associated dependencies for a single app.
// Returns true if any changes were made.
func bumpAppDependencies(
	_ context.Context,
	input cmdtools.CommandInput,
	app appdef.App,
	version string,
	payloadDeps *payloadDependencies,
	dryRun bool,
) (bool, error) {
	printer := input.Printer()
	pkgPath := filepath.Join(app.Path, "package.json")

	// Check if package.json exists.
	if !pkgjson.Exists(input.FS, pkgPath) {
		printer.Printf("⚠️  Skipping %s - package.json not found at %s\n", app.Name, pkgPath)
		return false, nil
	}

	// Read package.json.
	pkg, err := pkgjson.Read(input.FS, pkgPath)
	if err != nil {
		return false, err
	}

	// Create a matcher that matches any dependency in Payload's template.
	// This includes all dependencies from the template, plus payload and @payloadcms/* packages.
	matcher := func(name string) bool {
		// Match if in template's dependencies.
		if _, ok := payloadDeps.AllDeps[name]; ok {
			return true
		}
		// Also match payload and @payloadcms/* packages (even if they had workspace:*).
		if name == "payload" || (len(name) > len("@payloadcms/") && name[0:len("@payloadcms/")] == "@payloadcms/") {
			return true
		}
		return false
	}

	// Check if this package has any matchable dependencies.
	if !pkg.HasAnyDependency(matcher) {
		printer.Printf("• %s - no Payload-related dependencies found\n", app.Name)
		return false, nil
	}

	// Create a version formatter that:
	// - Uses the target version for payload and @payloadcms/* packages
	// - Uses template's version for other dependencies
	// - Respects exact vs caret formatting based on dependency type
	versionFormatter := func(name, _ string) string {
		// For payload and @payloadcms/* packages, use the target version.
		isPayloadPackage := name == "payload" || (len(name) > len("@payloadcms/") && name[0:len("@payloadcms/")] == "@payloadcms/")
		if isPayloadPackage {
			useExact := pkg.IsDevDependency(name)
			return pkgjson.FormatVersion(version, useExact)
		}

		// For other dependencies, use template's version.
		if templateVer, ok := payloadDeps.AllDeps[name]; ok {
			return templateVer
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
		if err = pkgjson.Write(input.FS, pkgPath, pkg); err != nil {
			return false, err
		}
	}

	return true, nil
}

// runPnpmInstall runs pnpm install in the app directory to update the lockfile.
func runPnpmInstall(ctx context.Context, input cmdtools.CommandInput, app appdef.App) error {
	printer := input.Printer()
	spinner := input.Spinner()

	printer.LineBreak()
	printer.Printf("Installing dependencies for %s...\n", app.Name)
	spinner.Start()
	defer spinner.Stop()

	appDir := filepath.Join(input.BaseDir, app.Path)
	cmd := exec.CommandContext(ctx, "pnpm", "install")
	cmd.Dir = appDir

	output, err := cmd.CombinedOutput()
	spinner.Stop()

	if err != nil {
		printer.Error(fmt.Sprintf("Failed to run pnpm install for %s", app.Name))
		if len(output) > 0 {
			printer.Println(string(output))
		}
		return err
	}

	printer.Success(fmt.Sprintf("Dependencies installed for %s", app.Name))
	return nil
}

// runPnpmMigrate runs pnpm migrate:create in the app directory to create database migrations.
func runPnpmMigrate(ctx context.Context, input cmdtools.CommandInput, app appdef.App) error {
	printer := input.Printer()
	spinner := input.Spinner()

	printer.LineBreak()
	printer.Printf("Creating database migrations for %s...\n", app.Name)
	spinner.Start()
	defer spinner.Stop()

	appDir := filepath.Join(input.BaseDir, app.Path)
	cmd := exec.CommandContext(ctx, "pnpm", "migrate:create")
	cmd.Dir = appDir

	output, err := cmd.CombinedOutput()
	spinner.Stop()

	if err != nil {
		printer.Error(fmt.Sprintf("Failed to run pnpm migrate:create for %s", app.Name))
		if len(output) > 0 {
			printer.Println(string(output))
		}
		return err
	}

	printer.Success(fmt.Sprintf("Database migrations created for %s", app.Name))
	return nil
}
