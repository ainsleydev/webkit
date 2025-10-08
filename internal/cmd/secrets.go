package cmd

import (
	"github.com/urfave/cli/v3"
)

var secretsCmd = &cli.Command{
	Name:        "secrets",
	Usage:       "Manage SOPS-encrypted secret files",
	Description: "Commands for working with secret files defined in app.json",
	Commands: []*cli.Command{
		{
			Name:        "sync",
			Usage:       "Sync secret placeholders from app.json",
			Description: "Reads app.json and adds placeholder entries for all secrets with source: 'sops'",
			//Action:      cmdtools.Wrap(operations.SyncSecrets),
		},
		{
			Name:        "sync",
			Usage:       "Sync secret placeholders from app.json to SOPS files",
			Description: "Reads app.json and adds placeholder entries for all secrets with source: 'sops'.",
			//Action: cmdtools.Wrap(operations.SyncSecrets),
		},
		{
			Name:        "validate",
			Usage:       "Validate that all secrets from app.json exist in secret files",
			Description: "Ensures every secret referenced in app.json has a corresponding entry in SOPS files",
			//Action:      cmdtools.Wr(operations.ValidateSecrets),
		},
		{
			Name:  "encrypt",
			Usage: "Encrypt secret files with SOPS",
			Description: "Encrypts all plaintext secret files in the secrets/ directory using SOPS and age. " +
				"Requires age key to be configured in .sops.yaml. " +
				"Files are encrypted in-place.",
			//Action: cmdtools.WrapCommand(operations.EncryptSecrets),
		},
		{
			Name:  "decrypt",
			Usage: "Decrypt secret files with SOPS",
			Description: "Decrypts all encrypted secret files in the secrets/ directory using SOPS and age. " +
				"Requires age key to be available (SOPS_AGE_KEY env var or ~/.config/webkit/age.key). " +
				"Files are decrypted in-place. Use with caution - do not commit decrypted files.",
			//Action: cmdtools.WrapCommand(operations.DecryptSecrets),
		},
	},
}
