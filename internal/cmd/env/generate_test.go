package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
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

// TestGenerateDefaultSOPSResolution is a regression test for the stale-pointer bug
// where targetApp was captured as a range copy before ResolveForEnvironment ran,
// causing env.default SOPS vars to appear as empty strings in the generated file.
// It simulates the resolution step and confirms that reading from the updated slice
// element (not the pre-resolution copy) produces non-nil values.
func TestGenerateDefaultSOPSResolution(t *testing.T) {
	t.Parallel()

	t.Run("Stale copy has nil SOPS value", func(t *testing.T) {
		t.Parallel()

		apps := []appdef.App{
			{
				Name: "cms",
				Env: appdef.Environment{
					Default: appdef.EnvVar{
						"SECRET_KEY": {Source: appdef.EnvSourceSOPS, Value: nil},
					},
					Production: appdef.EnvVar{
						"NODE_ENV": {Source: appdef.EnvSourceValue, Value: "production"},
					},
				},
			},
		}

		// Simulate what ResolveForEnvironment writes: resolved value into the slice element.
		apps[0].Env.Production["SECRET_KEY"] = appdef.EnvValue{
			Source: appdef.EnvSourceSOPS,
			Value:  "super-secret",
		}

		// A range copy (as the old code did) captured before resolution
		// would still hold Value: nil for SECRET_KEY.
		var staleCopy *appdef.App
		origApps := []appdef.App{
			{
				Name: "cms",
				Env: appdef.Environment{
					Default: appdef.EnvVar{
						"SECRET_KEY": {Source: appdef.EnvSourceSOPS, Value: nil},
					},
					Production: appdef.EnvVar{
						"NODE_ENV": {Source: appdef.EnvSourceValue, Value: "production"},
					},
				},
			},
		}
		for _, app := range origApps {
			if app.Name == "cms" {
				staleCopy = &app //nolint:gosec // intentional copy to reproduce the bug
				break
			}
		}
		require.NotNil(t, staleCopy)

		staleVars, err := staleCopy.MergeEnvironments(appdef.Environment{}).GetVarsForEnvironment(env.Production)
		require.NoError(t, err)
		assert.Nil(t, staleVars["SECRET_KEY"].Value, "stale copy should show the bug: Value is nil")
	})

	t.Run("Slice element has resolved SOPS value", func(t *testing.T) {
		t.Parallel()

		apps := []appdef.App{
			{
				Name: "cms",
				Env: appdef.Environment{
					Default: appdef.EnvVar{
						"SECRET_KEY": {Source: appdef.EnvSourceSOPS, Value: nil},
					},
					Production: appdef.EnvVar{
						"NODE_ENV": {Source: appdef.EnvSourceValue, Value: "production"},
					},
				},
			},
		}

		// Simulate ResolveForEnvironment writing the resolved value into the Production map.
		apps[0].Env.Production["SECRET_KEY"] = appdef.EnvValue{
			Source: appdef.EnvSourceSOPS,
			Value:  "super-secret",
		}

		// Reading from the slice element (as the fixed code does) yields the resolved value.
		vars, err := apps[0].MergeEnvironments(appdef.Environment{}).GetVarsForEnvironment(env.Production)
		require.NoError(t, err)
		assert.Equal(t, "super-secret", vars["SECRET_KEY"].Value, "slice element should have resolved value")
		assert.Equal(t, "production", vars["NODE_ENV"].Value)
	})
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
