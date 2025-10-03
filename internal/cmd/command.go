package cmd

import (
	"context"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/app"
	"github.com/ainsleydev/webkit/internal/scaffold"
)

// runCommand is the signature for command handlers. Each command
// should implement this function signature to run.
type runCommand func(ctx context.Context, input commandInput) error

// commandInput provides dependencies and context to command handlers.
type commandInput struct {
	FS        afero.Fs
	AppDef    app.Definition
	Command   *cli.Command
	Generator cgtools.Generator
}

// wrapCommand wraps a RunCommand to work with urfave/cli.
func wrapCommand(command runCommand) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		fs := afero.NewOsFs()
		input := commandInput{
			Command: c,
			AppDef:  app.Definition{},
			FS:      fs,
		}
		return command(ctx, input)
	}
}
