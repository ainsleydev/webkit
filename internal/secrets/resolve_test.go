package secrets

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/mocks"
	"github.com/ainsleydev/webkit/pkg/env"
)

func writeTempSecret(t *testing.T, content string) (tmpDir, secretPath string) {
	t.Helper()

	tmpDir = t.TempDir()
	secretPath = filepath.Join(tmpDir, FilePathFromEnv(env.Development))

	require.NoError(t, os.MkdirAll(filepath.Dir(secretPath), 0700))
	require.NoError(t, os.WriteFile(secretPath, []byte(content), 0600))

	return tmpDir, secretPath
}

func TestResolve(t *testing.T) {
	t.Run("Unknown Source", func(t *testing.T) {
		def := &appdef.Definition{
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Dev: map[string]appdef.EnvValue{
						"FOO": {Source: "unknown"},
					},
				},
			},
		}

		err := Resolve(t.Context(), def, ResolveConfig{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown env source type")
	})

	t.Run("App Resolve Error", func(t *testing.T) {
		def := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "failing-app",
					Env: appdef.Environment{
						Dev: map[string]appdef.EnvValue{
							"FOO": {Source: "unknown"}, // Triggers unknown source error
						},
					},
				},
			},
		}

		err := Resolve(t.Context(), def, ResolveConfig{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "resolving app \"failing-app\" env")
		assert.Contains(t, err.Error(), "unknown env source type")
	})

	t.Run("Decrypt Fails", func(t *testing.T) {
		tmpDir := t.TempDir()
		secretPath := filepath.Join(tmpDir, FilePathFromEnv(env.Development))

		def := &appdef.Definition{
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Dev: map[string]appdef.EnvValue{
						"API_KEY": {Source: appdef.EnvSourceSOPS},
					},
				},
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := mocks.NewMockEncrypterDecrypter(ctrl)
		mockClient.EXPECT().
			Decrypt(secretPath).
			Return(assert.AnError)

		cfg := ResolveConfig{
			SOPSClient: mockClient,
			BaseDir:    tmpDir,
		}

		err := Resolve(t.Context(), def, cfg)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "resolving shared env")
	})

	t.Run("Secret Not Found", func(t *testing.T) {
		tmpDir, secretPath := writeTempSecret(t, "OTHER_KEY: value")

		def := &appdef.Definition{
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Dev: map[string]appdef.EnvValue{
						"API_KEY": {Source: appdef.EnvSourceSOPS},
					},
				},
			},
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockClient := mocks.NewMockEncrypterDecrypter(ctrl)
		mockClient.EXPECT().Decrypt(secretPath).Return(nil)
		mockClient.EXPECT().Encrypt(secretPath).Return(nil)

		cfg := ResolveConfig{
			SOPSClient: mockClient,
			BaseDir:    tmpDir,
		}

		err := Resolve(t.Context(), def, cfg)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "secret 'API_KEY' not found")
	})

	t.Run("Success", func(t *testing.T) {
		content := `
API_KEY: supersecret
DB_PASS: dbpass123
`

		tmpDir, secretPath := writeTempSecret(t, content)

		def := &appdef.Definition{
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Dev: map[string]appdef.EnvValue{
						"API_KEY": {Source: appdef.EnvSourceSOPS},
					},
				},
			},
			Apps: []appdef.App{
				{
					Name: "test-app",
					Env: appdef.Environment{
						Dev: map[string]appdef.EnvValue{
							"DB_PASS": {Source: appdef.EnvSourceSOPS},
						},
					},
				},
			},
		}

		ctrl := gomock.NewController(t)
		mockClient := mocks.NewMockEncrypterDecrypter(ctrl)

		mockClient.EXPECT().
			Decrypt(secretPath).
			Return(nil).
			Times(2)
		mockClient.EXPECT().
			Encrypt(secretPath).
			Return(nil).
			Times(2)

		cfg := ResolveConfig{
			SOPSClient: mockClient,
			BaseDir:    tmpDir,
		}

		err := Resolve(t.Context(), def, cfg)
		require.NoError(t, err)
		assert.Equal(t, def.Shared.Env.Dev["API_KEY"].Value, "supersecret")
		assert.Equal(t, def.Apps[0].Env.Dev["DB_PASS"].Value, "dbpass123")
	})

}
