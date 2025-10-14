package secrets

import (
	"github.com/urfave/cli/v3"
)

// Command defines the secret commands for interacting, generating and
// manipulating the resources/sops/{env} SOPS encrypted YAML files.
var Command = &cli.Command{
	Name:        "secrets",
	Usage:       "Manage SOPS-encrypted secret files",
	Description: "Command for working with secret files defined in app.json",
	Commands: []*cli.Command{
		ScaffoldCmd,
		SyncCmd,
		EncryptCmd,
		DecryptCmd,
		GetCmd,
		ValidateCmd,
	},
}
