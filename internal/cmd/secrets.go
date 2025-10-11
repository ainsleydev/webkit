package cmd

import (
	"fmt"
	"slices"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/cmd/internal/operations/secrets"
	"github.com/ainsleydev/webkit/pkg/env"
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
			Action:      cmdtools.Wrap(secrets.Sync),
		},
		{
			Name:        "validate",
			Usage:       "Validate that all secrets from app.json exist in secret files",
			Description: "Ensures every secret referenced in app.json has a corresponding entry in SOPS files",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "check-orphans",
					Usage:   "Report keys in SOPS files not referenced in app.json",
					Aliases: []string{"o"},
				},
				&cli.BoolFlag{
					Name:    "allow-encrypted",
					Usage:   "Attempt to validate encrypted files (requires SOPS/age access)",
					Aliases: []string{"e"},
				},
			},
			Action: cmdtools.Wrap(secrets.Validate),
		},
		{
			Name:        "encrypt",
			Usage:       "Encrypt secret files with SOPS",
			Description: "Encrypts all plaintext secret files in the secrets/ directory using SOPS and age.",
			Action:      cmdtools.Wrap(secrets.EncryptFiles),
		},
		{
			Name:        "decrypt",
			Usage:       "Decrypt secret files with SOPS",
			Description: "Decrypts all encrypted secret files in the secrets/ directory using SOPS and age.",
			Action:      cmdtools.Wrap(secrets.DecryptFiles),
		},
		{
			Name:        "get",
			Usage:       "Retrieve a secret from a specific environment",
			Description: "Fetches the value of a secret from the chosen environment (development, staging, production)",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "env",
					Usage:    "Environment to fetch the secret from (development, staging, production)",
					Aliases:  []string{"e"},
					Required: true,
					Validator: func(s string) error {
						if !slices.Contains(env.All, s) {
							return fmt.Errorf("invalid environment: %s", s)
						}
						return nil
					},
				},
				&cli.StringFlag{
					Name:     "key",
					Usage:    "The key/name of the secret to retrieve",
					Aliases:  []string{"k"},
					Required: true,
				},
			},
			Action: cmdtools.Wrap(secrets.Get),
		},
	},
}
