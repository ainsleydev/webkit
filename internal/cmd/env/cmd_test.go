package env

import (
	"fmt"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/internal/cmdtools"
	"github.com/ainsleydev/webkit/internal/secrets/age"
	"github.com/ainsleydev/webkit/internal/state/manifest"
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
			// Updated: our custom marshaller doesn't quote simple values
			assert.Contains(t, string(content), "FOO=bar")
		})
	}
}

func TestMarshalEnvWithoutQuotes(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input map[string]string
		want  map[string]string // key -> expected format in output
	}{
		"Simple values without quotes": {
			input: map[string]string{
				"DATABASE_URI": "postgresql://user:pass@host:5432/db?sslmode=require",
				"API_KEY":      "abc123xyz",
				"PORT":         "3000",
			},
			want: map[string]string{
				"DATABASE_URI": "DATABASE_URI=postgresql://user:pass@host:5432/db?sslmode=require",
				"API_KEY":      "API_KEY=abc123xyz",
				"PORT":         "PORT=3000",
			},
		},
		"Values with spaces need quotes": {
			input: map[string]string{
				"MESSAGE":     "hello world",
				"DESCRIPTION": "This is a description with spaces",
			},
			want: map[string]string{
				"MESSAGE":     "MESSAGE=\"hello world\"",
				"DESCRIPTION": "DESCRIPTION=\"This is a description with spaces\"",
			},
		},
		"Empty values need quotes": {
			input: map[string]string{
				"EMPTY": "",
			},
			want: map[string]string{
				"EMPTY": "EMPTY=\"\"",
			},
		},
		"Mixed values": {
			input: map[string]string{
				"SIMPLE":     "value",
				"WITH_SPACE": "value with space",
				"EMPTY":      "",
				"URL":        "https://example.com/path?query=value",
			},
			want: map[string]string{
				"SIMPLE":     "SIMPLE=value",
				"WITH_SPACE": "WITH_SPACE=\"value with space\"",
				"EMPTY":      "EMPTY=\"\"",
				"URL":        "URL=https://example.com/path?query=value",
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := marshalEnvWithoutQuotes(test.input)
			for key, expectedLine := range test.want {
				assert.Contains(t, got, expectedLine,
					fmt.Sprintf("Expected %s to be formatted as: %s", key, expectedLine))
			}
		})
	}
}

func TestMarshalEnvWithoutQuotes_AlphabeticalOrder(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input map[string]string
		want  string
	}{
		"Keys sorted alphabetically": {
			input: map[string]string{
				"ZEBRA":    "last",
				"APPLE":    "first",
				"MANGO":    "middle",
				"BANANA":   "second",
				"DATABASE": "db",
			},
			want: "APPLE=first\nBANANA=second\nDATABASE=db\nMANGO=middle\nZEBRA=last\n",
		},
		"Numbers and letters sorted": {
			input: map[string]string{
				"VAR_2": "two",
				"VAR_1": "one",
				"VAR_3": "three",
			},
			want: "VAR_1=one\nVAR_2=two\nVAR_3=three\n",
		},
		"Single key": {
			input: map[string]string{
				"ONLY_ONE": "value",
			},
			want: "ONLY_ONE=value\n",
		},
		"Empty map": {
			input: map[string]string{},
			want:  "",
		},
		"Keys with underscores sorted": {
			input: map[string]string{
				"Z_CONFIG":    "z",
				"A_SETTING":   "a",
				"M_VARIABLE":  "m",
				"B_PARAMETER": "b",
			},
			want: "A_SETTING=a\nB_PARAMETER=b\nM_VARIABLE=m\nZ_CONFIG=z\n",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := marshalEnvWithoutQuotes(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}
