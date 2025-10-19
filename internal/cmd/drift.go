package cmd

import (
	"context"
	"errors"

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
func drift(_ context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()

	mani, err := manifest.Load(input.FS)
	if err != nil && errors.Is(err, manifest.ErrNoManifest) {
		return nil
	} else if err != nil {
		return err
	}

	drifted := manifest.DetectDrift(input.FS, mani)
	if drifted == nil {
		printer.Success("No drift detected")
		return nil
	}

	printer.Error("Drift found")
	printer.Println("Action Required: Run webkit update to sync your project.")
	printer.LineBreak()

	// Group by reason for better output
	var modified []string
	var deleted []string

	for _, file := range drifted {
		switch file.Reason {
		case manifest.DriftReasonModified:
			modified = append(modified, file.Path)
		case manifest.DriftReasonDeleted:
			deleted = append(deleted, file.Path)
		}
	}

	if len(modified) > 0 {
		printer.Println("Modified files:")
		printer.List(modified)
		printer.LineBreak()
	}

	if len(deleted) > 0 {
		printer.LineBreak()
		printer.Println("Deleted files:")
		printer.List(deleted)
		printer.LineBreak()
	}

	printer.LineBreak()

	return cmdtools.ExitWithCode(1)
}
