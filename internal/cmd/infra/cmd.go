package infra

import (
	"context"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
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

func initTerraform(ctx context.Context, input cmdtools.CommandInput) (*infra.Terraform, func(), error) {
	appDef := input.AppDef()
	printer := input.Printer()
	spinner := input.Spinner()

	printer.Println("Resolving Secrets...")
	spinner.Start()

	// Resolve all secrets from SOPS so we can pass them
	// to Terraform unmasked.
	err := secrets.Resolve(ctx, appDef, secrets.ResolveConfig{
		SOPSClient: input.SOPSClient(),
		BaseDir:    input.BaseDir,
	})
	if err != nil {
		return nil, func() {}, err
	}

	spinner.Stop()
	printer.Println("Initializing Terraform...")
	spinner.Start()

	tf, err := infra.NewTerraform(ctx, appDef)
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
