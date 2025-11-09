package cmd

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
)

var validateCmd = &cli.Command{
	Name:        "validate",
	Usage:       "Validate app.json configuration",
	Description: "Validates the app.json file for correctness, including required fields, domain formats, paths, and environment variable references",
	Action:      cmdtools.Wrap(validate),
}

func validate(_ context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()

	printer.Info("Validating app.json...")
	printer.LineBreak()

	// Load the app definition (this will parse and apply defaults)
	def := input.AppDef()

	// Run validation
	errs := def.Validate(input.FS)

	// Display results
	if errs == nil {
		printer.Success("Validation passed! No errors found.")
		return nil
	}

	// Display all errors
	printer.Error(fmt.Sprintf("Validation failed with %d error(s):", len(errs)))
	printer.LineBreak()

	if len(errs) > 0 {
		for i, err := range errs {
			printer.Println(fmt.Sprintf("  %d. %s", i+1, err.Error()))
		}
	}

	return errors.New("validation failed")
}
