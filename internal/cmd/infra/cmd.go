package infra

import (
	"context"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/infra"
)

// Command defines the infra commands for provisioning and managing
// cloud infrastructure based on app.json definitions.
var Command = &cli.Command{
	Name:        "infra",
	Usage:       "Provision and manage cloud infrastructure",
	Description: "Commands for planning and applying infrastructure changes defined in app.json",
	Commands: []*cli.Command{
		PlanCmd,
		ApplyCmd,
		DestroyCmd,
	},
	Before: func(ctx context.Context, command *cli.Command) (context.Context, error) {
		_, err := infra.ParseTFEnvironment()
		if err != nil {
			// TODO, could make these look a bit sexier.
			return ctx, errors.Wrap(err, "must include infra variables in PATH")
		}
		return ctx, nil
	},
}
