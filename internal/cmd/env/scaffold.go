package env

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
)

var ScaffoldCmd = &cli.Command{
	Name:   "scaffold",
	Usage:  "Generate empty env files for every app.",
	Action: cmdtools.Wrap(Scaffold),
}

// Scaffold generates the blank .env files located for every
// app defined in the definition. Prepends the WebKit notice.
func Scaffold(_ context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()

	for _, env := range environmentsWithDotEnv {
		for _, app := range appDef.Apps {
			err := writeMapToFile(writeArgs{
				FS:          input.FS,
				Vars:        nil, // Just generate the notice.
				App:         app,
				Environment: env,
				Manifest:    input.Manifest,
				IsScaffold:  true,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
