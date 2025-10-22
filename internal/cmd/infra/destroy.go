package infra

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/pkg/env"
)

var DestroyCmd = &cli.Command{
	Name:   "destroy",
	Usage:  "Tears down infrastructure defined in app.json",
	Action: cmdtools.Wrap(Destroy),
}

func Destroy(ctx context.Context, input cmdtools.CommandInput) error {
	printer := input.Printer()
	spinner := input.Spinner()

	// Ask for confirmation before destroying
	if !confirm("Are you sure you want to destroy all resources? This action cannot be undone.") {
		printer.Warn("Destroy aborted by user.")
		return nil
	}

	tf, cleanup, err := initTerraform(ctx, input)
	if err != nil {
		return err
	}
	defer cleanup()

	printer.Println("Destroying Resources...")
	spinner.Start()

	destroyOutput, err := tf.Destroy(ctx, env.Production)
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
	fmt.Printf("%s [y/N]: ", prompt) //nolint
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}
