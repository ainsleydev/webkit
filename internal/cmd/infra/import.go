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
	Usage: "Import existing infrastructure resources or apps into Terraform state",
	Description: `Import allows you to bring existing cloud resources or apps under Terraform management.

This is useful when:
  - Migrating from manually created infrastructure
  - Adopting webkit for an existing project
  - Recovering from state loss

Examples:
  # Import a resource (database, storage, etc.)
  webkit infra import --resource db --id ca9f591d-f38h-462a-a5c6-5a8a74838081

  # Import an app
  webkit infra import --app web --id a1b2c3d4-e5f6-7890-abcd-ef1234567890`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "resource",
			Aliases:  []string{"r"},
			Usage:    "Name of the resource in app.json to import (mutually exclusive with --app)",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "app",
			Aliases:  []string{"a"},
			Usage:    "Name of the app in app.json to import (mutually exclusive with --resource)",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "id",
			Usage:    "Provider-specific resource/app ID (e.g., DigitalOcean cluster ID or app ID)",
			Required: true,
		},
		&cli.StringFlag{
			Name:    "env",
			Aliases: []string{"e"},
			Usage:   "Environment to import into (development, staging, production)",
			Value:   env.Production.String(),
		},
	},
	Action: cmdtools.Wrap(Import),
}

// Import executes the import operation for the specified resource or app.
func Import(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()
	spinner := input.Spinner()

	// Validate that exactly one of --resource or --app is provided
	resourceName := input.Command.String("resource")
	appName := input.Command.String("app")

	if resourceName == "" && appName == "" {
		return fmt.Errorf("either --resource or --app must be specified")
	}
	if resourceName != "" && appName != "" {
		return fmt.Errorf("--resource and --app are mutually exclusive, specify only one")
	}

	// Determine if we're importing an app or resource
	isApp := appName != ""
	name := resourceName
	if isApp {
		name = appName
	}

	resourceID := input.Command.String("id")
	environment := env.Environment(input.Command.String("env"))

	itemType := "resource"
	if isApp {
		itemType = "app"
	}

	printer.Info(fmt.Sprintf("Importing %s %q with ID %q into %s environment",
		itemType, name, resourceID, environment))
	printer.Print("")

	// Initialise Terraform.
	tf, cleanup, err := initTerraform(ctx, input)
	defer cleanup()
	if err != nil {
		return err
	}

	printer.Print(fmt.Sprintf("Importing %s...", itemType))
	spinner.Start()

	result, err := tf.Import(ctx, infra.ImportInput{
		ResourceName: name,
		ResourceID:   resourceID,
		Environment:  environment,
		IsApp:        isApp,
	})

	spinner.Stop()

	if err != nil {
		printer.Error("Import failed")
		if result.Output != "" {
			printer.Print(result.Output)
		}
		return err
	}

	printer.Success(fmt.Sprintf("Successfully imported %d Terraform resource(s)", len(result.ImportedResources)))
	printer.Info("Imported Terraform addresses:")
	for _, addr := range result.ImportedResources {
		printer.Print("  - " + addr)
	}

	printer.Info("Next steps:")
	printer.Print("  1. Run 'webkit infra plan' to verify the import")
	printer.Print("  2. If there are configuration differences, update app.json to match")
	printer.Print("  3. Run 'webkit infra apply' to finalise any adjustments")

	return nil
}
