package infra

import (
	"context"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
	"github.com/ainsleydev/webkit/internal/secrets"
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

	printer.Println("Resolving Secrets...")
	spinner.Start()

	// Resolve all secrets from SOPS so we can pass them
	// to Terraform unmasked.
	resolveConfig := secrets.ResolveConfig{
		SOPSClient: input.SOPSClient(),
		BaseDir:    input.BaseDir,
	}

	// Ensure secrets are always re-encrypted, even if there's a panic or error
	defer secrets.EnsureEncrypted(resolveConfig)

	err := secrets.Resolve(ctx, appDef, resolveConfig)
	if err != nil {
		return nil, func() {}, err
	}

	spinner.Stop()
	printer.Println("Initializing Terraform...")
	spinner.Start()

	tf, err := newTerraform(ctx, appDef, input.Manifest)
	teardown := func() {
		tf.Cleanup()
	}
	if err != nil {
		return nil, teardown, err
	}

	if err = tf.Init(ctx); err != nil {
		return nil, teardown, err
	}

	spinner.Stop()

	return tf, teardown, nil
}
