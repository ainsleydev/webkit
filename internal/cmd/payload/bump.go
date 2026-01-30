package payload

import (
	"context"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/fsext"
	"github.com/ainsleydev/webkit/internal/ghapi"
	"github.com/ainsleydev/webkit/internal/pkgjson"
	"github.com/ainsleydev/webkit/internal/util/executil"
)

var BumpCmd = &cli.Command{
	Name:  "bump",
	Usage: "Bump Payload CMS and associated dependencies to the latest version",
	Description: "Fetches the latest stable Payload CMS release from GitHub and updates all " +
		"Payload-related dependencies in the current directory's package.json. " +
		"This includes payload, @payloadcms/* packages, AND all dependencies that Payload itself uses " +
		"(e.g., lexical, @lexical/headless, etc.) to ensure version compatibility. " +
		"Run this command from within a Payload project directory.",
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

// Bump updates all Payload CMS and associated dependencies to the latest version.
// It fetches Payload's template package.json to determine which dependencies to update,
// ensuring compatibility with the target Payload version.
// This command should be run from within a Payload project directory.
func Bump(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()

	// Check if we're in a Payload project directory
	pkgPath := "package.json"
	if !fsext.Exists(input.FS, pkgPath) {
		return errors.New("package.json not found in current directory. Please run this command from a Payload project directory")
	}

	// Read package.json to verify it's a Payload project
	pkg, err := pkgjson.Read(input.FS, pkgPath)
	if err != nil {
		return errors.Wrap(err, "reading package.json")
	}

	// Check if this is a Payload project by looking for Payload dependencies
	hasPayloadDeps := pkg.HasAnyDependency(func(name string) bool {
		return name == "payload" || (len(name) > len("@payloadcms/") && name[0:len("@payloadcms/")] == "@payloadcms/")
	})

	if !hasPayloadDeps {
		return errors.New("no Payload dependencies found in package.json. This doesn't appear to be a Payload project")
	}

	// Create GitHub client with the token that comes
	// out of the box with Terraform.
	var ghClient ghapi.Client
	if token := os.Getenv("GITHUB_TOKEN_CLASSIC"); token != "" {
		ghClient = ghapi.New(token)
	} else {
		ghClient = ghapi.NewWithoutAuth()
	}

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

		printer.Success(fmt.Sprintf("Latest version: %s\n", targetVersion))
	}

	// Fetch Payload's dependencies from the template to know what to bump.
	printer.Println("Fetching Payload's template dependencies...")
	payloadDeps, err := fetchPayloadDependencies(ctx)
	if err != nil {
		return errors.Wrap(err, "fetching Payload dependencies")
	}
	printer.Success(fmt.Sprintf("Found %d dependencies in Payload template\n", len(payloadDeps.AllDeps)))

	isDryRun := input.Command.Bool("dry-run")
	if isDryRun {
		printer.Println("🔍 DRY RUN - No changes will be made")
		printer.LineBreak()
	}

	// Process the current directory's Payload app
	changed, err := bumpAppDependencies(input, targetVersion, payloadDeps, isDryRun)
	if err != nil {
		return errors.Wrap(err, "bumping dependencies")
	}

	printer.LineBreak()

	if !changed {
		printer.Println("All Payload dependencies are already up to date!")
		return nil
	}

	if isDryRun {
		printer.Println("Dry run complete. Run without --dry-run to apply changes.")
		return nil
	}

	printer.Success("Successfully bumped Payload dependencies!")

	// Run pnpm install unless --no-install is specified
	if !input.Command.Bool("no-install") {
		if err := runPnpmInstall(ctx, input); err != nil {
			return errors.Wrap(err, "running pnpm install")
		}
	}

	// Run pnpm migrate:create unless --no-migrate is specified
	if !input.Command.Bool("no-migrate") {
		if err := runPnpmMigrate(ctx, input); err != nil {
			return errors.Wrap(err, "running pnpm migrate:create")
		}
	}

	if input.Command.Bool("no-install") {
		printer.Println("\n💡 Run 'pnpm install' in each app directory to update your lockfile")
	}

	return nil
}

