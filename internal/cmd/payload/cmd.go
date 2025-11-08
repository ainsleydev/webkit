package payload

import (
	"github.com/urfave/cli/v3"
)

// Command defines the payload commands for managing Payload CMS projects.
var Command = &cli.Command{
	Name:        "payload",
	Usage:       "Manage Payload CMS projects",
	Description: "Commands for working with Payload CMS applications defined in app.json",
	Commands: []*cli.Command{
		BumpCmd,
	},
}
