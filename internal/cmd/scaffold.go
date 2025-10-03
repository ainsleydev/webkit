package cmd

import (
	"github.com/urfave/cli/v3"
)

var scaffoldCmd = &cli.Command{
	Name:   "scaffold",
	Hidden: true,
	Commands: []*cli.Command{
		createCodeStyleFilesCmd,
		createGithubSettingsCmd,
	},
}
