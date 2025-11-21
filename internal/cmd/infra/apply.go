package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/pkg/env"
)

var ApplyCmd = &cli.Command{
	Name:  "apply",
	Usage: "Creates or updates infrastructure based off the apps and resources defined in app.json",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "silent",
			Aliases: []string{"s"},
			Usage:   "Suppress informational output (only show Terraform output)",
		},
		&cli.BoolFlag{
			Name:  "refresh-only",
			Usage: "Sync Terraform state with actual infrastructure without making changes (uses 'terraform apply -refresh-only')",
		},
	},
	Action: cmdtools.Wrap(Apply),
}

func Apply(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()
	refreshOnly := input.Command.Bool("refresh-only")

	if refreshOnly {
		printer.Info("Syncing Terraform state with actual infrastructure (refresh-only mode)")
	} else {
		printer.Info("Generating executive plan from app definition")
	}
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

	if refreshOnly {
		printer.Println("Refreshing State...")
	} else {
		printer.Println("Applying Changes...")
	}
	spinner.Start()

	result, err := tf.Apply(ctx, env.Production, refreshOnly)
	if err != nil {
		// Write error output directly to stdout (not through printer)
		fmt.Print(result.Output) //nolint:forbidigo
		if refreshOnly {
			return errors.New("executing terraform apply -refresh-only")
		}
		return errors.New("executing terraform apply")
	}

	spinner.Stop()

	// Write output directly to stdout (not through printer)
	fmt.Print(result.Output) //nolint:forbidigo
	if refreshOnly {
		printer.Success("Refresh succeeded, state is now in sync with actual infrastructure")
	} else {
		printer.Success("Apply succeeded, see console output")
	}

	return nil
}
