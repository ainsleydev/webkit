package testutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/kaptinlin/jsonschema"
	"github.com/stretchr/testify/require"
)

// SchemaValidator wraps a compiled JSON Schema for validation.
type SchemaValidator struct {
	schema *jsonschema.Schema
}

// SchemaFromURL fetches a JSON Schema from the given URL and compiles it.
func SchemaFromURL(t *testing.T, url string) (*SchemaValidator, error) {
	t.Helper()

	compiler := jsonschema.NewCompiler()
	re2CompatibleLoaders(compiler)
	schema, err := compiler.GetSchema(url)
	require.NoError(t, err)

	return &SchemaValidator{schema: schema}, nil
}

// SchemaFromBytes compiles a JSON Schema from a byte slice.
func SchemaFromBytes(t *testing.T, data []byte) (*SchemaValidator, error) {
	t.Helper()

	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile(data)
	require.NoError(t, err)

	return &SchemaValidator{schema: schema}, nil
}

// ValidateYAML validates YAML bytes against the schema.
func (v *SchemaValidator) ValidateYAML(yamlData []byte) error {
	jsonData, err := yaml.YAMLToJSON(yamlData)
	if err != nil {
		return fmt.Errorf("failed to convert YAML to JSON: %w", err)
	}
	return v.ValidateJSON(jsonData)
}

// ValidateJSON validates JSON bytes against the schema.
func (v *SchemaValidator) ValidateJSON(jsonData []byte) error {
	return handleResult(v.schema.ValidateJSON(jsonData))
}

// ValidateMap validates a map against the schema.
func (v *SchemaValidator) ValidateMap(data map[string]any) error {
	return handleResult(v.schema.ValidateMap(data))
}

// Validate validates any data against the schema.
func (v *SchemaValidator) Validate(data any) error {
	return handleResult(v.schema.Validate(data))
}

func handleResult(result *jsonschema.EvaluationResult) error {
	if result.IsValid() {
		return nil
	}

	// Only get errors if actually invalid
	details := result.GetDetailedErrors()
	if len(details) > 0 {
		indent, _ := json.MarshalIndent(details, "", "  ")
		return fmt.Errorf("schema validation failed:\n%s", indent)
	}

	return errors.New("schema is invalid")
}

// re2CompatibleLoaders wraps the compiler's HTTP(S) loaders so that every
// fetched schema, including remote $refs, has its regexes made compatible
// with Go's RE2 engine before compiling.
//
// Third-party schemas (e.g. schemastore.org) use ECMA-262 lookarounds such
// as (?!...) that RE2 can't parse, which fails validation outright.
func re2CompatibleLoaders(c *jsonschema.Compiler) {
	for _, scheme := range []string{"http", "https"} {
		next, ok := c.Loaders[scheme]
		if !ok {
			continue
		}
		c.RegisterLoader(scheme, func(url string) (io.ReadCloser, error) {
			body, err := next(url)
			if err != nil {
				return nil, err
			}
			defer body.Close() //nolint:errcheck

			data, err := io.ReadAll(body)
			if err != nil {
				return nil, err
			}
			return io.NopCloser(bytes.NewReader(rewriteSchemaRegexes(data))), nil
		})
	}
}

// rewriteSchemaRegexes rewrites "pattern" values and "patternProperties"
// keys that RE2 can't compile by removing their lookaround groups. The
// data is returned unchanged if it isn't JSON or nothing needs rewriting.
func rewriteSchemaRegexes(data []byte) []byte {
	var doc any
	if err := json.Unmarshal(data, &doc); err != nil {
		return data
	}
	if !rewriteRegexes(doc) {
		return data
	}
	out, err := json.Marshal(doc)
	if err != nil {
		return data
	}
	return out
}

func rewriteRegexes(node any) bool {
	changed := false
	switch v := node.(type) {
	case map[string]any:
		for key, val := range v {
			switch key {
			case "pattern":
				if p, ok := val.(string); ok {
					if fixed, ok := re2Compatible(p); ok {
						v[key] = fixed
						changed = true
						continue
					}
				}
			case "patternProperties":
				if props, ok := val.(map[string]any); ok {
					for p, sub := range props {
						if fixed, ok := re2Compatible(p); ok {
							delete(props, p)
							props[fixed] = sub
							changed = true
						}
					}
				}
			}
			if rewriteRegexes(val) {
				changed = true
			}
		}
	case []any:
		for _, val := range v {
			if rewriteRegexes(val) {
				changed = true
			}
		}
	}
	return changed
}

// re2Compatible returns p without lookaround groups when p doesn't
// compile with RE2 but the stripped version does. A lookaround only adds
// a condition, so removing it gives a pattern that matches a superset of
// the original: validation gets slightly looser, never stricter.
func re2Compatible(p string) (string, bool) {
	if _, err := regexp.Compile(p); err == nil {
		return "", false
	}
	stripped := stripLookarounds(p)
	if stripped == p {
		return "", false
	}
	if _, err := regexp.Compile(stripped); err != nil {
		return "", false
	}
	return stripped, true
}

// stripLookarounds removes (?=...), (?!...), (?<=...) and (?<!...) groups
// from a regex, respecting escapes and character classes.
func stripLookarounds(p string) string {
	var out strings.Builder
	for i := 0; i < len(p); {
		if end, ok := lookaroundEnd(p, i); ok {
			i = end
			continue
		}
		next := regexTokenEnd(p, i)
		out.WriteString(p[i:next])
		i = next
	}
	return out.String()
}

// lookaroundEnd reports whether a lookaround group starts at i and, if
// so, returns the index just past its closing parenthesis.
func lookaroundEnd(p string, i int) (int, bool) {
	rest := p[i:]
	if !strings.HasPrefix(rest, "(?=") && !strings.HasPrefix(rest, "(?!") &&
		!strings.HasPrefix(rest, "(?<=") && !strings.HasPrefix(rest, "(?<!") {
		return 0, false
	}
	depth := 0
	for j := i; j < len(p); {
		switch p[j] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return j + 1, true
			}
		}
		j = regexTokenEnd(p, j)
	}
	return 0, false // Unbalanced; leave the pattern alone.
}

// regexTokenEnd returns the index just past the token starting at i,
// treating an escape sequence or a whole character class as one token.
func regexTokenEnd(p string, i int) int {
	switch p[i] {
	case '\\':
		return min(i+2, len(p))
	case '[':
		j := i + 1
		if j < len(p) && p[j] == '^' {
			j++
		}
		if j < len(p) && p[j] == ']' {
			j++ // A leading ] is a literal.
		}
		for j < len(p) {
			switch p[j] {
			case '\\':
				j += 2
				continue
			case ']':
				return j + 1
			}
			j++
		}
		return len(p)
	default:
		return i + 1
	}
}
