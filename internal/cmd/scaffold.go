package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/env"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/cmd/internal/operations"
	"github.com/ainsleydev/webkit/internal/cmd/secrets"
)

var scaffoldCmd = &cli.Command{
	Name:        "scaffold",
	Usage:       "Scaffold individual project components",
	Description: "Generate standalone project files without modifying existing generated templates.",
	Commands: []*cli.Command{
		{
			Name:   "code-style",
			Usage:  "Generate code style configuration files",
			Action: cmdtools.Wrap(operations.CreateCodeStyleFiles),
		},
		{
			Name:   "git",
			Usage:  "Generate Git and GitHub configuration files",
			Action: cmdtools.Wrap(operations.CreateGitSettings),
		},
		{
			Name:   "package-json",
			Usage:  "Generate root package.json file",
			Action: cmdtools.Wrap(operations.CreatePackageJson),
		},
		{
			Name:   "cicd",
			Usage:  "Generate GitHub Actions workflow files",
			Action: cmdtools.Wrap(operations.CreateCICD),
		},
		{
			Name:   "pnpm-workspace",
			Usage:  "Generate pnpm-workspace file if there are any compatible apps.",
			Action: cmdtools.Wrap(operations.CreatePNPMWorkspace),
		},
		{
			Name:   "turbo",
			Usage:  "Generate turbo.json file if there are any compatible apps.",
			Action: cmdtools.Wrap(operations.CreateTurboJson),
		},
		{
			Name:   "secrets",
			Usage:  "Generate empty SOPS secret files and configuration",
			Action: cmdtools.Wrap(secrets.Scaffold),
		},
		{
			Name:   "env",
			Usage:  "TODO",
			Action: cmdtools.Wrap(env.Scaffold),
		},
	},
}
