package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/cmd/internal/operations"
)

var scaffoldCmd = &cli.Command{
	Name: "scaffold",
	Commands: []*cli.Command{
		{
			Name:   "code-style",
			Action: cmdtools.Wrap(operations.CreateCodeStyleFiles),
		},
		{
			Name:   "git",
			Action: cmdtools.Wrap(operations.CreateGitSettings),
		},
		{
			Name:   "package-json",
			Action: cmdtools.Wrap(operations.CreatePackageJson),
		},
		{
			Name:   "cicd",
			Action: cmdtools.Wrap(operations.CreateCICD),
		},
	},
}
