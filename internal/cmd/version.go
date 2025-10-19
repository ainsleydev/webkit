package cmd

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
)

var versionCmd = &cli.Command{
	Name:  "version",
	Usage: "Prints the current version of WebKit",
	Action: cmdtools.Wrap(func(_ context.Context, input cmdtools.CommandInput) error {
		input.Printer().Print("Version: " + "TODO")
		return nil
	}),
}
