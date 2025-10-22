package age

import (
	"testing"

	"filippo.io/age"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProvider(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	require.NoError(t, err)

	t.Run("Private Key Error", func(t *testing.T) {
		t.Setenv(KeyEnvVar, "invalid-key-format")

		provider, err := NewProvider()
		assert.Error(t, err)
		assert.Nil(t, provider)
		assert.Contains(t, err.Error(), "invalid age key format")
	})

	t.Run("Success", func(t *testing.T) {
		t.Setenv(KeyEnvVar, identity.String())

		provider, err := NewProvider()
		require.NoError(t, err)
		assert.NotNil(t, provider)
		assert.Equal(t, identity.String(), provider.privateKey)
		assert.Equal(t, identity.Recipient().String(), provider.publicKey)
	})
}

func TestProvider_EncryptArgs(t *testing.T) {
	t.Parallel()

	identity, err := age.GenerateX25519Identity()
	require.NoError(t, err)

	provider := &Provider{
		privateKey: identity.String(),
		publicKey:  identity.Recipient().String(),
	}

	got, err := provider.EncryptArgs()
	require.NoError(t, err)

	t.Log("Returns correct arguments")
	{
		assert.Equal(t, []string{"--age", identity.Recipient().String()}, got)
	}

	t.Log("Is deterministic on multiple calls")
	{
		got2, err2 := provider.EncryptArgs()
		require.NoError(t, err2)
		assert.Equal(t, got, got2)
	}
}

func TestProvider_DecryptArgs(t *testing.T) {
	t.Parallel()

	provider := &Provider{}

	got, err := provider.DecryptArgs()
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestProvider_Environment(t *testing.T) {
	t.Parallel()

	identity, err := age.GenerateX25519Identity()
	require.NoError(t, err)

	provider := &Provider{
		privateKey: identity.String(),
		publicKey:  identity.Recipient().String(),
	}

	env := provider.Environment()

	require.NotNil(t, env)
	assert.Len(t, env, 1)
	assert.Contains(t, env, "SOPS_AGE_KEY")

	val := env["SOPS_AGE_KEY"]
	assert.Equal(t, identity.String(), val, "Should return the private key")
	assert.NotContains(t, val, identity.Recipient().String(), "Should not include public key")
	assert.Equal(t, env, provider.Environment(), "Should return same result on multiple calls")
}
