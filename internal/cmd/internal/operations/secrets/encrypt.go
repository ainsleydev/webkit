package secrets

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
	"github.com/ainsleydev/webkit/pkg/env"
)

func EncryptFiles(ctx context.Context, input cmdtools.CommandInput) error {
	fmt.Println("Encrypting secret files...")

	client, err := getSopsClient()
	if err != nil {
		return err
	}

	for _, e := range env.All {
		path := filepath.Join(input.BaseDir, secrets.FilePath, e+".yaml")
		err = client.Encrypt(path)
		if errors.Is(err, sops.ErrAlreadyEncrypted) {
			continue
		} else if err != nil {
			slog.ErrorContext(ctx, "Failed to encrypt secret file", "error", err, "file", path)
		}
	}

	fmt.Println("Successfully encrypted secret files")

	return nil
}
