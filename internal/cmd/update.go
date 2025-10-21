package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/cicd"
	"github.com/ainsleydev/webkit/internal/cmd/env"
	"github.com/ainsleydev/webkit/internal/cmd/files"
	"github.com/ainsleydev/webkit/internal/cmd/secrets"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
)

var updateCmd = &cli.Command{
	Name:        "update",
	Usage:       "Update project dependencies from app.json",
	Description: "Rebuilds all generated files based on current app.json configuration",
	Action:      cmdtools.Wrap(update),
}

// runner defines an operation to execute during the update process,
// pairing a command function with a descriptive name for logging.
type runner struct {
	command cmdtools.RunCommand
	name    string
}

// updateOps defines all operations to run during an update, in order.
// The order is significant: manifest and definition must be first,
// and environment/secrets sync should be last.
var updateOps = []runner{
	{files.Manifest, "Manifest: Scaffold manifest files"},
	{files.Definition, "Definition: Update webkit_version in app.json"},
	{env.Scaffold, "Env: Scaffold .env files"},
	{secrets.Scaffold, "Secrets: Scaffold secret files"},

	// Alphabetically
	{files.CodeStyle, "Files: Create code style files"},
	{files.GitSettings, "Files: Create git settings"},
	{files.PackageJSON, "Files: Create package.json"},
	{files.TurboJSON, "Files: Create turbo.json"},
	{cicd.PR, "CICD: Create PR workflows"},
	{cicd.BackupWorkflow, "CICD: Create backup workflows"},
	{cicd.ActionTemplates, "CICD: Create action templates"},

	// Lastly
	{env.Scaffold, "Env: Sync .env files"},
	{secrets.Sync, "Secrets: Sync secret files"},
}

// update regenerates all project files based on the current app.json configuration.
// It loads the previous manifest, executes all update operations, saves the new manifest,
// and cleans up any orphaned files from the previous generation.
func update(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()

	printer.Info("Updating project dependencies...")
	printer.LineBreak()

	// 1. Load previous manifest
	oldManifest, err := manifest.Load(input.FS)
	if err != nil && !errors.Is(err, manifest.ErrNoManifest) {
		return errors.Wrap(err, "loading manifest")
	}

	// 2. Generate all files (they auto-track to new manifest)
	for _, op := range updateOps {
		printer.Printf("🏃 %v\n", op.name)
		if err = op.command(ctx, input); err != nil {
			return err
		}
	}

	// 3. Save new manifest
	if err = input.Manifest.Save(input.FS); err != nil {
		return errors.Wrap(err, "saving manifest")
	}

	// 4. Cleanup orphaned files
	if oldManifest != nil {
		newManifest, err := manifest.Load(input.FS)
		if err != nil {
			return errors.Wrap(err, "loading manifest")
		}
		if err = manifest.Cleanup(input.FS, oldManifest, newManifest, input.Printer()); err != nil {
			return errors.Wrap(err, "cleaning up manifest")
		}
	}

	printer.LineBreak()
	printer.Success("Successfully updated project dependencies!")

	return nil
}
