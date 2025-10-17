package infra

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
	"github.com/ainsleydev/webkit/pkg/env"
)

var PlanCmd = &cli.Command{
	Name:   "plan",
	Usage:  "Generates an executive plan from the apps and resources defined in app.json",
	Action: cmdtools.Wrap(Plan),
}

func Plan(ctx context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()
	printer := input.Printer()

	printer.Info("Generating executive plan from app definition")
	spinner := input.Spinner()

	terraform, err := infra.NewTerraform(ctx)
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
	printer.Println("Making Plan...")
	spinner.Start()

	plan, err := terraform.Plan(ctx, env.Production, appDef)
	if err != nil {
		return err
	}

	spinner.Stop()

	if plan.ResourceChanges != nil {
		printer.Printf("\nResource Changes:\n")
		for _, rc := range plan.ResourceChanges {
			printer.Printf("  - %s (%s): %v\n", rc.Address, rc.Type, rc.Change.Actions)
		}
	}

	if plan.OutputChanges != nil {
		printer.Printf("\nOutput Changes:\n")
		for name, change := range plan.OutputChanges {
			printer.Printf("  - %s: %v -> %v\n", name, change.Before, change.After)
		}
	}

	return nil
}
