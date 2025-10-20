package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
)

var driftCmd = &cli.Command{
	Name:        "drift",
	Usage:       "Detect manual modifications to generated files",
	Description: "Checks if generated files have been manually modified or deleted since the last webkit update",
	Action:      cmdtools.Wrap(drift),
}

// drift detects if tracked files have been manually modified or deleted.
//
// Note: This only detects drift from the last webkit update. If templates have
// been updated in a newer version of WebKit but webkit update hasn't been run,
// this will not detect that the files are outdated.
func drift(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()

	printer.Info("Checking for drift...")

	// Run update in memory to see what should be generated
	memFS := afero.NewMemMapFs()
	memTracker := manifest.NewTracker()

	memInput := cmdtools.CommandInput{
		FS:          memFS,
		AppDefCache: input.AppDef(),
		Manifest:    memTracker,
		BaseDir:     input.BaseDir,
	}
	memInput.Printer().SetWriter(io.Discard)

	if err := update(ctx, memInput); err != nil {
		return errors.Wrap(err, "running update for in-mem manifest")
	}

	// Compare actual vs expected
	drifted, err := manifest.DetectDrift(input.FS, memFS)
	if err != nil {
		return errors.Wrap(err, "detecting drift")
	}

	if len(drifted) == 0 {
		printer.Success("✓ No drift detected - all files are up to date")
		return nil
	}

	// Group by type
	modifiedFiles := manifest.DriftReasonModified.FilterEntries(drifted)
	outdatedFiles := manifest.DriftReasonOutdated.FilterEntries(drifted)
	newFiles := manifest.DriftReasonNew.FilterEntries(drifted)
	deletedFiles := manifest.DriftReasonDeleted.FilterEntries(drifted)

	// Report findings
	if len(modifiedFiles) > 0 {
		printer.Error("⚠ Manual modifications detected:")
		printer.Println("  These files were manually edited:")
		for _, d := range modifiedFiles {
			printer.Println(fmt.Sprintf("    • %s", d.Path))
		}
		printer.LineBreak()
	}

	if len(outdatedFiles) > 0 {
		printer.Error("Outdated files detected:")
		printer.Println("app.json changed, these files need regeneration:")
		for _, d := range outdatedFiles {
			printer.Println(fmt.Sprintf("    • %s", d.Path))
		}
		printer.LineBreak()
	}

	if len(newFiles) > 0 {
		printer.Error("Missing files detected:")
		printer.Println("These files should exist:")
		for _, d := range newFiles {
			printer.Println(fmt.Sprintf("    • %s", d.Path))
		}
		printer.LineBreak()
	}

	if len(deletedFiles) > 0 {
		printer.Warn("Orphaned files detected:")
		printer.Println("These files should be removed:")
		for _, d := range deletedFiles {
			printer.Println(fmt.Sprintf("    • %s", d.Path))
		}
		printer.LineBreak()
	}

	printer.Info("Run 'webkit update' to sync all files")

	return cmdtools.ExitWithCode(1)
}
