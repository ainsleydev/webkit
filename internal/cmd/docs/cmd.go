package docs

import (
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
)

// Command defines the docs commands for creating WebKit documentation
// from Webkit enabled projects such as AGENTS.md
var Command = &cli.Command{
	Name:        "docs",
	Usage:       "Manage Webkit documentation",
	Description: "TODO",
	Commands: []*cli.Command{
		{
			Name:        "agents",
			Usage:       "Agents AGENTS.md file",
			Description: "Generates AGENTS.md file from base template and optional custom content",
			Action:      cmdtools.Wrap(Agents),
		},
	},
}
