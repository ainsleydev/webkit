package secrets

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/mocks"
)

func TestEncrypt(t *testing.T) {
	t.Run("SOPS Error", func(t *testing.T) {
		input, _ := setup(t, &appdef.Definition{})

		ctrl := gomock.NewController(t)
		mock := mocks.NewMockEncrypterDecrypter(ctrl)
		mock.EXPECT().
			Encrypt(gomock.Any()).
			Return(errors.New("sops encrypt error")).
			Times(3)

		input.SOPSCache = mock

		err := Encrypt(t.Context(), input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "sops encrypt error")
	})

	t.Run("Encrypts Successfully", func(t *testing.T) {
		input, _ := setupEncryptedProdFile(t, `KEY: "1234"`)

		err := Encrypt(t.Context(), input)
		assert.NoError(t, err)

		// This should trigger sops.ErrAlreadyEncrypted, but not
		// return any error.
		err = Encrypt(t.Context(), input)
		assert.NoError(t, err)
	})
}
