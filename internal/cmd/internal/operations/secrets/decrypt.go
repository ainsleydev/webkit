package secrets

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
	"github.com/ainsleydev/webkit/pkg/env"
)

var DecryptCmd = &cli.Command{
	Name:        "decrypt",
	Usage:       "Decrypt secret files with SOPS",
	Description: "Decrypts all encrypted secret files in the secrets/ directory using SOPS and age.",
	Action:      cmdtools.Wrap(Decrypt),
}

func Decrypt(_ context.Context, input cmdtools.CommandInput) error {
	client, err := input.SOPSClient()
	if err != nil {
		return err
	}

	fmt.Println("Decrypting secret files...")

	var errs []error
	for _, e := range env.All {
		path := filepath.Join(input.BaseDir, secrets.FilePath, e+".yaml")
		err = client.Decrypt(path)
		if errors.Is(err, sops.ErrNotEncrypted) {
			continue
		} else if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	fmt.Println("Successfully decrypted secret files")

	return nil
}
