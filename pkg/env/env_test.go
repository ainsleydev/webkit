package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfig(t *testing.T) {
	type env struct {
		Key string `env:"ENV_KEY,required"`
	}

	err := os.WriteFile(".env", []byte("ENV_KEY=test"), os.ModePerm)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, os.Remove(".env"))
	})

	t.Run("No File", func(t *testing.T) {
		err := ParseConfig(&env{}, "wrong")
		require.Error(t, err)
	})

	t.Run("Parse error", func(t *testing.T) {
		err = ParseConfig("wrong", ".env")
		require.Error(t, err)
	})

	t.Run("OK", func(t *testing.T) {
		var cfg env
		err = ParseConfig(&cfg, ".env")
		require.NoError(t, err)
		assert.Equal(t, "test", cfg.Key)
	})
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

func TestGet(t *testing.T) {
	fallback := "fallback"

	t.Run("Empty", func(t *testing.T) {
		got := Get("", fallback)
		assert.Equal(t, fallback, got)
	})

	t.Run("OK", func(t *testing.T) {
		t.Setenv("EXISTING_KEY", "existing_value")
		got := Get("EXISTING_KEY", fallback)
		assert.Equal(t, "existing_value", got)
	})
}

func TestGetOrError(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		got, err := GetOrError("NON_EXISTING_KEY")
		assert.Error(t, err)
		assert.Empty(t, got)
	})

	t.Run("OK", func(t *testing.T) {
		t.Setenv("EXISTING_KEY", "existing_value")
		got, err := GetOrError("EXISTING_KEY")
		assert.NoError(t, err)
		assert.Equal(t, "existing_value", got)
	})
}

func TestAppEnvironment(t *testing.T) {
	t.Setenv(AppEnvironmentKey, "edge_case_value")
	got := AppEnvironment()
	want := "edge_case_value"
	assert.Equal(t, want, got)
}

func TestIsDevelopment(t *testing.T) {
	t.Run("Development", func(t *testing.T) {
		t.Setenv(AppEnvironmentKey, Development)
		got := IsDevelopment()
		assert.True(t, got)
	})
	t.Run("Empty", func(t *testing.T) {
		t.Setenv(AppEnvironmentKey, "")
		got := IsDevelopment()
		assert.True(t, got)
	})
}

func TestIsStaging(t *testing.T) {
	t.Setenv(AppEnvironmentKey, Staging)
	got := IsStaging()
	assert.True(t, got)
}

func TestIsProduction(t *testing.T) {
	t.Setenv(AppEnvironmentKey, Production)
	got := IsProduction()
	assert.True(t, got)
}
