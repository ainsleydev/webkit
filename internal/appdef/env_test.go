package appdef

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	e.Walk(func(envName string, name string, value EnvValue) {
		got = append(got, envName+":"+name+"="+value.Value)
	})

	want := []string{
		"development:DEBUG=true",
		"staging:DEBUG=true",
		"production:DEBUG=false",
	}

	assert.ElementsMatch(t, want, got)
}

func TestMergeVars(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		base     EnvVar
		override EnvVar
		want     EnvVar
	}{
		"Both Non Empty": {
			base: EnvVar{
				"FOO": {Source: EnvSourceValue, Value: "bar"},
				"BAZ": {Source: EnvSourceValue, Value: "qux"},
			},
			override: EnvVar{
				"FOO": {Source: EnvSourceValue, Value: "override"},
				"NEW": {Source: EnvSourceValue, Value: "val"},
			},
			want: EnvVar{
				"FOO": {Source: EnvSourceValue, Value: "override"},
				"BAZ": {Source: EnvSourceValue, Value: "qux"},
				"NEW": {Source: EnvSourceValue, Value: "val"},
			},
		},
		"Base Empty, Override Non Empty": {
			base: EnvVar{},
			override: EnvVar{
				"FOO": {Source: EnvSourceValue, Value: "val"},
			},
			want: EnvVar{
				"FOO": {Source: EnvSourceValue, Value: "val"},
			},
		},
		"Base Non Empty, Override Empty": {
			base: EnvVar{
				"FOO": {Source: EnvSourceValue, Value: "val"},
			},
			override: EnvVar{},
			want: EnvVar{
				"FOO": {Source: EnvSourceValue, Value: "val"},
			},
		},
		"Both Empty": {
			base:     EnvVar{},
			override: EnvVar{},
			want:     EnvVar{},
		},
		"Nil Base, Non Empty Override": {
			base: nil,
			override: EnvVar{
				"FOO": {Source: EnvSourceValue, Value: "val"},
			},
			want: EnvVar{
				"FOO": {Source: EnvSourceValue, Value: "val"},
			},
		},
		"Nil Override, NonEmpty Base": {
			base: EnvVar{
				"FOO": {Source: EnvSourceValue, Value: "val"},
			},
			override: nil,
			want: EnvVar{
				"FOO": {Source: EnvSourceValue, Value: "val"},
			},
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
