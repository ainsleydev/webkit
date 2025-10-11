package secrets

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/internal/secrets/age"
	"github.com/ainsleydev/webkit/internal/secrets/sops"
	"github.com/ainsleydev/webkit/pkg/env"
)

func DecryptFiles(ctx context.Context, input cmdtools.CommandInput) error {
	fmt.Println("Decrypting secret files...")

	prov, err := age.NewProvider()
	if err != nil {
		return err
	}

	client := sops.NewClient(prov)

	for _, e := range []string{env.Development, env.Staging, env.Production} {
		path := filepath.Join(input.BaseDir, secrets.FilePath, e+".yaml")
		if err = client.Decrypt(path); err != nil {
			slog.ErrorContext(ctx, "Failed to decrypt secret file", "error", err, "file", path)
		}
	}

	fmt.Println("Successfully decrypted secret files")

	return nil
}
