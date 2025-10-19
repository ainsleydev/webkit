package age

import (
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/config"
)

func TestReadIdentity(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	require.NoError(t, err)

	t.Run("From Environment Variable", func(t *testing.T) {
		t.Setenv(KeyEnvVar, identity.String())

		got, err := ReadIdentity()
		require.NoError(t, err)
		assert.Equal(t, identity.String(), got.String())
	})

	t.Run("From Config File", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)
		require.NoError(t, os.Unsetenv(KeyEnvVar))

		// Write key to config
		err := WritePrivateKey(identity.String())
		require.NoError(t, err)

		// Read it back
		got, err := ReadIdentity()
		require.NoError(t, err)
		assert.Equal(t, identity.String(), got.String())
	})

	t.Run("Environment Takes Precedence Over File", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		// Write one key to file
		fileIdentity, err := age.GenerateX25519Identity()
		require.NoError(t, err)
		err = WritePrivateKey(fileIdentity.String())
		require.NoError(t, err)

		// Set different key in environment
		envIdentity, err := age.GenerateX25519Identity()
		require.NoError(t, err)
		t.Setenv(KeyEnvVar, envIdentity.String())

		// Should get the env key
		got, err := ReadIdentity()
		require.NoError(t, err)
		assert.Equal(t, envIdentity.String(), got.String())
		assert.NotEqual(t, fileIdentity.String(), got.String())
	})

	t.Run("File Not Found", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)
		require.NoError(t, os.Unsetenv(KeyEnvVar))

		_, err = ReadIdentity()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reading age key from")
		assert.Contains(t, err.Error(), "age.key")
	})

	t.Run("Invalid Key From Environment", func(t *testing.T) {
		t.Setenv(KeyEnvVar, "not-a-valid-age-key")

		_, err = ReadIdentity()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid age key format")
		assert.Contains(t, err.Error(), "SOPS_AGE_KEY environment variable")
	})

	t.Run("Invalid Key From File", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)
		require.NoError(t, os.Unsetenv(KeyEnvVar))

		// Write invalid key directly to bypass validation
		err = config.Write(KeyFileName, []byte("not-a-valid-age-key"), 0o600)
		require.NoError(t, err)

		_, err = ReadIdentity()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid age key format")
		assert.Contains(t, err.Error(), "age.key")
	})

	t.Run("Empty Key From Environment", func(t *testing.T) {
		t.Setenv(KeyEnvVar, "")
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		// Empty env var should fall through to file
		_, err = ReadIdentity()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reading age key")
	})

	t.Run("Empty File From Environment", func(t *testing.T) {
		t.Setenv(KeyEnvVar, "")
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		err = config.Write(KeyFileName, []byte(""), os.ModePerm)
		require.NoError(t, err)

		_, err = ReadIdentity()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no SOPS_AGE_KEY key found")
	})
}

func TestWritePrivateKey(t *testing.T) {
	identity, err := age.GenerateX25519Identity()
	require.NoError(t, err)

	t.Run("Invalid Key Format", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		err = WritePrivateKey("not-a-valid-age-key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid age key format")

		// File should not have been created
		_, err = config.Read(KeyFileName)
		assert.Error(t, err)
	})

	t.Run("Empty Key", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		err = WritePrivateKey("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid age key format")
	})

	t.Run("Write Error", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		dir, err := config.Dir()
		require.NoError(t, err)
		require.NoError(t, os.MkdirAll(dir, 0o755))

		// Create a directory where the file should go
		filePath := filepath.Join(dir, KeyFileName)
		require.NoError(t, os.Mkdir(filePath, 0o755)) // <-- Dir instead of a file

		err = WritePrivateKey(identity.String())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "writing age key")
	})

	t.Run("Success", func(t *testing.T) {
		tmpHome := t.TempDir()
		t.Setenv("HOME", tmpHome)

		err = WritePrivateKey(identity.String())
		require.NoError(t, err)

		// Verify it was written correctly
		data, err := config.Read(KeyFileName)
		require.NoError(t, err)
		assert.Equal(t, identity.String(), string(data))
	})
}

func TestAgeKey_ReadWriteRoundTrip(t *testing.T) {
	identity, err := NewIdentity()
	require.NoError(t, err)

	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	require.NoError(t, os.Unsetenv(KeyEnvVar))

	err = WritePrivateKey(identity.String())
	require.NoError(t, err)

	key, err := ReadIdentity()
	require.NoError(t, err)

	assert.Equal(t, identity.String(), key.String())
}

func TestExtractPublicKey(t *testing.T) {
	t.Parallel()

	validIdentity, _ := age.GenerateX25519Identity()
	validPrivateKey := validIdentity.String()
	validPublicKey := validIdentity.Recipient().String()

	tt := map[string]struct {
		input   string
		want    string
		wantErr bool
	}{
		"Valid Age Private Key": {
			input:   validPrivateKey,
			want:    validPublicKey,
			wantErr: false,
		},
		"Empty String": {
			input:   "",
			want:    "",
			wantErr: true,
		},
		"Random String": {
			input:   "not-a-valid-key",
			want:    "",
			wantErr: true,
		},
		"Whitespace Only": {
			input:   "   \n\t",
			want:    "",
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := extractPublicKey(test.input)
			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}
