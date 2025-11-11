package infra

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/pkg/env"
)

var PlanCmd = &cli.Command{
	Name:   "plan",
	Usage:  "Generates an executive plan from the apps and resources defined in app.json",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "silent",
			Aliases: []string{"s"},
			Usage:   "Suppress informational output (only show Terraform output)",
		},
	},
	Action: cmdtools.Wrap(Plan),
}

func Plan(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()

	printer.Info("Generating executive plan from app definition")
	spinner := input.Spinner()

	// Filter definition to only include Terraform-managed items.
	appDef := input.AppDef()
	filtered, skipped := appDef.FilterTerraformManaged()

	// Show skipped items if any.
	if len(skipped.Apps) > 0 || len(skipped.Resources) > 0 {
		printer.Print("")
		printer.Info("The following items are not managed by Terraform:")
		if len(skipped.Apps) > 0 {
			printer.Print("  Apps:")
			for _, app := range skipped.Apps {
				printer.Print("    - " + app)
			}
		}
		if len(skipped.Resources) > 0 {
			printer.Print("  Resources:")
			for _, resource := range skipped.Resources {
				printer.Print("    - " + resource)
			}
		}
		printer.Print("")
	}

	// Use filtered definition for Terraform.
	tf, cleanup, err := initTerraformWithDefinition(ctx, input, filtered)
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

	// Write plan output directly to stdout (not through printer)
	fmt.Print(plan.Output) //nolint:forbidigo
	printer.Success("Plan generated, see console output")

	return nil
}
