package infra

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

// Mock environment for testing
var testEnv = TFEnvironment{
	BackBlazeKeyID:              "bb-key",
	BackBlazeApplicationKey:     "bb-secret",
	DigitalOceanAPIKey:          "do-token",
	DigitalOceanSpacesAccessKey: "spaces-id",
	DigitalOceanSpacesSecretKey: "spaces-secret",
}

func TestWriteBackendConfig(t *testing.T) {
	t.Parallel()

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewReadOnlyFs(afero.NewMemMapFs())
		infraDir := "/infra"

		err := writeBackendConfig(fs, infraDir, testEnv)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "writing backend config")
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewMemMapFs()
		infraDir := "/infra"

		err := writeBackendConfig(fs, infraDir, testEnv)
		assert.NoError(t, err)

		path := filepath.Join(infraDir, backendFileName)

		t.Log("File Exists")
		{
			exists, err := afero.Exists(fs, path)
			assert.NoError(t, err)
			assert.True(t, exists, "backend file should exist")
		}

		t.Log("File Contents")
		{
			content, err := afero.ReadFile(fs, path)
			assert.NoError(t, err)
			expected := `access_key = "bb-key"
secret_key = "bb-secret"`
			assert.Equal(t, expected, string(content))
		}

		t.Log("Permissions")
		{
			info, err := fs.Stat(path)
			assert.NoError(t, err)
			assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
		}
	})
}

func TestWriteProviderConfig(t *testing.T) {
	t.Parallel()

	t.Run("Failure", func(t *testing.T) {
		t.Parallel()

		fs := afero.NewReadOnlyFs(afero.NewMemMapFs())
		infraDir := "/infra"

		err := writeProviderConfig(fs, infraDir, testEnv)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "writing provider config")
	})

	t.Run("Success", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		infraDir := "/infra"

		err := fs.MkdirAll(infraDir, os.ModePerm)
		assert.NoError(t, err)

		err = writeProviderConfig(fs, infraDir, testEnv)
		assert.NoError(t, err)

		path := filepath.Join(infraDir, providerFileName)

		t.Log("File Exists")
		{
			exists, err := afero.Exists(fs, path)
			assert.NoError(t, err)
			assert.True(t, exists, "provider file should exist")
		}

		t.Log("File Contents")
		{
			content, err := afero.ReadFile(fs, path)
			assert.NoError(t, err)
			expected := `do_token = "do-token"
spaces_access_id  = "spaces-id"
spaces_secret_key = "spaces-secret"`
			assert.Equal(t, expected, string(content))
		}

		t.Log("Permissions")
		{
			info, err := fs.Stat(path)
			assert.NoError(t, err)
			assert.Equal(t, os.FileMode(0600), info.Mode().Perm())
		}
	})
}
