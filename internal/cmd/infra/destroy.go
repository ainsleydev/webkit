package infra

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/infra"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/pkg/env"
)

var DestroyCmd = &cli.Command{
	Name:   "destroy",
	Usage:  "Tears down infrastructure defined in app.json",
	Action: cmdtools.Wrap(Destroy),
}

func Destroy(ctx context.Context, input cmdtools.CommandInput) error {
	appDef := input.AppDef()
	printer := input.Printer()

	// Ask for confirmation before destroying
	if !confirm("Are you sure you want to destroy all resources? This action cannot be undone.") {
		printer.Warn("Destroy aborted by user.")
		return nil
	}

	// Resolve all secrets from SOPS so we can pass them
	// to Terraform unmasked.
	err := secrets.Resolve(ctx, appDef, secrets.ResolveConfig{
		SOPSClient: input.SOPSClient(),
		BaseDir:    input.BaseDir,
	})
	if err != nil {
		return err
	}

	printer.Info("Generating destruction plan from app definition")
	spinner := input.Spinner()

	terraform, err := infra.NewTerraform(ctx, appDef)
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
	printer.Println("Destroying Resources...")
	spinner.Start()

	destroyOutput, err := terraform.Destroy(ctx, env.Production)
	if err != nil {
		printer.Print(destroyOutput.Output)
		return errors.New("executing terraform destroy")
	}

	spinner.Stop()
	printer.Print(destroyOutput.Output)
	printer.Success("Destroy succeeded, see console output")

	return nil
}

func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}
