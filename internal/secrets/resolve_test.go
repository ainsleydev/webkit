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

	require.NoError(t, os.MkdirAll(filepath.Dir(secretPath), 0o700))
	require.NoError(t, os.WriteFile(secretPath, []byte(content), 0o600))

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

func TestResolveForEnvironment(t *testing.T) {
	t.Run("Only Resolves Target Environment", func(t *testing.T) {
		def := &appdef.Definition{
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Dev: map[string]appdef.EnvValue{
						"DEV_VAR": {Source: appdef.EnvSourceValue, Value: "dev"},
					},
					Production: map[string]appdef.EnvValue{
						"PROD_VAR": {Source: appdef.EnvSourceValue, Value: "prod"},
					},
				},
			},
			Apps: []appdef.App{
				{
					Name: "test-app",
					Env: appdef.Environment{
						Dev: map[string]appdef.EnvValue{
							"APP_DEV": {Source: appdef.EnvSourceValue, Value: "app-dev"},
						},
						Production: map[string]appdef.EnvValue{
							"APP_PROD": {Source: appdef.EnvSourceValue, Value: "app-prod"},
						},
					},
				},
			},
		}

		// Resolve only production
		err := ResolveForEnvironment(t.Context(), def, env.Production, ResolveConfig{})
		require.NoError(t, err)

		// Production should be resolved
		assert.Equal(t, "prod", def.Shared.Env.Production["PROD_VAR"].Value)
		assert.Equal(t, "app-prod", def.Apps[0].Env.Production["APP_PROD"].Value)

		// Dev should still have original values
		assert.Equal(t, "dev", def.Shared.Env.Dev["DEV_VAR"].Value)
		assert.Equal(t, "app-dev", def.Apps[0].Env.Dev["APP_DEV"].Value)
	})

	t.Run("Resolves Resource References Only For Target Environment", func(t *testing.T) {
		def := &appdef.Definition{
			Shared: appdef.Shared{},
			Apps: []appdef.App{
				{
					Name: "test-app",
					Env: appdef.Environment{
						Dev: map[string]appdef.EnvValue{
							"DATABASE_URI": {Source: appdef.EnvSourceResource, Value: "db.connection_url"},
						},
						Production: map[string]appdef.EnvValue{
							"DATABASE_URI": {Source: appdef.EnvSourceResource, Value: "db.connection_url"},
						},
					},
				},
			},
		}

		// Only provide Terraform outputs for production
		tfOutputs := &TerraformOutputProvider{
			OutputKey{
				Environment:  env.Production,
				ResourceName: "db",
				OutputName:   "connection_url",
			}: "postgresql://prod-db:5432",
		}

		cfg := ResolveConfig{
			TerraformOutput: tfOutputs,
		}

		// Should succeed - only resolves production, doesn't touch dev
		err := ResolveForEnvironment(t.Context(), def, env.Production, cfg)
		require.NoError(t, err)

		// Production should have resolved value
		assert.Equal(t, "postgresql://prod-db:5432", def.Apps[0].Env.Production["DATABASE_URI"].Value)

		// Dev should still have original resource reference (not resolved)
		assert.Equal(t, "db.connection_url", def.Apps[0].Env.Dev["DATABASE_URI"].Value)
	})

	t.Run("Resolves Defaults For Target Environment", func(t *testing.T) {
		def := &appdef.Definition{
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Default: map[string]appdef.EnvValue{
						"SHARED_VAR": {Source: appdef.EnvSourceValue, Value: "shared"},
					},
					Production: map[string]appdef.EnvValue{
						"PROD_VAR": {Source: appdef.EnvSourceValue, Value: "prod"},
					},
				},
			},
		}

		err := ResolveForEnvironment(t.Context(), def, env.Production, ResolveConfig{})
		require.NoError(t, err)

		// Both default and production-specific should be resolved
		assert.Equal(t, "shared", def.Shared.Env.Default["SHARED_VAR"].Value)
		assert.Equal(t, "prod", def.Shared.Env.Production["PROD_VAR"].Value)
	})

	t.Run("Error In Target Environment", func(t *testing.T) {
		def := &appdef.Definition{
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Production: map[string]appdef.EnvValue{
						"BAD_VAR": {Source: "unknown"},
					},
				},
			},
		}

		err := ResolveForEnvironment(t.Context(), def, env.Production, ResolveConfig{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown env source type")
	})

	t.Run("Missing Terraform Outputs For Target Environment", func(t *testing.T) {
		def := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "test-app",
					Env: appdef.Environment{
						Production: map[string]appdef.EnvValue{
							"DATABASE_URI": {Source: appdef.EnvSourceResource, Value: "db.connection_url"},
						},
					},
				},
			},
		}

		// Terraform outputs only for dev, not production
		tfOutputs := &TerraformOutputProvider{
			OutputKey{
				Environment:  env.Development,
				ResourceName: "db",
				OutputName:   "connection_url",
			}: "postgresql://dev-db:5432",
		}

		cfg := ResolveConfig{
			TerraformOutput: tfOutputs,
		}

		err := ResolveForEnvironment(t.Context(), def, env.Production, cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "terraform output not found")
	})
}
