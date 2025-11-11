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
	require.NoError(t, err)
	require.NotNil(t, schema)

	var schemaMap map[string]any
	err = json.Unmarshal(schema, &schemaMap)
	require.NoError(t, err)

	// Verify required top-level fields
	assert.Contains(t, schemaMap, "title")
	assert.Contains(t, schemaMap, "description")
	assert.Contains(t, schemaMap, "$schema")

	// Verify metadata values
	assert.Equal(t, "WebKit App Definition", schemaMap["title"])
	assert.Equal(t, "Schema for webkit app.json configuration files", schemaMap["description"])
	assert.Equal(t, "http://json-schema.org/draft-07/schema#", schemaMap["$schema"])

	// Verify the schema defines properties for Definition
	assert.Contains(t, schemaMap, "properties")
	properties, ok := schemaMap["properties"].(map[string]any)
	require.True(t, ok)

	// Verify key Definition fields are present
	expectedFields := []string{"webkit_version", "project", "apps", "resources", "shared"}
	for _, field := range expectedFields {
		assert.Contains(t, properties, field, "schema should define property %q", field)
	}
}
