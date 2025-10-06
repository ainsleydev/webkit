package testutil

import (
	"encoding/json"
	"errors"
	"fmt"
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
