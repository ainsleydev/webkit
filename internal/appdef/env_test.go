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
