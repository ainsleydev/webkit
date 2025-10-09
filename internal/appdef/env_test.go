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
