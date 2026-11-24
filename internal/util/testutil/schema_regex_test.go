package testutil

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStripLookarounds(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"No lookaround":        {input: `^[a-z]+$`, want: `^[a-z]+$`},
		"Negative lookahead":   {input: `^(?![0-9]+$)[^.][\s\S]*$`, want: `^[^.][\s\S]*$`},
		"Positive lookahead":   {input: `^(?=.*\d)\w+$`, want: `^\w+$`},
		"Lookbehinds":          {input: `(?<=a)b(?<!c)`, want: `b`},
		"Nested groups":        {input: `^(?!(a|b)+$)x`, want: `^x`},
		"Paren in class":       {input: `^(?![)(]+$)y`, want: `^y`},
		"Escaped paren":        {input: `^(?!\)+$)z`, want: `^z`},
		"Named group kept":     {input: `(?<name>a)(?P<n2>b)`, want: `(?<name>a)(?P<n2>b)`},
		"Escaped open paren":   {input: `\(?!x`, want: `\(?!x`},
		"Unbalanced untouched": {input: `^(?!abc`, want: `^(?!abc`},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, stripLookarounds(test.input))
		})
	}
}

func TestRewriteSchemaRegexes(t *testing.T) {
	t.Parallel()

	t.Run("Rewrites unsupported patterns", func(t *testing.T) {
		t.Parallel()

		in := []byte(`{
			"properties": {
				"name": {"type": "string", "pattern": "^(?!\\d)\\w+$"},
				"pattern": {"type": "string"}
			},
			"patternProperties": {
				"^(?![0-9]+$)[^.][\\s\\S]*$": {"type": "string"},
				"^x-": {"type": "number"}
			},
			"anyOf": [{"pattern": "^ok$"}]
		}`)

		var got map[string]any
		require.NoError(t, json.Unmarshal(rewriteSchemaRegexes(in), &got))

		name := got["properties"].(map[string]any)["name"].(map[string]any)
		assert.Equal(t, `^\w+$`, name["pattern"])

		props := got["patternProperties"].(map[string]any)
		assert.Contains(t, props, `^[^.][\s\S]*$`)
		assert.Contains(t, props, `^x-`)
		assert.Len(t, props, 2)

		assert.Equal(t, "^ok$", got["anyOf"].([]any)[0].(map[string]any)["pattern"])

		for p := range props {
			_, err := regexp.Compile(p)
			assert.NoError(t, err, p)
		}
	})

	t.Run("Unchanged when nothing to rewrite", func(t *testing.T) {
		t.Parallel()
		in := []byte(`{"pattern": "^[a-z]+$"}`)
		assert.Equal(t, in, rewriteSchemaRegexes(in))
	})

	t.Run("Non-JSON passed through", func(t *testing.T) {
		t.Parallel()
		in := []byte(`not json`)
		assert.Equal(t, in, rewriteSchemaRegexes(in))
	})
}
