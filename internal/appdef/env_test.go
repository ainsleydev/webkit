package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ainsleydev/webkit/pkg/env"
)

func TestEnvSource_String(t *testing.T) {
	t.Parallel()

	got := EnvSourceValue.String()
	assert.Equal(t, "value", got)
	assert.IsType(t, "", got)
}

func TestEnvironment_Walk(t *testing.T) {
	t.Parallel()

	e := Environment{
		Dev:        EnvVar{"DEBUG": {Value: "true"}},
		Production: EnvVar{"DEBUG": {Value: "false"}},
	}

	var got []string
	e.Walk(func(envName string, envVars EnvVar) {
		got = append(got, envName)
	})

	want := []string{env.Development, env.Production}
	assert.ElementsMatch(t, want, got)
}

func TestEnvValue_ParseSOPSPath(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   EnvValue
		want    SOPSPath
		wantErr bool
	}{
		"Valid SOPS path": {
			input: EnvValue{
				Source: EnvSourceSOPS,
				Path:   "secrets/production.yaml:PAYLOAD_SECRET",
			},
			want: SOPSPath{
				File: "secrets/production.yaml",
				Key:  "PAYLOAD_SECRET",
			},
			wantErr: false,
		},
		"Valid SOPS path with nested key": {
			input: EnvValue{
				Source: EnvSourceSOPS,
				Path:   "secrets/staging.yaml:DATABASE_PASSWORD",
			},
			want: SOPSPath{
				File: "secrets/staging.yaml",
				Key:  "DATABASE_PASSWORD",
			},
			wantErr: false,
		},
		"Not a SOPS source": {
			input: EnvValue{
				Source: EnvSourceValue,
				Path:   "secrets/production.yaml:PAYLOAD_SECRET",
			},
			want:    SOPSPath{},
			wantErr: true,
		},
		"Resource source": {
			input: EnvValue{
				Source: EnvSourceResource,
				Path:   "secrets/production.yaml:PAYLOAD_SECRET",
			},
			want:    SOPSPath{},
			wantErr: true,
		},
		"Invalid format no colon": {
			input: EnvValue{
				Source: EnvSourceSOPS,
				Path:   "secrets/production.yaml",
			},
			want:    SOPSPath{},
			wantErr: true,
		},
		"Invalid format multiple colons": {
			input: EnvValue{
				Source: EnvSourceSOPS,
				Path:   "secrets/production.yaml:KEY:EXTRA",
			},
			want:    SOPSPath{},
			wantErr: true,
		},
		"Empty path": {
			input: EnvValue{
				Source: EnvSourceSOPS,
				Path:   "",
			},
			want:    SOPSPath{},
			wantErr: true,
		},
		"Path with subdirectories": {
			input: EnvValue{
				Source: EnvSourceSOPS,
				Path:   "config/secrets/prod/app.yaml:API_KEY",
			},
			want: SOPSPath{
				File: "config/secrets/prod/app.yaml",
				Key:  "API_KEY",
			},
			wantErr: false,
		},
		"Empty file path": {
			input: EnvValue{
				Source: EnvSourceSOPS,
				Path:   ":KEY",
			},
			want: SOPSPath{
				File: "",
				Key:  "KEY",
			},
			wantErr: false,
		},
		"Empty key": {
			input: EnvValue{
				Source: EnvSourceSOPS,
				Path:   "secrets/production.yaml:",
			},
			want: SOPSPath{
				File: "secrets/production.yaml",
				Key:  "",
			},
			wantErr: false,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := test.input.ParseSOPSPath()
			assert.Equal(t, test.want, got)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}
