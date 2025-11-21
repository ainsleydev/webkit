package infra

import (
	"context"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
)

var ExecCmd = &cli.Command{
	Name:  "exec",
	Usage: "Execute arbitrary Terraform commands",
	Description: `Execute any Terraform command with webkit's environment and configuration.

Examples:
  webkit infra exec -- state list
  webkit infra exec -- state show 'module.apps.peekaping_app.site["blog"]'
  webkit infra exec -- workspace list
  webkit infra exec -- show
  webkit infra exec -- validate`,
	ArgsUsage: "-- <terraform-command> [args...]",
	Action:    cmdtools.Wrap(Exec),
}

// Exec executes arbitrary Terraform commands with webkit's environment.
func Exec(ctx context.Context, input cmdtools.CommandInput) error {
	cmd := input.Command
	spinner := input.Spinner()

	args := cmd.Args().Slice()
	if len(args) == 0 {
		return errors.New("no terraform command provided (use -- before terraform args)")
	}

	appDef := input.AppDef()
	filtered, _ := appDef.FilterTerraformManaged()

	tf, cleanup, err := initTerraformWithDefinition(ctx, input, filtered)
	if err != nil {
		return err
	}
	defer cleanup()

	spinner.Start()

	// Execute terraform command with -chdir to access configs in temp dir
	// while keeping user's CWD for file path resolution.
	chdir := "-chdir=" + tf.WorkDir()
	tfCmd := exec.CommandContext(ctx, "terraform", append([]string{chdir}, args...)...)
	tfCmd.Stdout = os.Stdout
	tfCmd.Stderr = os.Stderr
	tfCmd.Stdin = os.Stdin
	tfCmd.Env = os.Environ()

	spinner.Stop()

	if err := tfCmd.Run(); err != nil {
		return errors.Wrap(err, "executing terraform command")
	}

	return nil
}
