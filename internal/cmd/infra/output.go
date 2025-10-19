package infra

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/pkg/env"
)

var OutputCmd = &cli.Command{
	Name:        "output",
	Usage:       "Retrieve Terraform outputs from a specific environment",
	Description: "Fetches infrastructure outputs for resources provisioned by Terraform",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "env",
			Usage:    "Environment to fetch outputs from (development, staging, production)",
			Aliases:  []string{"e"},
			Required: true,
		},
		&cli.StringFlag{
			Name:    "resource",
			Usage:   "Filter outputs to a specific resource name",
			Aliases: []string{"r"},
		},
		&cli.StringFlag{
			Name:    "app",
			Usage:   "Filter outputs to a specific app name",
			Aliases: []string{"a"},
		},
	},
	Action: cmdtools.Wrap(Output),
}

// Output retrieves Terraform outputs from the specified environment.
// Supports filtering by resource name or app name.
func Output(ctx context.Context, input cmdtools.CommandInput) error {
	cmd := input.Command
	printer := input.Printer()
	spinner := input.Spinner()

	envStr := cmd.String("env")
	resource := cmd.String("resource")
	app := cmd.String("app")

	tf, cleanup, err := initTerraform(ctx, input)
	if err != nil {
		return err
	}
	defer cleanup()

	printer.Println("Retrieving outputs...")
	spinner.Start()

	result, err := tf.Output(ctx, env.Environment(envStr))
	if err != nil {
		spinner.Stop()
		return errors.Wrap(err, "retrieving terraform outputs")
	}

	spinner.Stop()

	// Filter based on flags
	var output any
	switch {
	case resource != "":
		// Filter to specific resource
		resourceOutput, exists := result.Resources[resource]
		if !exists {
			return fmt.Errorf("resource '%s' not found in outputs", resource)
		}
		output = resourceOutput
		printer.Success(fmt.Sprintf("Outputs for resource '%s' from environment: %s\n", resource, envStr))
	case app != "":
		// Filter to specific app
		appOutput, exists := result.Apps[app]
		if !exists {
			return fmt.Errorf("app '%s' not found in outputs", app)
		}
		output = appOutput
		printer.Success(fmt.Sprintf("Outputs for app '%s' from environment: %s\n", app, envStr))
	default:
		// Show all outputs
		output = result
		printer.Success(fmt.Sprintf("All outputs from environment: %s\n", envStr))
	}

	// Pretty print JSON
	indent, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return errors.Wrap(err, "serializing terraform outputs")
	}
	fmt.Println(string(indent)) //nolint:forbidigo

	return nil
}
