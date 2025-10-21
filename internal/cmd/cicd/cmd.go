// Package cicd provides commands for managing GitHub Actions workflows and CI/CD configurations.
package cicd

import (
	"path/filepath"

	"github.com/urfave/cli/v3"
)

// Command defines the CI/CD commands for generating and managing
// GitHub workflow and action file artifacts.
var Command = &cli.Command{
	Name:        "cicd",
	Usage:       "Manage github workflows",
	Description: "Command for working with github CI/CD workflows",
	Commands: []*cli.Command{
		ActionsCmd,
		BackupCmd,
		PRCmd,
	},
}

var (
	actionsPath   = filepath.Join(".github", "actions")
	workflowsPath = filepath.Join(".github", "workflows")
)
