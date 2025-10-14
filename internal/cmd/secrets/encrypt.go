package secrets

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
	"github.com/ainsleydev/webkit/pkg/env"
)

var EncryptCmd = &cli.Command{
	Name:        "encrypt",
	Usage:       "Encrypt secret files with SOPS",
	Description: "Encrypts all plaintext secret files in the secrets/ directory using SOPS and age.",
	Action:      cmdtools.Wrap(Encrypt),
}

func Encrypt(_ context.Context, input cmdtools.CommandInput) error {
	client := input.SOPSClient()

	input.Printer().Printf("Encrypting secret files...\n")

	var errs []error
	for _, e := range env.All {
		path := filepath.Join(input.BaseDir, secrets.FilePath, e.String()+".yaml")
		err := client.Encrypt(path)
		if errors.Is(err, sops.ErrAlreadyEncrypted) {
			continue
		} else if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	input.Printer().Success("Successfully encrypted secret files")

	return nil
}