const (
	payloadOwner            = "payloadcms"
	payloadRepo             = "payload"
	payloadBlankTemplateURL = "https://raw.githubusercontent.com/payloadcms/payload/refs/heads/main/templates/blank/package.json"
)

// payloadDependencies contains all dependencies from Payload's package.json.
// This is used to determine which dependencies should be bumped when updating.
type payloadDependencies struct {
	Dependencies    map[string]string
	DevDependencies map[string]string
	AllDeps         map[string]string // Combined for easier lookup
}

// ============================================================================
// Helpers
// ============================================================================

// fetchPayloadDependencies fetches Payload CMS's template package.json from GitHub
// and extracts all its dependencies (excluding workspace:* versions).
//
// This uses the blank template's package.json which represents the actual dependencies
// a new Payload project would have. For payload and @payloadcms/* packages that use
// workspace:*, we'll fetch the actual version from GitHub releases.
func fetchPayloadDependencies(ctx context.Context) (*payloadDependencies, error) {
	// Fetch template package.json from GitHub (always use main branch).
	pkg, err := pkgjson.FetchFromRemote(ctx, payloadBlankTemplateURL)
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

// bumpAppDependencies updates Payload and associated dependencies in the current directory.
// Returns true if any changes were made.
func bumpAppDependencies(
	input cmdtools.CommandInput,
	version string,
	payloadDeps *payloadDependencies,
	dryRun bool,
) (bool, error) {
	printer := input.Printer()
	pkgPath := "package.json"

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

	// Display skipped downgrades.
	if len(result.Skipped) > 0 {
		printer.Printf("⏭ Skipping %d dependencies (would downgrade):\n", len(result.Skipped))
		for _, dep := range result.Skipped {
			currentVer := pkg.Dependencies[dep]
			if currentVer == "" {
				currentVer = pkg.DevDependencies[dep]
			}
			if currentVer == "" {
				currentVer = pkg.PeerDependencies[dep]
			}
			printer.Printf("   %s: %s (keeping current version)\n", dep, currentVer)
		}
	}

	if len(result.Updated) == 0 {
		printer.Println("✓ All dependencies already up to date")
		return false, nil
	}

	// Display changes.
	printer.Printf("📦 Updating dependencies in %s\n", pkgPath)
	for i, dep := range result.Updated {
		oldVer := result.OldVersions[dep]
		newVer := versionFormatter(dep, oldVer)

		lineBreak := "\n"
		if i == len(result.Updated)-1 {
			lineBreak = "" // no line break for the last one
		}

		if dryRun {
			printer.Printf("   %s: %s → %s%s", dep, oldVer, newVer, lineBreak)
		} else {
			printer.Printf("   ✓ %s: %s → %s%s", dep, oldVer, newVer, lineBreak)
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

// runPnpmInstall runs pnpm install in the current directory to update the lockfile.
func runPnpmInstall(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()
	spinner := input.Spinner()

	printer.LineBreak()
	printer.Println("Installing dependencies...")
	spinner.Start()
	defer spinner.Stop()

	cmd := executil.NewCommand("pnpm", "install")
	cmd.Dir = input.BaseDir

	result, err := input.Runner.Run(ctx, cmd)
	spinner.Stop()

	if err != nil {
		printer.Error("Failed to run pnpm install")
		if result.Output != "" {
			printer.Println(result.Output)
		}
		return err
	}

	printer.Success("Dependencies installed")
	return nil
}

// runPnpmMigrate runs pnpm migrate:create in the current directory to create database migrations.
func runPnpmMigrate(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()
	spinner := input.Spinner()

	printer.LineBreak()
	printer.Println("Creating database migrations...")
	spinner.Start()
	defer spinner.Stop()

	cmd := executil.NewCommand("pnpm", "migrate:create")
	cmd.Dir = input.BaseDir

	result, err := input.Runner.Run(ctx, cmd)
	spinner.Stop()

	if err != nil {
		printer.Error("Failed to run pnpm migrate:create")
		if result.Output != "" {
			printer.Println(result.Output)
		}
		return err
	}

	printer.Success("Database migrations created")
	return nil
}
