package infra

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
)

var PlanCmd = &cli.Command{
	Name:   "plan",
	Usage:  "Generates an executive plan from the apps and resources defined in app.json",
	Action: cmdtools.Wrap(Plan),
}

func Plan(ctx context.Context, input cmdtools.CommandInput) error {
	_ = input.AppDef()
	printer := input.Printer()

	printer.Info("Generating executive plan from app definition")
	spinner := input.Spinner()

	terraform, err := infra.NewTerraform(ctx)
	if err != nil {
		return err
	}
	defer terraform.Cleanup()

	printer.Println("Initializing Terraform...")
	spinner.Start()

	if err = terraform.Init(ctx); err != nil {
		return err
	}

	spinner.Stop()
	printer.Println("Making Plan...")
	spinner.Start()

	if err = terraform.Plan(ctx); err != nil {
		return err
	}

	spinner.Stop()

	state, err := terraform.Show(context.Background())
	if err != nil {
		return err
	}

	fmt.Println(state)

	return nil
}
