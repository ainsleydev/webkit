package sops

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func fakeCmdOutput(output string, code int) *exec.Cmd {
	script := fmt.Sprintf(`echo "%s"; exit %d`, output, code)
	return exec.Command("bash", "-c", script)
}

func newClient(t *testing.T, provider Provider, cmd *exec.Cmd) *Client {
	t.Helper()

	client := NewClient(provider)
	require.NotNil(t, client)

	client.exec = func(ctx context.Context, name string, arg ...string) *exec.Cmd {
		return cmd
	}

	return client
}

func TestClient_Encrypt(t *testing.T) {
	t.Parallel()

	t.Run("Already Encrypted", func(t *testing.T) {
		t.Parallel()

		client := newClient(t,
			&fakeProvider{},
			fakeCmdOutput("contains a top-level entry called 'sops', with error", 1),
		)

		err := client.Encrypt("file.yaml")
		fmt.Println(err)
		assert.ErrorIs(t, err, ErrAlreadyEncrypted)
	})

	t.Run("Provider Error", func(t *testing.T) {
		t.Parallel()

		client := newClient(t,
			&fakeProvider{err: fmt.Errorf("provider failure")},
			nil,
		)

		err := client.Encrypt("file.yaml")
		assert.Error(t, err)
		assert.EqualError(t, err, "provider failure")
	})

	t.Run("Encrypt CLI Failure", func(t *testing.T) {
		t.Parallel()

		client := newClient(t,
			&fakeProvider{},
			fakeCmdOutput("some error", 1),
		)

		err := client.Encrypt("file.yaml")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "sops encrypt failed")
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		client := newClient(t,
			&fakeProvider{},
			fakeCmdOutput("encrypted", 0),
		)

		err := client.Encrypt("file.yaml")
		assert.NoError(t, err)
	})
}

func TestClient_Decrypt(t *testing.T) {
	t.Parallel()

	t.Run("Already Decrypted", func(t *testing.T) {
		t.Parallel()

		client := newClient(t,
			&fakeProvider{},
			fakeCmdOutput("sops metadata not found", 1),
		)

		err := client.Decrypt("file.yaml")
		assert.ErrorIs(t, err, ErrNotEncrypted)
	})

	t.Run("CLI Failure", func(t *testing.T) {
		t.Parallel()

		client := newClient(t,
			&fakeProvider{},
			fakeCmdOutput("unexpected error", 1),
		)

		err := client.Decrypt("file.yaml")
		require.Error(t, err)
		assert.True(t, strings.Contains(err.Error(), "sops decrypt failed"))
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		client := newClient(t,
			&fakeProvider{},
			fakeCmdOutput("decrypted", 0),
		)

		err := client.Decrypt("file.yaml")
		assert.NoError(t, err)
	})
}
