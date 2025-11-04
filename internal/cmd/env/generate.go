package env

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/pkg/env"
)

var GenerateCmd = &cli.Command{
	Name:        "generate",
	Usage:       "Generate env file for a specific app and environment",
	Description: "Generates a .env file for a specific app and environment with a custom output path. Useful for VM deployments.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "app",
			Usage:    "Name of the app to generate env file for",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "environment",
			Aliases:  []string{"env"},
			Usage:    "Target environment (development, staging, production)",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "output",
			Usage: "Output path for the .env file (defaults to {app.Path}/.env.{environment})",
		},
	},
	Action: cmdtools.Wrap(Generate),
}

type generateArgs struct {
	AppName     string
	Environment string
	OutputPath  string
}

// Generate creates a .env file for a specific app and environment.
func Generate(ctx context.Context, input cmdtools.CommandInput) error {
	args := generateArgs{
		AppName:     input.Command.String("app"),
		Environment: input.Command.String("environment"),
		OutputPath:  input.Command.String("output"),
	}

	return generateEnvFile(ctx, input, args)
}

// generateEnvFile contains the core logic for generating env files.
func generateEnvFile(ctx context.Context, input cmdtools.CommandInput, args generateArgs) error {
	appDef := input.AppDef()
	printer := input.Printer()

	environment, err := env.Parse(args.Environment)
	if err != nil {
		return fmt.Errorf("invalid environment: %w", err)
	}

	var targetApp *appdef.App
	for _, app := range appDef.Apps {
		if app.Name == args.AppName {
			targetApp = &app
			break
		}
	}

	if targetApp == nil {
		return fmt.Errorf("app '%s' not found in app.json", args.AppName)
	}

	err = secrets.Resolve(ctx, appDef, secrets.ResolveConfig{
		SOPSClient: input.SOPSClient(),
		BaseDir:    input.BaseDir,
	})
	if err != nil {
		return err
	}

	mergedApp := targetApp.MergeEnvironments(appDef.Shared.Env)

	vars, err := getEnvironmentVars(mergedApp, environment)
	if err != nil {
		return err
	}

	if len(vars) == 0 {
		printer.Warn(fmt.Sprintf("No environment variables defined for app '%s' in environment '%s'", args.AppName, environment))
		return nil
	}

	err = writeMapToFile(writeArgs{
		Input:            input,
		Vars:             vars,
		App:              *targetApp,
		Environment:      environment,
		CustomOutputPath: args.OutputPath,
	})
	if err != nil {
		return err
	}

	outputPath := args.OutputPath
	if outputPath == "" {
		outputPath = fmt.Sprintf("%s/.env%s", targetApp.Path, envSuffix(environment))
	}

	printer.Success(fmt.Sprintf("Generated env file for app '%s' (%s) at: %s", args.AppName, environment, outputPath))

	return nil
}

// envSuffix returns the .env file suffix for an environment.
func envSuffix(environment env.Environment) string {
	if environment == env.Development {
		return ""
	}
	return fmt.Sprintf(".%s", environment)
}
