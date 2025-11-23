package secrets

import (
	"errors"
	"io"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/state/manifest"
	"github.com/ainsleydev/webkit/internal/mocks"
	"github.com/ainsleydev/webkit/internal/secrets/age"
)

func TestDecrypt(t *testing.T) {
	t.Run("SOPS Error", func(t *testing.T) {
		input, _ := setup(t, &appdef.Definition{})

		ctrl := gomock.NewController(t)
		mock := mocks.NewMockEncrypterDecrypter(ctrl)
		mock.EXPECT().
			Decrypt(gomock.Any()).
			Return(errors.New("sops decrypt error")).
			Times(3)

		input.SOPSCache = mock

		err := Decrypt(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "sops decrypt error")
	})

	t.Run("Decrypts Successfully", func(t *testing.T) {
		ageIdentity, err := age.NewIdentity()
		require.NoError(t, err)

		tmpDir := t.TempDir()
		t.Setenv(age.KeyEnvVar, ageIdentity.String())

		input := cmdtools.CommandInput{
			FS:       afero.NewBasePathFs(afero.NewOsFs(), tmpDir),
			BaseDir:  tmpDir,
			Command:  GetCmd,
			Manifest: manifest.NewTracker(),
		}
		input.Printer().SetWriter(io.Discard)

		err = Scaffold(t.Context(), input)
		assert.NoError(t, err)

		err = Decrypt(t.Context(), input)
		assert.NoError(t, err)

		// This should trigger sops.ErrNotEncrypted, but not
		// return any error.
		err = Decrypt(t.Context(), input)
		assert.NoError(t, err)
	})
}
