package env

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
