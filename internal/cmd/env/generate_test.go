package env

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/secrets/age"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestGetEnvironmentVars(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   appdef.Environment
		env     env.Environment
		want    appdef.EnvVar
		wantErr bool
	}{
		"Development": {
			input: appdef.Environment{
				Dev: appdef.EnvVar{"FOO": {Value: "bar"}},
			},
			env: env.Development,
			want: appdef.EnvVar{
				"FOO": {Value: "bar"},
			},
		},
		"Staging": {
			input: appdef.Environment{
				Staging: appdef.EnvVar{"BAZ": {Value: "qux"}},
			},
			env: env.Staging,
			want: appdef.EnvVar{
				"BAZ": {Value: "qux"},
			},
		},
		"Production": {
			input: appdef.Environment{
				Production: appdef.EnvVar{"PROD": {Value: "value"}},
			},
			env: env.Production,
			want: appdef.EnvVar{
				"PROD": {Value: "value"},
			},
		},
		"Unsupported": {
			input:   appdef.Environment{},
			env:     env.Environment("unknown"),
			wantErr: true,
		},
		"Empty Dev": {
			input: appdef.Environment{
				Dev: appdef.EnvVar{},
			},
			env:  env.Development,
			want: appdef.EnvVar{},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := test.input.GetVarsForEnvironment(test.env)
			assert.Equal(t, test.wantErr, err != nil)

			if !test.wantErr {
				assert.Equal(t, test.want, got)
			}
		})
	}
}

func TestEnvSuffix(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input env.Environment
		want  string
	}{
		"Development": {input: env.Development, want: ""},
		"Staging":     {input: env.Staging, want: ".staging"},
		"Production":  {input: env.Production, want: ".production"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := envSuffix(test.input)
			assert.Equal(t, test.want, got)
		})
	}
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
