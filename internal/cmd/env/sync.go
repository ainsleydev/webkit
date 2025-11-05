package env

import (
	"context"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
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
	spinner := input.Spinner()

	// Check if we need to fetch Terraform outputs (only if there are resource references).
	var tfOutputs *secrets.TerraformOutputProvider
	var err error
	if hasResourceReferences(appDef) {
		printer.Println("Fetching Terraform outputs...")
		spinner.Start()

		tfOutputs, err = fetchAllTerraformOutputs(ctx, input)
		if err != nil {
			spinner.Stop()
			return errors.Wrap(err, "fetching terraform outputs")
		}

		spinner.Stop()
	}

	err = secrets.Resolve(ctx, appDef, secrets.ResolveConfig{
		SOPSClient:      input.SOPSClient(),
		BaseDir:         input.BaseDir,
		TerraformOutput: tfOutputs,
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

// hasResourceReferences checks if the definition contains any environment
// variables with source="resource" that need Terraform outputs.
func hasResourceReferences(def *appdef.Definition) bool {
	// Check shared environment.
	hasResource := false
	def.Shared.Env.Walk(func(entry appdef.EnvWalkEntry) {
		if entry.Source == appdef.EnvSourceResource {
			hasResource = true
		}
	})
	if hasResource {
		return true
	}

	// Check app environments.
	for _, app := range def.Apps {
		app.Env.Walk(func(entry appdef.EnvWalkEntry) {
			if entry.Source == appdef.EnvSourceResource {
				hasResource = true
			}
		})
		if hasResource {
			return true
		}
	}

	return false
}

// fetchAllTerraformOutputs fetches Terraform outputs for all environments that have .env files.
func fetchAllTerraformOutputs(
	ctx context.Context,
	input cmdtools.CommandInput,
) (*secrets.TerraformOutputProvider, error) {
	provider := &secrets.TerraformOutputProvider{
		Outputs: make(map[env.Environment]map[string]map[string]any),
	}

	// Initialize Terraform manager once.
	tf, err := infra.NewTerraform(ctx, input.AppDef(), input.Manifest)
	if err != nil {
		return nil, errors.Wrap(err, "creating terraform manager")
	}

	// Initialize terraform (copies templates to temp dir).
	if err := tf.Init(ctx); err != nil {
		return nil, errors.Wrap(err, "initialising terraform")
	}
	defer tf.Cleanup()

	// Fetch outputs for all environments that need .env files.
	for _, environment := range environmentsWithDotEnv {
		result, err := tf.Output(ctx, environment)
		if err != nil {
			return nil, errors.Wrap(err, "retrieving terraform outputs for "+string(environment))
		}
		provider.Outputs[environment] = result.Resources
	}

	return provider, nil
}
