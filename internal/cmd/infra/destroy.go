package infra

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/pkg/env"
)

var DestroyCmd = &cli.Command{
	Name:   "destroy",
	Usage:  "Tears down infrastructure defined in app.json",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "silent",
			Aliases: []string{"s"},
			Usage:   "Suppress informational output (only show Terraform output)",
		},
	},
	Action: cmdtools.Wrap(Destroy),
}

func Destroy(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()
	spinner := input.Spinner()

	// Ask for confirmation before destroying
	if !confirm("Are you sure you want to destroy all resources? This action cannot be undone.") {
		printer.Warn("Destroy aborted by user.")
		return nil
	}

	// Filter definition to only include Terraform-managed items.
	appDef := input.AppDef()
	filtered, skipped := appDef.FilterTerraformManaged()

	// Show skipped items if any.
	if len(skipped.Apps) > 0 || len(skipped.Resources) > 0 {
		printer.Println("")
		printer.Info("The following items are not managed by Terraform:")
		if len(skipped.Apps) > 0 {
			printer.Println("  Apps:")
			for _, app := range skipped.Apps {
				printer.Println("    - " + app)
			}
		}
		if len(skipped.Resources) > 0 {
			printer.Println("  Resources:")
			for _, resource := range skipped.Resources {
				printer.Println("    - " + resource)
			}
		}
		printer.Println("")
	}

	// Use filtered definition for Terraform.
	tf, cleanup, err := initTerraformWithDefinition(ctx, input, filtered)
	if err != nil {
		return err
	}
	defer cleanup()

	printer.Println("Destroying Resources...")
	spinner.Start()

	destroyOutput, err := tf.Destroy(ctx, env.Production)
	if err != nil {
		printer.Print(destroyOutput.Output)
		return errors.New("executing terraform destroy")
	}

	spinner.Stop()
	printer.Print(destroyOutput.Output)
	printer.Success("Destroy succeeded, see console output")

	return nil
}

func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", prompt) //nolint
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}
