package infra

import (
	"context"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/pkg/env"
)

var PlanCmd = &cli.Command{
	Name:   "plan",
	Usage:  "Generates an executive plan from the apps and resources defined in app.json",
	Action: cmdtools.Wrap(Plan),
}

func Plan(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()

	printer.Info("Generating executive plan from app definition")
	spinner := input.Spinner()

	tf, cleanup, err := initTerraform(ctx, input)
	defer cleanup()
	if err != nil {
		return err
	}

	printer.Print("Making Plan...")
	spinner.Start()

	plan, err := tf.Plan(ctx, env.Production)
	if err != nil {
		return err
	}

	spinner.Stop()

	printer.Print(plan.Output)
	printer.Success("Plan generated, see console output")

	return nil
}
