package infra

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/pkg/env"
)

var RefreshCmd = &cli.Command{
	Name:  "refresh",
	Usage: "Syncs Terraform state with actual infrastructure without making changes",
	Description: `Refresh reads the current settings from all managed remote resources and updates
the Terraform state to match. This uses 'terraform apply -refresh-only', which is the
modern replacement for the deprecated 'terraform refresh' command.

This is useful when resources have been modified outside of Terraform, or when provider
bugs cause state inconsistencies (e.g., the peekaping provider's tag_ids/notification_ids issues).

Note: Refresh does not modify infrastructure; it only updates the state file to reflect reality.`,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "silent",
			Aliases: []string{"s"},
			Usage:   "Suppress informational output (only show Terraform output)",
		},
	},
	Action: cmdtools.Wrap(Refresh),
}

func Refresh(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()

	printer.Info("Syncing Terraform state with actual infrastructure")
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

	printer.Println("Refreshing State...")
	spinner.Start()

	result, err := tf.Refresh(ctx, env.Production)
	if err != nil {
		// Write error output directly to stdout (not through printer)
		fmt.Print(result.Output) //nolint:forbidigo
		return errors.New("executing terraform refresh")
	}

	spinner.Stop()

	// Write refresh output directly to stdout (not through printer)
	fmt.Print(result.Output) //nolint:forbidigo
	printer.Success("Refresh succeeded, state is now in sync with actual infrastructure")

	return nil
}
