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
