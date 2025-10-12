package secrets

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/executil"
	"github.com/ainsleydev/webkit/internal/secrets/age"
)

func TestSecretsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	if !executil.Exists("sops") {
		t.Skip("sops CLI not found in PATH; skipping integration test")
	}

	ctx := t.Context()
	tmpDir := t.TempDir()
	fs := afero.NewBasePathFs(afero.NewOsFs(), tmpDir)

	ageIdentity, err := age.NewIdentity()
	require.NoError(t, err)
	t.Setenv(age.KeyEnvVar, ageIdentity.String())

	input := cmdtools.CommandInput{
		FS:      fs,
		BaseDir: tmpDir,
		AppDefCache: &appdef.Definition{
			Shared: appdef.Shared{
				Env: appdef.Environment{
					Production: map[string]appdef.EnvValue{
						"SECRET_SHARED_VAR": {Value: "secret", Source: appdef.EnvSourceSOPS},
					},
				},
			},
			Apps: []appdef.App{
				{
					Name: "sops",
					Env: appdef.Environment{
						Dev: map[string]appdef.EnvValue{
							"SECRET_APP_VAR": {Value: "secret", Source: appdef.EnvSourceSOPS},
						},
						Production: map[string]appdef.EnvValue{
							"SECRET_APP_VAR": {Value: "secret", Source: appdef.EnvSourceSOPS},
						},
						Staging: map[string]appdef.EnvValue{
							"SECRET_APP_VAR": {Value: "secret", Source: appdef.EnvSourceSOPS},
						},
					},
				},
			},
		},
	}

	t.Log("Creates Files")
	{
		err = CreateFiles(ctx, input)
		assert.NoError(t, err)
	}

	// Sync doesnt write values for you, it only writes placeholders
	// so we need to write back to the file.

	t.Log("Syncs Vars to Environments")
	{
		err = Sync(ctx, input)
		assert.NoError(t, err)
	}

	t.Log("Encrypts Files")
	{
		err = Encrypt(ctx, input)
		assert.NoError(t, err)
	}

	t.Log("Decrypts Files")
	{
		err = Decrypt(ctx, input)
		assert.NoError(t, err)
	}

	t.Log("Gets Vars")
	{
		assert.NoError(t, Encrypt(ctx, input))

		input.Command = &cli.Command{
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "env"},
				&cli.StringFlag{Name: "key"},
			},
		}
		require.NoError(t, input.Command.Set("env", "production"))
		require.NoError(t, input.Command.Set("key", "SECRET_APP_VAR"))

		err = Get(ctx, input)
		require.NoError(t, err)

		fmt.Println(err)
	}

}
