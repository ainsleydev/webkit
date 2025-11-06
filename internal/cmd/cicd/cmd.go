package cicd

import (
	"path/filepath"

	"github.com/urfave/cli/v3"
)

// Command defines the env commands for interacting and generating
// github workflow/action file artifacts.
var Command = &cli.Command{
	Name:        "cicd",
	Usage:       "Manage github workflows",
	Description: "Command for working with github CI/CD workflows",
	Commands: []*cli.Command{
		ActionsCmd,
		BackupCmd,
		DeployAppCmd,
		DeployDigitalOceanVMCmd,
		DeployDigitalOceanContainerCmd,
		PRCmd,
		ReleaseCmd,
		VMMaintenanceCmd,
	},
}

var (
	actionsPath   = filepath.Join(".github", "actions")
	workflowsPath = filepath.Join(".github", "workflows")
)
