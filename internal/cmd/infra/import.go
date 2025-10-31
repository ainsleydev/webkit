package infra

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
	"github.com/ainsleydev/webkit/pkg/env"
)

// ImportCmd defines the command for importing existing infrastructure
// into Terraform state.
var ImportCmd = &cli.Command{
	Name:  "import",
	Usage: "Import existing infrastructure resources into Terraform state",
	Description: `Import allows you to bring existing cloud resources under Terraform management.

This is useful when:
  - Migrating from manually created infrastructure
  - Adopting webkit for an existing project
  - Recovering from state loss

Example:
  webkit infra import --resource db --id ca9f591d-f38h-462a-a5c6-5a8a74838081`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "resource",
			Aliases:  []string{"r"},
			Usage:    "Name of the resource in app.json to import",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "id",
			Usage:    "Provider-specific resource ID (e.g., DigitalOcean cluster ID)",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "env",
			Aliases: []string{"e"},
			Usage:   "Environment to import into (development, staging, production)",
			Value:   "production",
		},
	},
	Action: cmdtools.Wrap(Import),
}

// Import executes the import operation for the specified resource.
func Import(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()
	spinner := input.Spinner()

	resourceName := input.Command.String("resource")
	resourceID := input.Command.String("id")
	environment := env.Environment(input.Command.String("env"))

	printer.Info(fmt.Sprintf("Importing resource %q with ID %q into %s environment",
		resourceName, resourceID, environment))
	printer.Print("")

	// Initialise Terraform.
	tf, cleanup, err := initTerraform(ctx, input)
	defer cleanup()
	if err != nil {
		return err
	}

	printer.Print("Importing resources...")
	spinner.Start()

	result, err := tf.Import(ctx, infra.ImportInput{
		ResourceName: resourceName,
		ResourceID:   resourceID,
		Environment:  environment,
	})

	spinner.Stop()

	if err != nil {
		printer.Error("Import failed")
		if result.Output != "" {
			printer.Print(result.Output)
		}
		return nil
	}

	printer.Success(fmt.Sprintf("Successfully imported %d resource(s)", len(result.ImportedResources)))
	printer.Info("Imported Terraform resources:")
	for _, addr := range result.ImportedResources {
		printer.Print("  - " + addr)
	}

	printer.Info("Next steps:")
	printer.Print("  1. Run 'webkit infra plan' to verify the import")
	printer.Print("  2. If there are configuration differences, update app.json to match")
	printer.Print("  3. Run 'webkit infra apply' to finalise any adjustments")

	return nil
}
