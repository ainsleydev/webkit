package secrets

import (
	"errors"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets/age"
)

func TestEncrypt(t *testing.T) {
	ctx := t.Context()
	ageIdentity, err := age.NewIdentity()
	require.NoError(t, err)

	tmpDir := t.TempDir()

	t.Run("Client Error", func(t *testing.T) {
		input := cmdtools.CommandInput{Command: GetCmd}

		err = Encrypt(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "age")
	})

	t.Run("SOPS Error", func(t *testing.T) {
		input := cmdtools.CommandInput{
			FS:      afero.NewBasePathFs(afero.NewOsFs(), tmpDir),
			BaseDir: tmpDir,
			Command: GetCmd,
			SOPSCache: &mockEncrypterDecrypter{
				err: errors.New("sops encrypt error"),
			},
		}

		err = Encrypt(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "sops encrypt error")
	})

	t.Run("Encrypts Successfully", func(t *testing.T) {
		t.Setenv(age.KeyEnvVar, ageIdentity.String())
		fs := afero.NewBasePathFs(afero.NewOsFs(), tmpDir)

		input := cmdtools.CommandInput{
			FS:      fs,
			BaseDir: tmpDir,
			Command: GetCmd,
		}
		err = CreateFiles(ctx, input)
		assert.NoError(t, err)

		content := "KEY: VALUE"
		path := "resources/secrets/production.yaml"

		err = afero.WriteFile(fs, path, []byte(content), os.ModePerm)
		require.NoError(t, err)

		err = Encrypt(ctx, input)
		assert.NoError(t, err)

		// This should trigger sops.ErrAlreadyEncrypted, but not
		// return any error.
		err = Encrypt(ctx, input)
		assert.NoError(t, err)
	})
}
