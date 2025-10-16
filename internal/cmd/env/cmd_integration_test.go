//go:build integration
// +build integration

package env

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmd/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets/age"
	"github.com/ainsleydev/webkit/internal/util/executil"
)

func TestEnvIntegration(t *testing.T) {
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

	appDef := &appdef.Definition{
		Shared: appdef.Shared{
			Env: appdef.Environment{
				Production: map[string]appdef.EnvValue{
					"SHARED_KEY": {Value: "shared_value", Source: appdef.EnvSourceSOPS},
				},
			},
		},
		Apps: []appdef.App{
			{
				Name: "api",
				Path: "apps/api",
				Env: appdef.Environment{
					Dev: map[string]appdef.EnvValue{
						"API_DEV_KEY": {Value: "dev_value", Source: appdef.EnvSourceValue},
					},
					Production: map[string]appdef.EnvValue{
						"API_PROD_KEY": {Value: "prod_value", Source: appdef.EnvSourceValue},
					},
				},
			},
			{
				Name: "web",
				Path: "apps/web",
				Env: appdef.Environment{
					Dev: map[string]appdef.EnvValue{
						"WEB_DEV_KEY": {Value: "dev_value", Source: appdef.EnvSourceValue},
					},
					Production: map[string]appdef.EnvValue{
						"WEB_PROD_KEY": {Value: "prod_value", Source: appdef.EnvSourceSOPS},
					},
				},
			},
		},
	}

	input := cmdtools.CommandInput{
		FS:          fs,
		BaseDir:     tmpDir,
		AppDefCache: appDef,
	}

	err = Scaffold(ctx, input)
	require.NoError(t, err)

	for _, app := range appDef.Apps {
		for _, envName := range []string{"", "production"} {
			filename := ".env"
			if envName != "" {
				filename = ".env." + envName
			}
			path := filepath.Join(app.Path, filename)

			exists, err := afero.Exists(fs, path)
			require.NoError(t, err)
			assert.Truef(t, exists, "%s should exist after scaffold", path)
		}
	}

	err = Sync(ctx, input)
	require.NoError(t, err)

	for _, app := range appDef.Apps {
		devPath := filepath.Join(app.Path, ".env")
		prodPath := filepath.Join(app.Path, ".env.production")

		devFile, err := afero.ReadFile(fs, devPath)
		require.NoError(t, err)

		prodFile, err := afero.ReadFile(fs, prodPath)
		require.NoError(t, err)

		assert.Contains(t, string(devFile), "API_DEV_KEY")
		assert.Contains(t, string(prodFile), "API_PROD_KEY")
		assert.Contains(t, string(prodFile), "SHARED_KEY")
	}
}
