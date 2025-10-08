package sops

import (
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/config"
)

func TestAgeKeyRead(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	require.NoError(t, err)

	t.Run("From Environment Variable", func(t *testing.T) {
		t.Setenv(AgeKeyEnvVar, identity.String())

		key, err := AgeKeyRead()
		require.NoError(t, err)
		assert.Equal(t, identity.String(), key)
	})

	t.Run("From Config File", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)
		require.NoError(t, os.Unsetenv(AgeKeyEnvVar))

		// Write key to config
		err := AgeKeyWrite(identity.String())
		require.NoError(t, err)

		// Read it back
		key, err := AgeKeyRead()
		require.NoError(t, err)
		assert.Equal(t, identity.String(), key)
	})

	t.Run("Environment Takes Precedence Over File", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		// Write one key to file
		fileIdentity, err := age.GenerateX25519Identity()
		require.NoError(t, err)
		err = AgeKeyWrite(fileIdentity.String())
		require.NoError(t, err)

		// Set different key in environment
		envIdentity, err := age.GenerateX25519Identity()
		require.NoError(t, err)
		t.Setenv(AgeKeyEnvVar, envIdentity.String())

		// Should get the env key
		key, err := AgeKeyRead()
		require.NoError(t, err)
		assert.Equal(t, envIdentity.String(), key)
		assert.NotEqual(t, fileIdentity.String(), key)
	})

	t.Run("File Not Found", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)
		require.NoError(t, os.Unsetenv(AgeKeyEnvVar))

		_, err = AgeKeyRead()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reading age key from")
		assert.Contains(t, err.Error(), "age.key")
	})

	t.Run("Invalid Key From Environment", func(t *testing.T) {
		t.Setenv(AgeKeyEnvVar, "not-a-valid-age-key")

		_, err = AgeKeyRead()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid age key format")
		assert.Contains(t, err.Error(), "SOPS_AGE_KEY environment variable")
	})

	t.Run("Invalid Key From File", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)
		require.NoError(t, os.Unsetenv(AgeKeyEnvVar))

		// Write invalid key directly to bypass validation
		err = config.Write(AgeKeyFileName, []byte("not-a-valid-age-key"), 0600)
		require.NoError(t, err)

		_, err = AgeKeyRead()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid age key format")
		assert.Contains(t, err.Error(), "age.key")
	})

	t.Run("Empty Key From Environment", func(t *testing.T) {
		t.Setenv(AgeKeyEnvVar, "")
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		// Empty env var should fall through to file
		_, err = AgeKeyRead()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reading age key")
	})

	t.Run("Whitespace In Key From File", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)
		require.NoError(t, os.Unsetenv(AgeKeyEnvVar))

		// Write key with extra whitespace
		keyWithWhitespace := "\n" + identity.String() + "\n"
		err = config.Write(AgeKeyFileName, []byte(keyWithWhitespace), 0600)
		require.NoError(t, err)

		_, err = AgeKeyRead()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid age key format")
	})
}

func TestAgeKeyWrite(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	require.NoError(t, err)

	t.Run("Valid Key", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		err = AgeKeyWrite(identity.String())
		require.NoError(t, err)

		// Verify it was written correctly
		data, err := config.Read(AgeKeyFileName)
		require.NoError(t, err)
		assert.Equal(t, identity.String(), string(data))
	})

	t.Run("Creates Config Directory", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		// Config dir shouldn't exist yet
		dir, err := config.Dir()
		require.NoError(t, err)
		_, err = os.Stat(dir)
		assert.True(t, os.IsNotExist(err))

		// Write key should create directory
		err = AgeKeyWrite(identity.String())
		require.NoError(t, err)

		// Directory should exist now
		info, err := os.Stat(dir)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
	})

	t.Run("Overwrites Existing Key", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		// Write first key
		firstIdentity, err := age.GenerateX25519Identity()
		require.NoError(t, err)
		err = AgeKeyWrite(firstIdentity.String())
		require.NoError(t, err)

		// Write second key
		secondIdentity, err := age.GenerateX25519Identity()
		require.NoError(t, err)
		err = AgeKeyWrite(secondIdentity.String())
		require.NoError(t, err)

		// Should have second key
		data, err := config.Read(AgeKeyFileName)
		require.NoError(t, err)
		assert.Equal(t, secondIdentity.String(), string(data))
		assert.NotEqual(t, firstIdentity.String(), string(data))
	})

	t.Run("Invalid Key Format", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		err = AgeKeyWrite("not-a-valid-age-key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid age key format")

		// File should not have been created
		_, err = config.Read(AgeKeyFileName)
		assert.Error(t, err)
	})

	t.Run("Empty Key", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		err = AgeKeyWrite("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid age key format")
	})

	t.Run("Write Error", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		dir, err := config.Dir()
		require.NoError(t, err)
		require.NoError(t, os.MkdirAll(dir, 0755))

		// Create a directory where the file should go
		filePath := filepath.Join(dir, AgeKeyFileName)
		require.NoError(t, os.Mkdir(filePath, 0755)) // <-- Dir instead of a file

		err = AgeKeyWrite(identity.String())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "writing age key")
	})
}

func TestAgeKey_ReadWriteRoundTrip(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	require.NoError(t, err)

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	require.NoError(t, os.Unsetenv(AgeKeyEnvVar))

	err = AgeKeyWrite(identity.String())
	require.NoError(t, err)

	key, err := AgeKeyRead()
	require.NoError(t, err)

	assert.Equal(t, identity.String(), key)

	parsedIdentity, err := age.ParseX25519Identity(key)
	require.NoError(t, err)
	assert.Equal(t, identity.Recipient().String(), parsedIdentity.Recipient().String())
}
