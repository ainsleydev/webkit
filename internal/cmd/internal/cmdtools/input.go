package cmdtools

import (
	"context"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/printer"
	"github.com/ainsleydev/webkit/internal/secrets/age"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
	"github.com/ainsleydev/webkit/pkg/env"
)

// RunCommand is the signature for command handlers. Each command
// should implement this function signature to run.
type RunCommand func(ctx context.Context, input CommandInput) error

// CommandInput provides dependencies and context to command handlers.
type CommandInput struct {
	FS          afero.Fs
	Command     *cli.Command
	AppDefCache *appdef.Definition
	BaseDir     string
	SOPSCache   sops.EncrypterDecrypter
	Manifest    *manifest.Tracker
	printer     *printer.Console
}

// Wrap wraps a RunCommand to work with urfave/cli.
func Wrap(command RunCommand) cli.ActionFunc {
	return func(ctx context.Context, c *cli.Command) error {
		fs := afero.NewOsFs()
		dir := "./"

		if env.AppEnvironment() == env.Development {
			// Let's temporarily use playground so we don't override any shit.
			path := "./internal/playground"
			fs = afero.NewBasePathFs(afero.NewOsFs(), path)
			dir = path
		}

		input := CommandInput{
			Command:  c,
			FS:       fs,
			BaseDir:  dir,
			Manifest: manifest.NewTracker(),
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

// Printer returns a new console writer to stdout.
func (c *CommandInput) Printer() *printer.Console {
	if c.printer == nil {
		c.printer = printer.New(os.Stdout)
	}
	return c.printer
}

// SOPSClient returns a cached sops.Client or initialises it
// by using an age provider.
func (c *CommandInput) SOPSClient() sops.EncrypterDecrypter {
	if c.SOPSCache != nil {
		return c.SOPSCache
	}
	prov, err := age.NewProvider()
	if err != nil {
		Exit(err)
	}
	c.SOPSCache = sops.NewClient(prov)
	return c.SOPSCache
}

func (c *CommandInput) Spinner() *spinner.Spinner {
	return spinner.New(spinner.CharSets[9], 100*time.Millisecond)
}
