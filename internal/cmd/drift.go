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
	Description: "Detect file drift caused by outdated WebKit templates",
	Action:      cmdtools.Wrap(drift),
}

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
		printer.Println("Deleted files:")
		printer.List(deleted)
		printer.LineBreak()
	}

	return cmdtools.ExitWithCode(1)
}
