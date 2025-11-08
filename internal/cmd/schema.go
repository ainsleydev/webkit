package cmd

import (
	"context"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
)

var schemaCmd = &cli.Command{
	Name:        "schema",
	Usage:       "Generate JSON schema for app.json",
	Description: "Generates a JSON schema file that can be used for IDE autocomplete and validation",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Output file path for the generated schema",
			Value:   "schema.json",
		},
	},
	Action: cmdtools.Wrap(schema),
}

func schema(ctx context.Context, input cmdtools.CommandInput) error {
	outputPath := input.Command.String("output")

	input.Printer().Info("Generating JSON schema...")

	// Generate schema
	schemaData, err := appdef.GenerateSchema()
	if err != nil {
		return errors.Wrap(err, "generating schema")
	}

	// Write to file
	err = afero.WriteFile(input.FS, outputPath, schemaData, 0o644)
	if err != nil {
		return errors.Wrap(err, "writing schema file")
	}

	input.Printer().Success("Schema generated successfully at: " + outputPath)

	return nil
}
