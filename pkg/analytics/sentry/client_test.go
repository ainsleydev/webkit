package sentry

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ainsleydev/webkit/pkg/env"
)

func TestInitSentry(t *testing.T) {
	tests := map[string]struct {
		dsn       string
		dev       bool
		wantError bool
	}{
		"Valid DSN": {
			dsn:       "https://public@sentry.example.com/1",
			dev:       false,
			wantError: false,
		},
		"Invalid DSN": {
			dsn:       "invalid",
			dev:       false,
			wantError: true,
		},
		"Development Env": {
			dev:       true,
			wantError: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.dev {
				require.NoError(t, os.Setenv("APP_ENVIRONMENT", env.Development))
			}
			_, err := Init(test.dsn)
			assert.Equal(t, test.wantError, err != nil)
		})
	}
}
