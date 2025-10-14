package secrets

import (
	"errors"
	"io"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets/age"
)

func TestDecrypt(t *testing.T) {
	ctx := t.Context()
	ageIdentity, err := age.NewIdentity()
	require.NoError(t, err)

	tmpDir := t.TempDir()

	t.Run("SOPS Error", func(t *testing.T) {
		input := cmdtools.CommandInput{
			FS:      afero.NewBasePathFs(afero.NewOsFs(), tmpDir),
			BaseDir: tmpDir,
			Command: GetCmd,
			SOPSCache: &mockEncrypterDecrypter{
				err: errors.New("sops decrypt error"),
			},
		}
		input.Printer().SetWriter(io.Discard)

		err = Decrypt(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "sops decrypt error")
	})

	t.Run("Decrypts Successfully", func(t *testing.T) {
		t.Setenv(age.KeyEnvVar, ageIdentity.String())

		input := cmdtools.CommandInput{
			FS:      afero.NewBasePathFs(afero.NewOsFs(), tmpDir),
			BaseDir: tmpDir,
			Command: GetCmd,
		}
		input.Printer().SetWriter(io.Discard)

		err = Scaffold(ctx, input)
		assert.NoError(t, err)

		err = Decrypt(ctx, input)
		assert.NoError(t, err)

		// This should trigger sops.ErrNotEncrypted, but not
		// return any error.
		err = Decrypt(ctx, input)
		assert.NoError(t, err)
	})
}
