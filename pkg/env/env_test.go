package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	tt := map[string]struct {
		input any
		want  string
	}{
		"Existing Key": {
			input: "ENV_KEY",
			want:  "existing_value",
		},
		"Non Existing Key": {
			input: "NON_EXISTING_KEY",
			want:  "",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			err := os.WriteFile(".env", []byte("ENV_KEY=existing_value"), os.ModePerm)
			require.NoError(t, err)
			t.Cleanup(func() {
				require.NoError(t, os.Remove(".env"))
			})
			got := ParseConfig(test.input)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGet_WithExistingKey(t *testing.T) {
	tt := map[string]struct {
		input string
		want  string
	}{
		"Existing Key": {
			input: "ENV_KEY",
			want:  "existing_value",
		},
		"Fallback": {
			input: "ENV_KEY",
			want:  "fallback_value",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			if test.want != "" {
				t.Setenv(test.input, test.want)
			}
			got := Get(test.input, "fallback_value")
			assert.Equal(t, test.want, got)
		})
	}
}

func TestGetOrError(t *testing.T) {
	tt := map[string]struct {
		input   string
		wantErr bool
		want    string
	}{
		"Existing Key": {
			input:   "EXISTING_KEY",
			wantErr: false,
			want:    "existing_value",
		},
		"Non Existing Key": {
			input:   "NON_EXISTING_KEY",
			wantErr: true,
			want:    "",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			if test.want != "" {
				t.Setenv(test.input, test.want)
			}
			got, err := GetOrError(test.input)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestAppEnvironment(t *testing.T) {
	t.Setenv("APP_ENVIRONMENT", "edge_case_value")
	got := AppEnvironment()
	want := "edge_case_value"
	assert.Equal(t, want, got)
}

func TestIsDevelopment(t *testing.T) {
	t.Run("Development", func(t *testing.T) {
		t.Setenv("APP_ENVIRONMENT", Development)
		got := IsDevelopment()
		assert.True(t, got)
	})
	t.Run("Empty", func(t *testing.T) {
		t.Setenv("APP_ENVIRONMENT", "")
		got := IsDevelopment()
		assert.True(t, got)
	})
}

func TestIsStaging(t *testing.T) {
	t.Setenv("APP_ENVIRONMENT", Staging)
	got := IsStaging()
	assert.True(t, got)
}

func TestIsProduction(t *testing.T) {
	t.Setenv("APP_ENVIRONMENT", Production)
	got := IsProduction()
	assert.True(t, got)
}
