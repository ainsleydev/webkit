package sops

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/util/executil"
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

	t.Run("Decryption Fails", func(t *testing.T) {
		t.Parallel()

		file, teardown := setup(t)
		defer teardown()

		client, mem := newClient(&fakeProvider{})
		mem.AddStub("sops --decrypt", executil.Result{}, errors.New("unexpected failure"))

		got, err := DecryptFileToMap(client, file.Name())
		assert.Nil(t, got)
		assert.ErrorContains(t, err, "unexpected failure")
	})

	t.Run("Already Decrypted", func(t *testing.T) {
		t.Parallel()

		file, teardown := setup(t)
		defer teardown()
		_, err := file.WriteString("key: value\n")
		require.NoError(t, err)

		client, mem := newClient(&fakeProvider{})
		mem.AddStub("sops --decrypt --in-place", executil.Result{}, ErrNotEncrypted)

		data, err := DecryptFileToMap(client, file.Name())
		require.NoError(t, err)
		require.NotNil(t, data)
		assert.Equal(t, "value", data["key"])
	})

	t.Run("Read File Error", func(t *testing.T) {
		t.Parallel()

		client, mem := newClient(&fakeProvider{})
		mem.AddStub("sops --decrypt", executil.Result{}, nil)

		got, err := DecryptFileToMap(client, "wrong-path.yaml")
		assert.Nil(t, got)
		assert.ErrorContains(t, err, "failed to read sops file")
	})

	t.Run("Invalid YAML", func(t *testing.T) {
		t.Parallel()

		file, teardown := setup(t)
		defer teardown()

		_, err := file.WriteString("key: value\nunbalanced")
		require.NoError(t, err)

		client, mem := newClient(&fakeProvider{})
		mem.AddStub("sops --decrypt", executil.Result{
			Output: "key: value\nunbalanced",
		}, nil)

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

		client, mem := newClient(&fakeProvider{})
		mem.AddStub("sops --decrypt --in-place", executil.Result{
			Output: "",
		}, nil)

		mem.AddStub("sops --encrypt --age test-key --in-place", executil.Result{
			Output: "",
		}, nil)

		data, err := DecryptFileToMap(client, file.Name())
		require.NoError(t, err)
		require.NotNil(t, data)
		assert.Equal(t, "value", data["key"])
	})
}
