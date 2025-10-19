package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
)

var driftCmd = &cli.Command{
	Name:        "drift",
	Description: "Detect file drift caused by outdated WebKit templates",
	Action:      cmdtools.Wrap(drift),
}

func drift(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()

	mani, err := manifest.Load(input.FS)
	if err != nil {
		return err
	}

	files := manifest.DetectDrift(input.FS, mani)
	if files == nil {
		printer.Success("No drift detected")
		return nil
	}

	printer.Error("Drift found")
	printer.Println("Be sure to run webkit update to update the projects dependencies.")

	list := make([]string, len(files))
	for idx, file := range files {
		list[idx] = fmt.Sprintf("%s\n", file)
	}
	printer.LineBreak()
	printer.List(list)
	printer.LineBreak()

	os.Exit(1)

	return nil
}
