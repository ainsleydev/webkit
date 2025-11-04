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
  webkit infra import --app web --id a1b2c3d4-e5f6-7890-abcd-ef1234567890

  # Import a DigitalOcean project
  webkit infra import --project --id 12345678-abcd-1234-5678-1234567890ab`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "resource",
			Aliases:  []string{"r"},
			Usage:    "Name of the resource in app.json to import (mutually exclusive with --app and --project)",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "app",
			Aliases:  []string{"a"},
			Usage:    "Name of the app in app.json to import (mutually exclusive with --resource and --project)",
			Required: false,
		},
		&cli.BoolFlag{
			Name:     "project",
			Aliases:  []string{"p"},
			Usage:    "Import a DigitalOcean project (mutually exclusive with --resource and --app)",
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

type (
	// importType represents what kind of item is being imported.
	importType int
)

const (
	importTypeResource importType = iota
	importTypeApp
	importTypeProject
)

func (t importType) String() string {
	switch t {
	case importTypeApp:
		return "app"
	case importTypeProject:
		return "project"
	default:
		return "resource"
	}
}

// importTarget represents the validated import target from CLI flags.
type importTarget struct {
	Type        importType
	Name        string // Empty for projects
	ID          string
	Environment env.Environment
}

// parseImportFlags validates and extracts the import target from CLI flags.
func parseImportFlags(cmd *cli.Command) (importTarget, error) {
	resourceName := cmd.String("resource")
	appName := cmd.String("app")
	isProject := cmd.Bool("project")

	// Count how many flags are set.
	flagsSet := []bool{resourceName != "", appName != "", isProject}
	count := 0
	for _, set := range flagsSet {
		if set {
			count++
		}
	}

	if count == 0 {
		return importTarget{}, fmt.Errorf("one of --resource, --app, or --project must be specified")
	}
	if count > 1 {
		return importTarget{}, fmt.Errorf("--resource, --app, and --project are mutually exclusive")
	}

	target := importTarget{
		ID:          cmd.String("id"),
		Environment: env.Environment(cmd.String("env")),
	}

	switch {
	case isProject:
		target.Type = importTypeProject
	case appName != "":
		target.Type = importTypeApp
		target.Name = appName
	default:
		target.Type = importTypeResource
		target.Name = resourceName
	}

	return target, nil
}

// Import executes the import operation for the specified resource or app.
func Import(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()
	spinner := input.Spinner()

	target, err := parseImportFlags(input.Command)
	if err != nil {
		return err
	}

	// Print what we're importing.
	if target.Type == importTypeProject {
		printer.Info(fmt.Sprintf("Importing DigitalOcean project with ID %q into %s environment",
			target.ID, target.Environment))
	} else {
		printer.Info(fmt.Sprintf("Importing %s %q with ID %q into %s environment",
			target.Type, target.Name, target.ID, target.Environment))
	}
	printer.Print("")

	// Initialise Terraform.
	tf, cleanup, err := initTerraform(ctx, input)
	defer cleanup()
	if err != nil {
		return err
	}

	printer.Print(fmt.Sprintf("Importing %s...", target.Type))
	spinner.Start()

	result, err := tf.Import(ctx, infra.ImportInput{
		ResourceName: target.Name,
		ResourceID:   target.ID,
		Environment:  target.Environment,
		IsApp:        target.Type == importTypeApp,
		IsProject:    target.Type == importTypeProject,
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
