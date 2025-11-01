package appdef

import (
	"fmt"
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
		Staging:    EnvVar{"DEBUG": {Value: "true"}},
		Production: EnvVar{"DEBUG": {Value: "false"}},
	}

	var got []string
	e.Walk(func(entry EnvWalkEntry) {
		val := fmt.Sprintf("%s:%s=%v", entry.Environment, entry.Key, entry.Value)
		got = append(got, val)
	})

	want := []string{
		"development:DEBUG=true",
		"staging:DEBUG=true",
		"production:DEBUG=false",
	}

	assert.ElementsMatch(t, want, got)
}

func TestEnvironment_WalkE(t *testing.T) {
	t.Parallel()

	e := Environment{
		Dev:        EnvVar{"DEBUG": {Value: "true"}},
		Staging:    EnvVar{"DEBUG": {Value: "true"}},
		Production: EnvVar{"DEBUG": {Value: "false"}},
	}

	var got []string
	err := e.WalkE(func(entry EnvWalkEntry) error {
		if entry.Environment == env.Production {
			return fmt.Errorf("stop at production")
		}
		val := fmt.Sprintf("%s:%s=%v", entry.Environment, entry.Key, entry.Value)
		got = append(got, val)
		return nil
	})

	assert.Error(t, err)
	assert.ErrorContains(t, err, "production")
	assert.ElementsMatch(t, []string{
		"development:DEBUG=true",
		"staging:DEBUG=true",
	}, got)
}

func TestEnvironment_Walk_WithDefaults(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		env  Environment
		want []string
	}{
		"Default Only": {
			env: Environment{
				Default: EnvVar{"API_KEY": {Source: EnvSourceSOPS, Path: "secrets.yaml:API_KEY"}},
			},
			want: []string{
				"development:API_KEY=<nil>",
				"staging:API_KEY=<nil>",
				"production:API_KEY=<nil>",
			},
		},
		"Default With Override": {
			env: Environment{
				Default:    EnvVar{"API_KEY": {Source: EnvSourceSOPS}},
				Production: EnvVar{"API_KEY": {Source: EnvSourceValue, Value: "prod-key"}},
			},
			want: []string{
				"development:API_KEY=<nil>",  // from Default
				"staging:API_KEY=<nil>",      // from Default
				"production:API_KEY=<nil>",   // from Default
				"production:API_KEY=prod-key", // from Production (overrides)
			},
		},
		"Default Plus Specific": {
			env: Environment{
				Default: EnvVar{"SHARED_VAR": {Source: EnvSourceValue, Value: "shared"}},
				Dev:     EnvVar{"DEV_VAR": {Source: EnvSourceValue, Value: "dev-only"}},
			},
			want: []string{
				"development:SHARED_VAR=shared",
				"development:DEV_VAR=dev-only",
				"staging:SHARED_VAR=shared",
				"production:SHARED_VAR=shared",
			},
		},
		"Multiple Defaults Overridden": {
			env: Environment{
				Default: EnvVar{
					"VAR1": {Source: EnvSourceValue, Value: "default1"},
					"VAR2": {Source: EnvSourceValue, Value: "default2"},
				},
				Dev: EnvVar{
					"VAR1": {Source: EnvSourceValue, Value: "dev-override"},
				},
				Production: EnvVar{
					"VAR2": {Source: EnvSourceValue, Value: "prod-override"},
				},
			},
			want: []string{
				// Defaults applied to all environments
				"development:VAR1=default1",
				"development:VAR2=default2",
				"staging:VAR1=default1",
				"staging:VAR2=default2",
				"production:VAR1=default1",
				"production:VAR2=default2",
				// Environment-specific overrides
				"development:VAR1=dev-override",
				"production:VAR2=prod-override",
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var got []string
			test.env.Walk(func(entry EnvWalkEntry) {
				val := fmt.Sprintf("%s:%s=%v", entry.Environment, entry.Key, entry.Value)
				got = append(got, val)
			})

			assert.ElementsMatch(t, test.want, got)
		})
	}
}

func TestMergeVars(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		base     EnvVar
		override EnvVar
		want     EnvVar
	}{
		"Override Wins": {
			base:     EnvVar{"FOO": {Source: EnvSourceValue, Value: "bar"}},
			override: EnvVar{"FOO": {Source: EnvSourceValue, Value: "override"}},
			want:     EnvVar{"FOO": {Source: EnvSourceValue, Value: "override"}},
		},
		"Merge Both": {
			base:     EnvVar{"BAZ": {Source: EnvSourceValue, Value: "qux"}},
			override: EnvVar{"FOO": {Source: EnvSourceValue, Value: "bar"}},
			want: EnvVar{
				"BAZ": {Source: EnvSourceValue, Value: "qux"},
				"FOO": {Source: EnvSourceValue, Value: "bar"},
			},
		},
		"Nil Base": {
			base:     nil,
			override: EnvVar{"FOO": {Source: EnvSourceValue, Value: "val"}},
			want:     EnvVar{"FOO": {Source: EnvSourceValue, Value: "val"}},
		},
		"Nil Override": {
			base:     EnvVar{"FOO": {Source: EnvSourceValue, Value: "val"}},
			override: nil,
			want:     EnvVar{"FOO": {Source: EnvSourceValue, Value: "val"}},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := mergeVars(test.base, test.override)
			assert.Equal(t, test.want, got)
		})
	}
}
