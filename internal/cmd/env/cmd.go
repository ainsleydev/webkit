package env

import (
	"github.com/urfave/cli/v3"
)

// Command defines the env commands for interacting and generating
// env file artifacts.
var Command = &cli.Command{
	Name:        "env",
	Usage:       "Manage environment variables",
	Description: "Command for working with the environment files defined in app.json",
	Commands: []*cli.Command{
		ScaffoldCmd,
		SyncCmd,
	},
}
