package docs

import (
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmdtools"
)

// Command is the docs sub-command for managing agent documentation.
var Command = &cli.Command{
	Name:        "docs",
	Usage:       "Manage agent documentation",
	Description: "Generate and manage AGENTS.md documentation for AI coding assistants",
	Commands: []*cli.Command{
		{
			Name:        "generate",
			Usage:       "Generate AGENTS.md file",
			Description: "Generates AGENTS.md file from base template and optional custom content",
			Action:      cmdtools.Wrap(Generate),
		},
	},
}
