package env

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/pkg/env"
)

var SyncCmd = &cli.Command{
	Name:        "sync",
	Usage:       "Sync secrets to env files from app.json",
	Description: "Reads app.json and adds creates or updates .env files in the relevant app directories",
	Action:      cmdtools.Wrap(Sync),
}

// Sync
func Sync(ctx context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()
	printer := input.Printer()

	err := secrets.Resolve(ctx, appDef, secrets.ResolveConfig{
		SOPSClient: input.SOPSClient(),
		BaseDir:    input.BaseDir,
	})
	if err != nil {
		return err
	}

	for _, app := range appDef.Apps {
		mergedApp := app.MergeEnvironments(appDef.Shared.Env)

		for _, enviro := range environmentsWithDotEnv {
			vars, err := getEnvironmentVars(mergedApp, enviro)
			if err != nil {
				return err
			}

			if len(vars) == 0 {
				continue
			}

			err = writeMapToFile(writeArgs{
				Input:       input,
				Vars:        vars,
				App:         app,
				Environment: enviro,
				IsScaffold:  false,
			})
			if err != nil {
				return err
			}
		}

		printer.Success("Successfully synced environment files for app: " + app.Name)
	}

	return nil
}
