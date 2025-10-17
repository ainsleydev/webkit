package infra

import (
	"testing"

	"github.com/ainsleydev/webkit/internal/appdef"
	"github.com/ainsleydev/webkit/pkg/env"
)

func TestTFVarsFromDefinition(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		definition appdef.Definition
		env        env.Environment
	}{
		"Resource": {},
	}

	for name, _ := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

		})
	}
}
