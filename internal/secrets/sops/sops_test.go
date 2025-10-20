package sops

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/util/executil"
)

type fakeProvider struct {
	err error
}

func (f *fakeProvider) DecryptArgs() ([]string, error) {
	return []string{"--age", "test-key"}, f.err
}

func (f *fakeProvider) EncryptArgs() ([]string, error) {
	return []string{"--age", "test-key"}, f.err
}

func (f *fakeProvider) Environment() map[string]string {
	return map[string]string{"SOPS_AGE_KEY": "fake"}
}

func newClient(provider Provider) (*Client, *executil.MemRunner) {
	mem := executil.NewMemRunner()
	client := NewClient(provider)
	client.runner = mem
	return client, mem
}

func TestClient_Decrypt(t *testing.T) {
	t.Run("Already Decrypted", func(t *testing.T) {
		client, mem := newClient(&fakeProvider{})

		mem.AddStub("sops --decrypt", executil.Result{},
			errors.New("sops metadata not found"))

		err := client.Decrypt("file.yaml")
		assert.ErrorIs(t, err, ErrNotEncrypted)
	})

	t.Run("CLI Failure", func(t *testing.T) {
		client, mem := newClient(&fakeProvider{})

		mem.AddStub("sops --decrypt", executil.Result{
			Output: "unexpected error",
		}, fmt.Errorf("exit status 1"))

		err := client.Decrypt("file.yaml")
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "sops decrypt failed"))
	})

	t.Run("Success", func(t *testing.T) {
		client, mem := newClient(&fakeProvider{})

		mem.AddStub("sops --decrypt", executil.Result{
			Output: "decrypted",
		}, nil)

		err := client.Decrypt("file.yaml")
		assert.NoError(t, err)
	})
}

func TestClient_Encrypt(t *testing.T) {
	t.Run("Already Encrypted", func(t *testing.T) {
		client, mem := newClient(&fakeProvider{})

		mem.AddStub("sops --encrypt", executil.Result{},
			errors.New("contains a top-level entry called 'sops'"))

		err := client.Encrypt("file.yaml")
		assert.ErrorIs(t, err, ErrAlreadyEncrypted)
	})

	t.Run("Doesn't Error Empty File", func(t *testing.T) {
		client, mem := newClient(&fakeProvider{})

		mem.AddStub("sops --encrypt", executil.Result{},
			errors.New("it must contain at least one document"))

		err := client.Encrypt("file.yaml")
		assert.NoError(t, err)
	})

	t.Run("Provider Error", func(t *testing.T) {
		client, _ := newClient(&fakeProvider{err: fmt.Errorf("provider failure")})

		err := client.Encrypt("file.yaml")
		assert.Error(t, err)
		assert.EqualError(t, err, "provider failure")
	})

	t.Run("Encrypt CLI Failure", func(t *testing.T) {
		client, mem := newClient(&fakeProvider{})

		mem.AddStub("sops --encrypt", executil.Result{
			Output: "some error",
		}, fmt.Errorf("exit status 1"))

		err := client.Encrypt("file.yaml")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "sops encrypt failed")
	})

	t.Run("Success", func(t *testing.T) {
		client, mem := newClient(&fakeProvider{})

		mem.AddStub("sops --encrypt", executil.Result{
			Output: "encrypted",
		}, nil)

		err := client.Encrypt("file.yaml")
		assert.NoError(t, err)
	})
}
