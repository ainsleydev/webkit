package docs

import (
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
)

// Command is the sub-command for managing AI assistant documentation.
var Command = &cli.Command{
	Name:        "agents",
	Usage:       "Manage AI assistant documentation",
	Description: "Agents and manage documentation files for AI coding assistants",
	Commands: []*cli.Command{
		{
			Name:        "generate",
			Usage:       "Agents AGENTS.md file",
			Description: "Generates AGENTS.md file from base template and optional custom content",
			Action:      cmdtools.Wrap(Agents),
		},
	},
}
