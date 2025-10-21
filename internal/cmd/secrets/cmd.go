// Package secrets provides commands for managing SOPS-encrypted secret files
// stored in the resources/sops directory.
package secrets

import (
	"github.com/urfave/cli/v3"
)

// Command defines the secret commands for generating, encrypting, decrypting, and
// manipulating SOPS-encrypted YAML files in resources/sops/{env}.
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
