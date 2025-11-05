package env

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/manifest"
	"github.com/ainsleydev/webkit/internal/secrets/age"
	"github.com/ainsleydev/webkit/pkg/env"
)

func setup(t *testing.T, def *appdef.Definition) cmdtools.CommandInput {
	t.Helper()

	fs := afero.NewMemMapFs()
	input := cmdtools.CommandInput{
		FS:          fs,
		AppDefCache: def,
		Manifest:    manifest.NewTracker(),
	}

	return input
}

func TestWriteMapToFileCustomPath(t *testing.T) {
	ageIdentity, err := age.NewIdentity()
	require.NoError(t, err)
	t.Setenv(age.KeyEnvVar, ageIdentity.String())

	appDef := &appdef.Definition{
		Apps: []appdef.App{
			{
				Name: "test-app",
				Path: "./test-app",
				Env: appdef.Environment{
					Production: appdef.EnvVar{
						"FOO": {Value: "bar", Source: appdef.EnvSourceValue},
					},
				},
			},
		},
	}

	tt := map[string]struct {
		customPath string
		wantPath   string
	}{
		"Default Path": {
			customPath: "",
			wantPath:   "test-app/.env.production",
		},
		"Custom Path": {
			customPath: "/opt/myapp/.env",
			wantPath:   "/opt/myapp/.env",
		},
		"Custom Nested Path": {
			customPath: "/var/app/production/.env",
			wantPath:   "/var/app/production/.env",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			input := setup(t, appDef)

			err := writeMapToFile(writeArgs{
				Input:            input,
				Vars:             appDef.Apps[0].Env.Production,
				App:              appDef.Apps[0],
				Environment:      env.Production,
				CustomOutputPath: test.customPath,
			})
			require.NoError(t, err)

			exists, err := afero.Exists(input.FS, test.wantPath)
			require.NoError(t, err)
			assert.True(t, exists, fmt.Sprintf("Expected file at %s", test.wantPath))

			content, err := afero.ReadFile(input.FS, test.wantPath)
			require.NoError(t, err)
			assert.Contains(t, string(content), "FOO=\"bar\"")
		})
	}
}
