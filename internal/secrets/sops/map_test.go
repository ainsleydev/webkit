package sops

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecryptFileToMap(t *testing.T) {
	t.Parallel()

	setup := func(t *testing.T) (*os.File, func()) {
		t.Helper()

		file, err := os.CreateTemp("", "secrets-*.yaml")
		require.NoError(t, err)

		return file, func() {
			require.NoError(t, file.Close())
			require.NoError(t, os.Remove(file.Name()))
		}
	}

	t.Run("Decryption fails", func(t *testing.T) {
		t.Parallel()

		file, teardown := setup(t)
		defer teardown()

		client := newClient(t,
			&fakeProvider{},
			fakeCmdOutput("sops metadata not found", 1),
		)

		got, err := DecryptFileToMap(client, file.Name())
		assert.Nil(t, got)
		assert.ErrorIs(t, err, ErrNotEncrypted)
	})

	t.Run("Read File Error", func(t *testing.T) {
		t.Parallel()

		client := newClient(t,
			&fakeProvider{},
			fakeCmdOutput("", 0),
		)

		toMap, err := DecryptFileToMap(client, "wrong")
		assert.Nil(t, toMap)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to read sops file")
	})

	t.Run("Invalid YAML", func(t *testing.T) {
		t.Parallel()

		file, teardown := setup(t)
		defer teardown()

		// Write invalid YAML content
		_, err := file.WriteString("key: value\nunbalanced")
		require.NoError(t, err)

		client := newClient(t,
			&fakeProvider{},
			fakeCmdOutput("", 0),
		)

		got, err := DecryptFileToMap(client, file.Name())
		assert.Nil(t, got)
		assert.ErrorContains(t, err, "failed to parse sops content")
	})

	t.Run("Successful Decryption", func(t *testing.T) {
		t.Parallel()

		file, teardown := setup(t)
		defer teardown()

		_, err := file.WriteString("key: value\n")
		require.NoError(t, err)

		// Fake client just returns success
		client := newClient(t,
			&fakeProvider{},
			fakeCmdOutput("", 0),
		)

		data, err := DecryptFileToMap(client, file.Name())
		require.NoError(t, err)
		require.NotNil(t, data)

		assert.Equal(t, "value", data["key"])
	})
}
