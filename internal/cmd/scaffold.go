package cmd

import (
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/cicd"
	"github.com/ainsleydev/webkit/internal/cmd/env"
	"github.com/ainsleydev/webkit/internal/cmd/files"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
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
			Action: cmdtools.Wrap(files.CodeStyle),
		},
		{
			Name:   "git",
			Usage:  "Generate Git and GitHub configuration files",
			Action: cmdtools.Wrap(files.CreateGitSettings),
		},
		{
			Name:   "package-json",
			Usage:  "Generate root package.json file",
			Action: cmdtools.Wrap(files.CreatePackageJson),
		},
		{
			Name:   "cicd",
			Usage:  "Generate GitHub Actions workflow files",
			Action: cmdtools.Wrap(cicd.CreatePRWorkflow),
		},
		{
			Name:   "pnpm-workspace",
			Usage:  "Generate pnpm-workspace file if there are any compatible apps.",
			Action: cmdtools.Wrap(files.CreatePNPMWorkspace),
		},
		{
			Name:   "turbo",
			Usage:  "Generate turbo.json file if there are any compatible apps.",
			Action: cmdtools.Wrap(files.CreateTurboJson),
		},
		{
			Name:   "docker-ignore",
			Usage:  "Generate .dockerignore files for every app defined in the definition.",
			Action: cmdtools.Wrap(files.DockerIgnore),
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
