package secrets

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/secrets"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestSync(t *testing.T) {
	envFixture := appdef.Environment{
		Production: map[string]appdef.EnvValue{
			"SECRET_KEY": {
				Source: appdef.EnvSourceSOPS,
				Value:  "production",
			},
		},
	}

	t.Run("No Files", func(t *testing.T) {
		t.Parallel()

		input, buf := setup(t, &appdef.Definition{})

		got := Sync(t.Context(), input)
		assert.NoError(t, got)
		assert.Contains(t, buf.String(), "No secrets")
	})

	t.Run("Scaffold Error", func(t *testing.T) {
		t.Parallel()

		def := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "app1", Env: envFixture},
			},
		}

		input, buf := setup(t, def)

		got := Sync(t.Context(), input)
		assert.NoError(t, got, "No production.yaml file causes error")
		assert.Contains(t, buf.String(), "Missing file")
	})

	t.Run("Invalid YAML", func(t *testing.T) {
		t.Parallel()

		def := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "app1", Env: envFixture},
			},
		}

		input, buf := setup(t, def)

		path := secrets.FilePathFromEnv(env.Production)
		err := afero.WriteFile(input.FS, path, []byte(`wrong\Yaml`), os.ModePerm)
		require.NoError(t, err)

		got := Sync(t.Context(), input)
		assert.Error(t, got)
		assert.Contains(t, buf.String(), "invalid YAML")
	})

	t.Run("Already Encrypted", func(t *testing.T) {
		def := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "app1", Env: envFixture},
			},
		}

		input, buf := setupEncryptedProdFile(t, `KEY: "1234"`)
		input.AppDefCache = def

		// This will (hopefully) encrypt the files so it's skipped.
		err := Encrypt(t.Context(), input)
		require.NoError(t, err)

		got := Sync(t.Context(), input)
		assert.NoError(t, got)
		assert.Contains(t, buf.String(), "Encrypted (skipped)")
	})

	t.Run("Key Already Exists", func(t *testing.T) {
		def := &appdef.Definition{
			Apps: []appdef.App{
				{
					Name: "app1",
					Env: appdef.Environment{
						Production: map[string]appdef.EnvValue{
							"SECRET_KEY": {
								Source: appdef.EnvSourceSOPS,
								Value:  "production",
							},
						},
					},
				},
			},
		}

		input, buf := setup(t, def)
		path := secrets.FilePathFromEnv(env.Production)

		// Write a file that already contains the secret
		initialContent := `SECRET_KEY: "EXISTING_VALUE"`
		err := afero.WriteFile(input.FS, path, []byte(initialContent), 0o644)
		require.NoError(t, err)

		got := Sync(t.Context(), input)
		assert.NoError(t, got)
		out := buf.String()

		// Should indicate that the secret was skipped
		assert.Contains(t, out, "• 0 added, 1 skipped")

		// File content should remain unchanged
		file, err := afero.ReadFile(input.FS, path)
		require.NoError(t, err)
		assert.Contains(t, string(file), `SECRET_KEY: "EXISTING_VALUE"`)
	})

	t.Run("Sync's Shared", func(t *testing.T) {
		t.Parallel()

		def := &appdef.Definition{
			Shared: appdef.Shared{
				Env: envFixture,
			},
		}

		input, _ := setup(t, def)

		err := Scaffold(t.Context(), input)
		assert.NoError(t, err)

		got := Sync(t.Context(), input)
		assert.NoError(t, got)

		file, err := afero.ReadFile(input.FS, secrets.FilePathFromEnv(env.Production))
		assert.NoError(t, err)
		assert.Contains(t, string(file), `SECRET_KEY: "REPLACE_ME_SECRET_KEY"`)
	})

	t.Run("Sync's Apps", func(t *testing.T) {
		t.Parallel()

		def := &appdef.Definition{
			Apps: []appdef.App{
				{Name: "app1", Env: envFixture},
				{Env: envFixture},
			},
		}

		input, _ := setup(t, def)

		err := Scaffold(t.Context(), input)
		assert.NoError(t, err)

		got := Sync(t.Context(), input)
		assert.NoError(t, got)

		file, err := afero.ReadFile(input.FS, secrets.FilePathFromEnv(env.Production))
		assert.NoError(t, err)
		assert.Contains(t, string(file), `SECRET_KEY: "REPLACE_ME_SECRET_KEY"`)
	})
}
