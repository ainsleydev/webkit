package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/internal/appdef"
)

// NewAppDefWithDefaults creates an appdef.Definition and applies defaults.
//
// This helper reduces boilerplate in tests by handling the common pattern
// of creating a definition and calling ApplyDefaults().
func NewAppDefWithDefaults(t *testing.T, def *appdef.Definition) *appdef.Definition {
	t.Helper()
	require.NoError(t, def.ApplyDefaults())
	return def
}
