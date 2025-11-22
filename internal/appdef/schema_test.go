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

func TestGenerateSchema_CommandsNoCircularRef(t *testing.T) {
	t.Parallel()

	schema, err := GenerateSchema()
	require.NoError(t, err)

	var schemaMap map[string]any
	require.NoError(t, json.Unmarshal(schema, &schemaMap))

	// Navigate to commands.additionalProperties.oneOf
	definitions := schemaMap["definitions"].(map[string]any)
	appDef := definitions["AppdefApp"].(map[string]any)
	commands := appDef["properties"].(map[string]any)["commands"].(map[string]any)
	additionalProps := commands["additionalProperties"].(map[string]any)
	oneOf := additionalProps["oneOf"].([]any)

	// Ensure no "$ref": "#" circular reference exists
	for i, item := range oneOf {
		if itemMap, ok := item.(map[string]any); ok {
			if ref, hasRef := itemMap["$ref"].(string); hasRef {
				assert.NotEqual(t, "#", ref, "oneOf[%d] has circular reference", i)
			}
		}
	}

	// Verify object option has inline properties
	objectOption := oneOf[2].(map[string]any)
	assert.Equal(t, "object", objectOption["type"])
	assert.Contains(t, objectOption, "properties")
}
