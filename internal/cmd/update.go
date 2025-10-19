package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/cicd"
	"github.com/ainsleydev/webkit/internal/cmd/files"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/scaffold"
)

var updateCmd = &cli.Command{
	Name:        "update",
	Usage:       "Update project dependencies from app.json",
	Description: "Rebuilds all generated files based on current app.json configuration",
	Action:      cmdtools.Wrap(update),
}

var updateOps = []cmdtools.RunCommand{
	files.CodeStyle,
	files.CreateGitSettings,
	files.CreatePackageJson,
	cicd.CreatePRWorkflow,
	cicd.BackupResourcesWorkflow,
}

func update(ctx context.Context, input cmdtools.CommandInput) error {
	gen := scaffold.New(input.FS, input.Manifest)

	// 1. Load previous manifest
	oldManifest, err := manifest.Load(input.FS, ".webkit/generated.json")
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("loading manifest: %w", err)
	}

	// 2. Generate all files (they auto-track to new manifest)
	for _, op := range updateOps {
		if err := op(ctx, input); err != nil {
			return err
		}
	}

	// 3. Save new manifest
	if err := gen.Finalize(); err != nil {
		return fmt.Errorf("saving manifest: %w", err)
	}

	// 4. Cleanup orphaned files
	if oldManifest != nil {
		newManifest, _ := manifest.Load(input.FS, ".webkit/generated.json")
		if err := manifest.Cleanup(input.FS, oldManifest, newManifest, gen.Printer); err != nil {
			return fmt.Errorf("cleaning up: %w", err)
		}
	}

	return nil
}
