package infra

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
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

	tf, cleanup, err := initTerraform(ctx, input)
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

	return nil
}
