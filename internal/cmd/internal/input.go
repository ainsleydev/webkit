package cmdtools

import (
	"context"

	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
)

// RunCommand is the signature for command handlers. Each command
// should implement this function signature to run.
type RunCommand func(ctx context.Context, input CommandInput) error

// CommandInput provides dependencies and context to command handlers.
type CommandInput struct {
	FS          afero.Fs
	Command     *cli.Command
	AppDefCache *appdef.Definition
}

// WrapCommand wraps a RunCommand to work with urfave/cli.
func WrapCommand(command RunCommand) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		// Let's temporarily use playground so we don't override any shit.
		fs := afero.NewBasePathFs(afero.NewOsFs(), "./internal/playground")
		input := CommandInput{
			Command: c,
			FS:      fs,
		}
		return command(ctx, input)
	}
}

// AppDef retrieves the main app manifest from the root
// of the project. Exits without it.
func (c *CommandInput) AppDef() *appdef.Definition {
	if c.AppDefCache != nil {
		return c.AppDefCache
	}

	read, err := appdef.Read(c.FS)
	if err != nil {
		Exit(err)
	}
	c.AppDefCache = read

	return read
}
