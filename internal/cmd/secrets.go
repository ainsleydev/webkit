package cmd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/cmd/internal/operations"
	"github.com/ainsleydev/webkit/internal/secrets/age"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
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
			Action:      cmdtools.Wrap(operations.SecretsSync),
		},
		{
			Name:        "validate",
			Usage:       "Validate that all secrets from app.json exist in secret files",
			Description: "Ensures every secret referenced in app.json has a corresponding entry in SOPS files",
		},
		{
			Name:        "encrypt",
			Usage:       "Encrypt secret files with SOPS",
			Description: "Encrypts all plaintext secret files in the secrets/ directory using SOPS and age.",
			Action: cmdtools.Wrap(func(ctx context.Context, input cmdtools.CommandInput) error {
				fmt.Println("Encrypting secret files...")

				prov, err := age.NewProvider()
				if err != nil {
					return err
				}

				path := filepath.Join(input.BaseDir, "resources", "secrets", "production.yaml")
				return sops.NewClient(prov).Encrypt(path)
			}),
		},
		{
			Name:        "decrypt",
			Usage:       "Decrypt secret files with SOPS",
			Description: "Decrypts all encrypted secret files in the secrets/ directory using SOPS and age.",
			Action: cmdtools.Wrap(func(ctx context.Context, input cmdtools.CommandInput) error {
				fmt.Println("Decrypting secret files...")

				prov, err := age.NewProvider()
				if err != nil {
					return err
				}

				path := filepath.Join(input.BaseDir, "resources", "secrets", "production.yaml")
				return sops.NewClient(prov).Decrypt(path)
			}),
		},
	},
}
