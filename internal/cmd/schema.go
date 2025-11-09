package cmd

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
)

var schemaCmd = &cli.Command{
	Name:        "schema",
	Usage:       "Generate JSON schema for app.json",
	Description: "Generates a JSON schema file that can be used for IDE autocomplete and validation. By default, outputs to .webkit/schema.json for easy IDE integration.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "Output file path for the generated schema",
			Value:   ".webkit/schema.json",
		},
		&cli.BoolFlag{
			Name:  "stdout",
			Usage: "Output schema to stdout instead of a file",
		},
	},
	Action: cmdtools.Wrap(schema),
}

func schema(ctx context.Context, input cmdtools.CommandInput) error {
	outputPath := input.Command.String("output")
	stdout := input.Command.Bool("stdout")

	input.Printer().Info("Generating JSON schema...")

	// Generate schema
	schemaData, err := appdef.GenerateSchema()
	if err != nil {
		return errors.Wrap(err, "generating schema")
	}

	// Output to stdout if requested
	if stdout {
		fmt.Println(string(schemaData))
		return nil
	}

	// Ensure parent directory exists
	if err := input.FS.MkdirAll(".webkit", 0o755); err != nil {
		return errors.Wrap(err, "creating .webkit directory")
	}

	// Write to file
	err = afero.WriteFile(input.FS, outputPath, schemaData, 0o644)
	if err != nil {
		return errors.Wrap(err, "writing schema file")
	}

	input.Printer().Success("Schema generated successfully at: " + outputPath)

	return nil
}
