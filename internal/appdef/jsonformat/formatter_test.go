package jsonformat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"Simple environment variable with source and value": {
			input: `{
	"env": {
		"dev": {
			"DATABASE_URI": {
				"source": "value",
				"value": "file:./cms.db"
			}
		}
	}
}`,
			want: `{
	"env": {
		"dev": {
			"DATABASE_URI": {"source": "value", "value": "file:./cms.db"}
		}
	}
}`,
		},
		"Multiple environment variables": {
			input: `{
	"env": {
		"dev": {
			"DATABASE_URI": {
				"source": "value",
				"value": "file:./cms.db"
			},
			"FRONTEND_URL": {
				"source": "value",
				"value": "http://localhost:5173"
			}
		}
	}
}`,
			want: `{
	"env": {
		"dev": {
			"DATABASE_URI": {"source": "value", "value": "file:./cms.db"},
			"FRONTEND_URL": {"source": "value", "value": "http://localhost:5173"}
		}
	}
}`,
		},
		"Simple command": {
			input: `{
	"commands": {
		"build": {
			"command": "pnpm build"
		}
	}
}`,
			want: `{
	"commands": {
		"build": {"command": "pnpm build"}
	}
}`,
		},
		"Multiple commands": {
			input: `{
	"commands": {
		"build": {
			"command": "pnpm build"
		},
		"test": {
			"command": "pnpm test"
		}
	}
}`,
			want: `{
	"commands": {
		"build": {"command": "pnpm build"},
		"test": {"command": "pnpm test"}
	}
}`,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := Format([]byte(test.input))
			require.NoError(t, err)
			assert.Equal(t, test.want, string(got))
		})
	}
}

func TestFormat_WithRealJSON(t *testing.T) {
	// Test with actual JSON marshaling to ensure it works end-to-end.
	type envValue struct {
		Source string `json:"source"`
		Value  string `json:"value,omitempty"`
	}

	type commandSpec struct {
		Command string `json:"command"`
	}

	type app struct {
		Name     string                 `json:"name"`
		Env      map[string]map[string]envValue `json:"env"`
		Commands map[string]commandSpec `json:"commands"`
	}

	testApp := app{
		Name: "test",
		Env: map[string]map[string]envValue{
			"dev": {
				"DB_URL": {Source: "value", Value: "localhost"},
			},
		},
		Commands: map[string]commandSpec{
			"build": {Command: "make"},
		},
	}

	data, err := json.MarshalIndent(testApp, "", "\t")
	require.NoError(t, err)

	formatted, err := Format(data)
	require.NoError(t, err)

	result := string(formatted)
	assert.Contains(t, result, `"DB_URL": {"source": "value", "value": "localhost"}`)
	assert.Contains(t, result, `"build": {"command": "make"}`)
}
