package secrets

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets/age"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestGet(t *testing.T) {
	ctx := t.Context()
	ageIdentity, err := age.NewIdentity()
	require.NoError(t, err)

	tmpDir := t.TempDir()

	t.Run("Client Error", func(t *testing.T) {
		input := cmdtools.CommandInput{Command: GetCmd}

		err = Encrypt(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "age")
	})

	t.Run("Decode Error", func(t *testing.T) {
		t.Setenv(age.KeyEnvVar, ageIdentity.String())
		fs := afero.NewBasePathFs(afero.NewOsFs(), tmpDir)

		input := cmdtools.CommandInput{
			FS:      fs,
			BaseDir: tmpDir,
			Command: GetCmd,
		}
		err = CreateFiles(ctx, input)
		assert.NoError(t, err)

		content := "KEY: VALUE\ninvalid-yaml"
		path := "resources/secrets/production.yaml"

		err = afero.WriteFile(fs, path, []byte(content), os.ModePerm)
		require.NoError(t, err)

		err = Get(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "decoding sops to map")
	})

	t.Run("No Value", func(t *testing.T) {
		t.Setenv(age.KeyEnvVar, ageIdentity.String())
		fs := afero.NewBasePathFs(afero.NewOsFs(), tmpDir)

		input := cmdtools.CommandInput{
			FS:      fs,
			BaseDir: tmpDir,
			Command: GetCmd,
		}
		err = CreateFiles(ctx, input)
		assert.NoError(t, err)

		content := `KEY: "1234"`
		path := "resources/secrets/production.yaml"

		err = afero.WriteFile(fs, path, []byte(content), os.ModePerm)
		require.NoError(t, err)

		err = Encrypt(ctx, input)
		require.NoError(t, err)

		require.NoError(t, input.Command.Set("env", env.Production))
		require.NoError(t, input.Command.Set("key", "wrong"))

		err = Get(ctx, input)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "key wrong not found")
	})

	t.Run("Success", func(t *testing.T) {
		t.Setenv(age.KeyEnvVar, ageIdentity.String())
		fs := afero.NewBasePathFs(afero.NewOsFs(), tmpDir)

		input := cmdtools.CommandInput{
			FS:      fs,
			BaseDir: tmpDir,
			Command: GetCmd,
		}
		err = CreateFiles(ctx, input)
		assert.NoError(t, err)

		content := `KEY: "1234"`
		path := "resources/secrets/production.yaml"

		err = afero.WriteFile(fs, path, []byte(content), os.ModePerm)
		require.NoError(t, err)

		err = Encrypt(ctx, input)
		require.NoError(t, err)

		require.NoError(t, input.Command.Set("env", env.Production))
		require.NoError(t, input.Command.Set("key", "KEY"))

		err = Get(ctx, input)
		assert.NoError(t, err)
		// TODO: Assert that output has value.
	})
}
