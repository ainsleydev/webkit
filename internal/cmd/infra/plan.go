package infra

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/pkg/env"
)

var PlanCmd = &cli.Command{
	Name:   "plan",
	Usage:  "Generates an executive plan from the apps and resources defined in app.json",
	Action: cmdtools.Wrap(Plan),
}

func Plan(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()

	printer.Info("Generating executive plan from app definition")
	spinner := input.Spinner()

	tf, cleanup, err := initTerraform(ctx, input)
	defer cleanup()
	if err != nil {
		return err
	}

	printer.Print("Making Plan...")
	spinner.Start()

	plan, err := tf.Plan(ctx, env.Production)
	if err != nil {
		return err
	}

	spinner.Stop()

	// Show skipped items if any.
	if len(plan.Skipped.Apps) > 0 || len(plan.Skipped.Resources) > 0 {
		printer.Print("")
		printer.Info("The following items are not managed by Terraform:")
		if len(plan.Skipped.Apps) > 0 {
			printer.Print("  Apps:")
			for _, app := range plan.Skipped.Apps {
				printer.Print("    - " + app)
			}
		}
		if len(plan.Skipped.Resources) > 0 {
			printer.Print("  Resources:")
			for _, resource := range plan.Skipped.Resources {
				printer.Print("    - " + resource)
			}
		}
		printer.Print("")
	}

	printer.Print(plan.Output)
	printer.Success("Plan generated, see console output")

	return nil
}
