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

func validate(ctx context.Context, input cmdtools.CommandInput) error {
	input.Printer().Info("Validating app.json...")

	// Load the app definition (this will parse and apply defaults)
	def := input.AppDef()

	// Run validation
	errs := def.Validate(input.FS)

	// Display results
	if errs == nil {
		input.Printer().Success("✓ Validation passed! No errors found.")
		return nil
	}

	// Display all errors
	input.Printer().Error(fmt.Sprintf("✗ Validation failed with %d error(s):", len(errs)))
	input.Printer().LineBreak()

	for i, err := range errs {
		input.Printer().Error(fmt.Sprintf("  %d. %s", i+1, err.Error()))
	}

	input.Printer().LineBreak()

	return errors.New("validation failed")
}
