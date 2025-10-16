package infra

import (
	"github.com/urfave/cli/v3"
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
	},
}
