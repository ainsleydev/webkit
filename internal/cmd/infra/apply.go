package infra

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/pkg/env"
)

var ApplyCmd = &cli.Command{
	Name:   "apply",
	Usage:  "Creates or updates infrastructure based off the apps and resources defined in app.json",
	Action: cmdtools.Wrap(Apply),
}

func Apply(ctx context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()
	printer := input.Printer()

	printer.Info("Generating executive plan from app definition")
	spinner := input.Spinner()

	// Resolve all secrets from SOPS so we can pass them
	// to Terraform unmasked.
	err := secrets.Resolve(ctx, appDef, secrets.ResolveConfig{
		SOPSClient: input.SOPSClient(),
	})
	if err != nil {
		return err
	}

	terraform, err := infra.NewTerraform(ctx, appDef)
	if err != nil {
		return err
	}
	defer terraform.Cleanup()

	printer.Println("Initializing Terraform...")
	spinner.Start()

	if err = terraform.Init(ctx); err != nil {
		return err
	}

	spinner.Stop()
	printer.Println("Applying Changes...")
	spinner.Start()

	plan, err := terraform.Apply(ctx, env.Production)
	if err != nil {
		printer.Print(plan.Output)
		return errors.New("executing terraform apply")
	}

	spinner.Stop()

	printer.Print(plan.Output)
	printer.Success("Apply succeeded, see console output")

	return nil
}
