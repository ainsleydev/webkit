package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/cmd/internal/operations"
)

var scaffoldCmd = &cli.Command{
	Name:   "scaffold",
	Hidden: true,
	Commands: []*cli.Command{
		{
			Name:   "code-style",
			Action: cmdtools.WrapCommand(operations.CreateCodeStyleFiles),
		},
		{
			Name:   "git",
			Action: cmdtools.WrapCommand(operations.CreateGitSettings),
		},
		{
			Name:   "package-json",
			Action: cmdtools.WrapCommand(operations.CreatePackageJson),
		},
		{
			Name:   "cicd",
			Action: cmdtools.WrapCommand(operations.CreateCICD),
		},
	},
}
