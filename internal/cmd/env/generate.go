package env

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
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

// Generate creates a .env file for a specific app and environment.
func Generate(ctx context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()
	printer := input.Printer()
	spinner := input.Spinner()

	appName := input.Command.String("app")
	environmentStr := input.Command.String("environment")
	outputPath := input.Command.String("output")

	environment := env.Environment(environmentStr)

	var targetApp *appdef.App
	for _, app := range appDef.Apps {
		if app.Name == appName {
			targetApp = &app
			break
		}
	}

	if targetApp == nil {
		return fmt.Errorf("app '%s' not found in app.json", appName)
	}

	// Check if we need to fetch Terraform outputs (only if there are resource references).
	var tfOutputs *secrets.TerraformOutputProvider
	var err error
	if hasResourceReferences(appDef) {
		printer.Println("Fetching Terraform outputs...")
		spinner.Start()

		tfOutputs, err = fetchTerraformOutputs(ctx, input, environment)
		if err != nil {
			spinner.Stop()
			return errors.Wrap(err, "fetching terraform outputs")
		}

		spinner.Stop()
	}

	err = secrets.ResolveForEnvironment(ctx, appDef, environment, secrets.ResolveConfig{
		SOPSClient:      input.SOPSClient(),
		BaseDir:         input.BaseDir,
		TerraformOutput: tfOutputs,
	})
	if err != nil {
		return err
	}

	mergedApp := targetApp.MergeEnvironments(appDef.Shared.Env)

	vars, err := mergedApp.GetVarsForEnvironment(environment)
	if err != nil {
		return err
	}

	if len(vars) == 0 {
		printer.Warn(fmt.Sprintf("No environment variables defined for app '%s' in environment '%s'", appName, environment))
		return nil
	}

	err = writeMapToFile(writeArgs{
		Input:            input,
		Vars:             vars,
		App:              *targetApp,
		Environment:      environment,
		CustomOutputPath: outputPath,
	})
	if err != nil {
		return err
	}

	if outputPath == "" {
		outputPath = fmt.Sprintf("%s/.env%s", targetApp.Path, envSuffix(environment))
	}

	printer.Success(fmt.Sprintf("Generated env file for app '%s' (%s) at: %s", appName, environment, outputPath))

	return nil
}

// envSuffix returns the .env file suffix for an environment.
func envSuffix(environment env.Environment) string {
	if environment == env.Development {
		return ""
	}
	return fmt.Sprintf(".%s", environment)
}

// fetchTerraformOutputs fetches Terraform outputs for the specified environment.
// Returns a TerraformOutputProvider containing resource outputs.
func fetchTerraformOutputs(
	ctx context.Context,
	input cmdtools.CommandInput,
	environment env.Environment,
) (*secrets.TerraformOutputProvider, error) {
	tf, err := infra.NewTerraform(ctx, input.AppDef(), input.Manifest)
	if err != nil {
		return nil, errors.Wrap(err, "creating terraform manager")
	}

	if err := tf.Init(ctx); err != nil {
		return nil, errors.Wrap(err, "initialising terraform")
	}
	defer tf.Cleanup()

	result, err := tf.Output(ctx, environment)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving terraform outputs")
	}

	provider := make(secrets.TerraformOutputProvider)
	for resourceName, outputs := range result.Resources {
		for outputName, value := range outputs {
			key := secrets.OutputKey{
				Environment:  environment,
				ResourceName: resourceName,
				OutputName:   outputName,
			}
			provider[key] = value
		}
	}

	return &provider, nil
}
