package infra

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
	"github.com/ainsleydev/webkit/pkg/env"
)

var ApplyCmd = &cli.Command{
	Name:   "apply",
	Usage:  "Creates or updates infrastructure based off the apps and resources defined in app.json",
	Action: cmdtools.Wrap(Apply),
}

func Apply(ctx context.Context, input cmdtools.CommandInput) error {
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

	printer.Println("Applying Changes...")
	spinner.Start()

	plan, err := tf.Apply(ctx, env.Production)
	if err != nil {
		printer.Print(plan.Output)
		return errors.New("executing terraform apply")
	}

	spinner.Stop()

	printer.Print(plan.Output)
	printer.Success("Apply succeeded, see console output")

	// Auto-update postgres firewalls with app IPs
	if err := autoUpdateFirewalls(ctx, input, tf, filtered, env.Production); err != nil {
		// Log warning but don't fail the apply
		printer.Warning("Failed to auto-update postgres firewalls: " + err.Error())
	}

	return nil
}

// autoUpdateFirewalls automatically updates PostgreSQL firewall rules
// by discovering app IPs and re-applying terraform.
func autoUpdateFirewalls(ctx context.Context, input cmdtools.CommandInput, tf *infra.Terraform, appDef *appdef.Definition, environment env.Environment) error {
	printer := input.Printer()

	// Get terraform outputs to find app IPs
	outputs, err := tf.Output(ctx, environment)
	if err != nil {
		return err
	}

	// Create firewall updater
	tfEnv, err := infra.ParseTFEnvironment()
	if err != nil {
		return err
	}

	updater := infra.NewFirewallUpdater(tfEnv.DigitalOceanAPIKey)

	// Discover IPs and update firewall configs
	updates, err := updater.UpdateFirewalls(ctx, appDef, outputs, environment)
	if err != nil {
		return err
	}

	// If no updates were made, we're done
	if len(updates) == 0 {
		return nil
	}

	// Print summary of what we're updating
	printer.Print(infra.FormatUpdateSummary(updates))

	// Clear the terraform vars cache to force regeneration with updated configs
	tf.ClearVarsCache()

	// Re-apply terraform with updated firewall rules
	spinner := input.Spinner()
	spinner.Start()

	_, err = tf.Apply(ctx, environment)
	if err != nil {
		return errors.New("re-applying terraform with updated firewall rules")
	}

	spinner.Stop()
	printer.Success("Firewall rules updated successfully")

	return nil
}
