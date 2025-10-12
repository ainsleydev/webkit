package secrets

import (
	"errors"
	"io"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
)

func TestEncrypt(t *testing.T) {
	ctx := t.Context()

	tmpDir := t.TempDir()

	t.Run("SOPS Error", func(t *testing.T) {
		input := cmdtools.CommandInput{
			FS:      afero.NewBasePathFs(afero.NewOsFs(), tmpDir),
			BaseDir: tmpDir,
			Command: GetCmd,
			SOPSCache: &mockEncrypterDecrypter{
				err: errors.New("sops encrypt error"),
			},
		}
		input.Printer().SetWriter(io.Discard)

		err := Encrypt(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "sops encrypt error")
	})

	t.Run("Encrypts Successfully", func(t *testing.T) {
		input, _ := setupEncryptedProdFile(t, `KEY: "1234"`)

		err := Encrypt(ctx, input)
		assert.NoError(t, err)

		// This should trigger sops.ErrAlreadyEncrypted, but not
		// return any error.
		err = Encrypt(ctx, input)
		assert.NoError(t, err)
	})
}
