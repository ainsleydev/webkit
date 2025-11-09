package appdef

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateSchema(t *testing.T) {
	t.Parallel()

	schema, err := GenerateSchema()
	require.NoError(t, err, "GenerateSchema should not return an error")
	require.NotNil(t, schema, "schema should not be nil")

	// Verify the schema is valid JSON
	var schemaMap map[string]any
	err = json.Unmarshal(schema, &schemaMap)
	require.NoError(t, err, "schema should be valid JSON")

	// Verify required top-level fields
	assert.Contains(t, schemaMap, "title", "schema should contain a title")
	assert.Contains(t, schemaMap, "description", "schema should contain a description")
	assert.Contains(t, schemaMap, "$schema", "schema should contain a $schema field")

	// Verify metadata values
	assert.Equal(t, "WebKit App Definition", schemaMap["title"], "title should be 'WebKit App Definition'")
	assert.Equal(t, "Schema for webkit app.json configuration files", schemaMap["description"], "description should match expected value")
	assert.Equal(t, "http://json-schema.org/draft-07/schema#", schemaMap["$schema"], "$schema should be draft-07")

	// Verify the schema defines properties for Definition
	assert.Contains(t, schemaMap, "properties", "schema should contain properties")
	properties, ok := schemaMap["properties"].(map[string]any)
	require.True(t, ok, "properties should be a map")

	// Verify key Definition fields are present
	expectedFields := []string{"webkit_version", "project", "apps", "resources", "shared"}
	for _, field := range expectedFields {
		assert.Contains(t, properties, field, "schema should define property %q", field)
	}
}

func TestGenerateSchema_SchemaStructure(t *testing.T) {
	t.Parallel()

	schema, err := GenerateSchema()
	require.NoError(t, err, "GenerateSchema should not return an error")

	var schemaMap map[string]any
	err = json.Unmarshal(schema, &schemaMap)
	require.NoError(t, err, "schema should be valid JSON")

	// Verify properties exist
	properties, ok := schemaMap["properties"].(map[string]any)
	require.True(t, ok, "properties should be a map")

	// Verify webkit_version property
	webkitVersion, ok := properties["webkit_version"].(map[string]any)
	require.True(t, ok, "webkit_version should be defined")
	assert.NotNil(t, webkitVersion, "webkit_version should not be nil")

	// Verify apps property is an array
	apps, ok := properties["apps"].(map[string]any)
	require.True(t, ok, "apps should be defined")
	assert.NotNil(t, apps, "apps should not be nil")

	// Verify project property is an object
	project, ok := properties["project"].(map[string]any)
	require.True(t, ok, "project should be defined")
	assert.NotNil(t, project, "project should not be nil")
}

func TestGenerateSchema_ValidJSON(t *testing.T) {
	t.Parallel()

	schema, err := GenerateSchema()
	require.NoError(t, err, "GenerateSchema should not return an error")

	// Verify the schema is not empty
	assert.Greater(t, len(schema), 0, "schema should not be empty")

	// Verify the schema is properly formatted JSON with indentation
	assert.Contains(t, string(schema), "\n", "schema should be indented/formatted")

	// Verify it can be unmarshaled and remarshaled
	var schemaMap map[string]any
	err = json.Unmarshal(schema, &schemaMap)
	require.NoError(t, err, "schema should be valid JSON")

	_, err = json.Marshal(schemaMap)
	require.NoError(t, err, "schema should be re-marshalable")
}
