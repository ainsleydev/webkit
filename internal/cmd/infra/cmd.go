package infra

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/pkg/env"
)

// Command defines the infra commands for provisioning and managing
// cloud infrastructure based on app.json definitions.
var Command = &cli.Command{
	Name:        "infra",
	Usage:       "Provision and manage cloud infrastructure",
	Description: "Commands for planning and applying infrastructure changes defined in app.json",
	Commands: []*cli.Command{
		PlanCmd,
		ApplyCmd,
		DestroyCmd,
		OutputCmd,
		ImportCmd,
	},
	Before: func(ctx context.Context, command *cli.Command) (context.Context, error) {
		_, err := infra.ParseTFEnvironment()
		if err != nil {
			// TODO, could make these look a bit sexier.
			return ctx, errors.Wrap(err, "must include infra variables in PATH")
		}
		return ctx, nil
	},
}

var newTerraform = infra.NewTerraform

func initTerraform(ctx context.Context, input cmdtools.CommandInput) (infra.Manager, func(), error) {
	return initTerraformWithDefinition(ctx, input, input.AppDef())
}

func initTerraformWithDefinition(ctx context.Context, input cmdtools.CommandInput, appDef *appdef.Definition) (infra.Manager, func(), error) {
	printer := input.Printer()
	spinner := input.Spinner()

	spinner.Stop()
	if !input.Silent {
		printer.Println("Initializing Terraform...")
	}
	spinner.Start()

	tf, err := newTerraform(ctx, appDef, input.Manifest)
	teardown := func() {
		if tf != nil {
			tf.Cleanup()
		}
	}
	if err != nil {
		spinner.Stop()
		return nil, teardown, err
	}

	if err = tf.Init(ctx); err != nil {
		spinner.Stop()
		return nil, teardown, err
	}

	spinner.Stop()

	if !input.Silent {
		printer.Println("Resolving Secrets...")
	}
	spinner.Start()

	// Check if we need to fetch Terraform outputs (only if there are resource references).
	var tfOutputs *secrets.TerraformOutputProvider
	if hasResourceReferences(appDef) {
		if !input.Silent {
			printer.Println("Fetching Terraform outputs...")
		}
		spinner.Start()

		tfOutputs, err = fetchTerraformOutputs(ctx, tf, env.Production)
		if err != nil {
			spinner.Stop()
		}

		spinner.Stop()
	}

	// Resolve all secrets from SOPS so we can pass them
	// to Terraform unmasked.
	resolveConfig := secrets.ResolveConfig{
		SOPSClient:      input.SOPSClient(),
		BaseDir:         input.BaseDir,
		TerraformOutput: tfOutputs,
	}

	// Ensure secrets are always re-encrypted, even if there's a panic or error
	defer func() {
		for _, e := range []env.Environment{env.Development, env.Staging, env.Production} {
			_ = resolveConfig.SOPSClient.Encrypt(filepath.Join(resolveConfig.BaseDir, secrets.FilePathFromEnv(e)))
		}
	}()

	err = secrets.Resolve(ctx, appDef, resolveConfig)
	if err != nil {
		spinner.Stop()
		return nil, func() {}, err
	}

	return tf, teardown, nil
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

// fetchTerraformOutputs fetches Terraform outputs for the specified environment.
// Returns a TerraformOutputProvider containing resource outputs.
//
// This function uses an existing, initialized Terraform instance.
// See also: env/cmd.go:fetchTerraformOutputs for a similar function that manages
// the full Terraform lifecycle.
func fetchTerraformOutputs(
	ctx context.Context,
	tf *infra.Terraform,
	environment env.Environment,
) (*secrets.TerraformOutputProvider, error) {
	result, err := tf.Output(ctx, environment)
	if err != nil {
		return nil, errors.Wrap(err, "retrieving terraform outputs")
	}

	provider := secrets.TransformOutputs(result, environment)
	return &provider, nil
}
