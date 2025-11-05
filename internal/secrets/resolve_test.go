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

	t.Run("Resource Invalid Format", func(t *testing.T) {
		def := &appdef.Definition{
			Resources: []appdef.Resource{
				{Name: "db", Type: appdef.ResourceTypePostgres},
			},
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Production: map[string]appdef.EnvValue{
						"DATABASE_URI": {Source: appdef.EnvSourceResource, Value: "invalid-format"},
					},
				},
			},
		}

		err := Resolve(t.Context(), def, ResolveConfig{})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid resource reference format")
		assert.ErrorContains(t, err, "DATABASE_URI")
		assert.ErrorContains(t, err, "expected 'resource_name.output_name'")
	})

	t.Run("Resource Not Found In Definition", func(t *testing.T) {
		def := &appdef.Definition{
			Resources: []appdef.Resource{
				{Name: "db", Type: appdef.ResourceTypePostgres},
			},
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Production: map[string]appdef.EnvValue{
						"DATABASE_URI": {Source: appdef.EnvSourceResource, Value: "nonexistent.connection_url"},
					},
				},
			},
		}

		err := Resolve(t.Context(), def, ResolveConfig{})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "resource 'nonexistent' not found in definition")
		assert.ErrorContains(t, err, "DATABASE_URI")
	})

	t.Run("Resource Env Var Not Set", func(t *testing.T) {
		def := &appdef.Definition{
			Resources: []appdef.Resource{
				{Name: "db", Type: appdef.ResourceTypePostgres},
			},
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Production: map[string]appdef.EnvValue{
						"DATABASE_URI": {Source: appdef.EnvSourceResource, Value: "db.connection_url"},
					},
				},
			},
		}

		// Make sure the env var is not set
		t.Setenv("TF_PROD_DB_CONNECTION_URL", "")
		os.Unsetenv("TF_PROD_DB_CONNECTION_URL")

		err := Resolve(t.Context(), def, ResolveConfig{})
		assert.Error(t, err)
		assert.ErrorContains(t, err, "environment variable 'TF_PROD_DB_CONNECTION_URL' not set")
		assert.ErrorContains(t, err, "DATABASE_URI")
	})

	t.Run("Resource Success Production", func(t *testing.T) {
		t.Setenv("TF_PROD_DB_CONNECTION_URL", "postgresql://user:pass@host:5432/dbname")

		def := &appdef.Definition{
			Resources: []appdef.Resource{
				{Name: "db", Type: appdef.ResourceTypePostgres},
			},
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Production: map[string]appdef.EnvValue{
						"DATABASE_URI": {Source: appdef.EnvSourceResource, Value: "db.connection_url"},
					},
				},
			},
		}

		err := Resolve(t.Context(), def, ResolveConfig{})
		require.NoError(t, err)
		assert.Equal(t, "postgresql://user:pass@host:5432/dbname", def.Shared.Env.Production["DATABASE_URI"].Value)
	})

	t.Run("Resource Success Multiple Outputs", func(t *testing.T) {
		t.Setenv("TF_DEV_DB_CONNECTION_URL", "postgresql://localhost:5432/dev")
		t.Setenv("TF_DEV_DB_HOST", "localhost")
		t.Setenv("TF_DEV_DB_PORT", "5432")

		def := &appdef.Definition{
			Resources: []appdef.Resource{
				{Name: "db", Type: appdef.ResourceTypePostgres},
			},
			Apps: []appdef.App{
				{
					Name: "test-app",
					Env: appdef.Environment{
						Dev: map[string]appdef.EnvValue{
							"DATABASE_URL": {Source: appdef.EnvSourceResource, Value: "db.connection_url"},
							"DATABASE_HOST": {Source: appdef.EnvSourceResource, Value: "db.host"},
							"DATABASE_PORT": {Source: appdef.EnvSourceResource, Value: "db.port"},
						},
					},
				},
			},
		}

		err := Resolve(t.Context(), def, ResolveConfig{})
		require.NoError(t, err)
		assert.Equal(t, "postgresql://localhost:5432/dev", def.Apps[0].Env.Dev["DATABASE_URL"].Value)
		assert.Equal(t, "localhost", def.Apps[0].Env.Dev["DATABASE_HOST"].Value)
		assert.Equal(t, "5432", def.Apps[0].Env.Dev["DATABASE_PORT"].Value)
	})

	t.Run("Resource Success With Hyphens", func(t *testing.T) {
		t.Setenv("TF_STAGING_MY_APP_DB_CONNECTION_URL", "postgresql://staging:5432/app")

		def := &appdef.Definition{
			Resources: []appdef.Resource{
				{Name: "my-app-db", Type: appdef.ResourceTypePostgres},
			},
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Staging: map[string]appdef.EnvValue{
						"DATABASE_URI": {Source: appdef.EnvSourceResource, Value: "my-app-db.connection_url"},
					},
				},
			},
		}

		err := Resolve(t.Context(), def, ResolveConfig{})
		require.NoError(t, err)
		assert.Equal(t, "postgresql://staging:5432/app", def.Shared.Env.Staging["DATABASE_URI"].Value)
	})

	t.Run("Resource Mixed With Other Sources", func(t *testing.T) {
		content := `
API_KEY: supersecret
`
		tmpDir, _ := writeTempSecret(t, content)

		t.Setenv("TF_PROD_DB_CONNECTION_URL", "postgresql://prod:5432/db")

		def := &appdef.Definition{
			Resources: []appdef.Resource{
				{Name: "db", Type: appdef.ResourceTypePostgres},
			},
			Apps: []appdef.App{
				{
					Name: "test-app",
					Env: appdef.Environment{
						Production: map[string]appdef.EnvValue{
							"DATABASE_URL": {Source: appdef.EnvSourceResource, Value: "db.connection_url"},
							"API_KEY":      {Source: appdef.EnvSourceSOPS, Value: "API_KEY"},
							"FRONTEND_URL": {Source: appdef.EnvSourceValue, Value: "https://example.com"},
						},
					},
				},
			},
		}

		ctrl := gomock.NewController(t)
		mockClient := mocks.NewMockEncrypterDecrypter(ctrl)
		mockClient.EXPECT().Decrypt(filepath.Join(tmpDir, FilePathFromEnv(env.Production))).Return(nil)
		mockClient.EXPECT().Encrypt(filepath.Join(tmpDir, FilePathFromEnv(env.Production))).Return(nil)

		cfg := ResolveConfig{
			SOPSClient: mockClient,
			BaseDir:    tmpDir,
		}

		err := Resolve(t.Context(), def, cfg)
		require.NoError(t, err)
		assert.Equal(t, "postgresql://prod:5432/db", def.Apps[0].Env.Production["DATABASE_URL"].Value)
		assert.Equal(t, "supersecret", def.Apps[0].Env.Production["API_KEY"].Value)
		assert.Equal(t, "https://example.com", def.Apps[0].Env.Production["FRONTEND_URL"].Value)
	})
}
