package appdef

import (
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/swaggest/jsonschema-go"
)

// Generate it in the templates for copying it over to WebKit repos
// and the root so it's available to latch onto publicly.

//go:generate go run ../../cmd/schema/main.go --output=../templates/schema.json
//go:generate go run ../../cmd/schema/main.go --output=../../schema.json

// GenerateSchema generates a JSON schema for the Definition type.
func GenerateSchema() ([]byte, error) {
	reflector := jsonschema.Reflector{}

	schema, err := reflector.Reflect(&Definition{})
	if err != nil {
		return nil, errors.Wrap(err, "reflecting schema")
	}

	// Add metadata.
	schema.WithTitle("WebKit App Definition")
	schema.WithDescription("Schema for webkit app.json configuration files")
	schema.WithExtraPropertiesItem("$schema", "http://json-schema.org/draft-07/schema#")

	// Marshal to JSON first.
	data, err := json.Marshal(schema)
	if err != nil {
		return nil, errors.Wrap(err, "marshalling schema to JSON")
	}

	// Post-process the schema to fix OrderedMap type references.
	data, err = fixOrderedMapSchema(data)
	if err != nil {
		return nil, errors.Wrap(err, "fixing OrderedMap schema")
	}

	// Re-marshal with indentation for readability.
	var schemaMap map[string]any
	if err := json.Unmarshal(data, &schemaMap); err != nil {
		return nil, errors.Wrap(err, "unmarshalling schema for formatting")
	}

	data, err = json.MarshalIndent(schemaMap, "", "\t")
	if err != nil {
		return nil, errors.Wrap(err, "marshalling schema with indentation")
	}

	return data, nil
}

// fixOrderedMapSchema removes the TypesOrderedMap definition and replaces
// references to it with an inline object schema.
func fixOrderedMapSchema(data []byte) ([]byte, error) {
	var schemaMap map[string]any
	if err := json.Unmarshal(data, &schemaMap); err != nil {
		return nil, err
	}

	// Find and remove the TypesOrderedMap definition.
	definitions, ok := schemaMap["definitions"].(map[string]any)
	if !ok {
		return data, nil
	}

	// Find the OrderedMap definition key.
	var orderedMapKey string
	for key := range definitions {
		if len(key) > 15 && key[:15] == "TypesOrderedMap" {
			orderedMapKey = key
			break
		}
	}

	if orderedMapKey == "" {
		return data, nil // No OrderedMap found, nothing to fix.
	}

	// Get the value schema from the OrderedMap definition.
	orderedMapDef, ok := definitions[orderedMapKey].(map[string]any)
	if !ok {
		return data, nil
	}

	// Create the replacement schema for commands.
	commandsSchema := map[string]any{
		"type":        "object",
		"description": "Custom commands for linting, testing, formatting, and building",
	}

	// Get the additionalProperties from the OrderedMap definition.
	if additionalProps, ok := orderedMapDef["additionalProperties"]; ok {
		commandsSchema["additionalProperties"] = additionalProps
	}

	// Walk the schema and replace all references to the OrderedMap.
	replaceRefs(schemaMap, "#/definitions/"+orderedMapKey, commandsSchema)

	// Remove the OrderedMap definition.
	delete(definitions, orderedMapKey)

	// Fix circular references in oneOf schemas.
	// The CommandSpec oneOf contains "$ref": "#" which incorrectly points to root schema.
	fixCircularOneOfRefs(schemaMap)

	// Marshal the updated schema.
	return json.Marshal(schemaMap)
}

// replaceRefs recursively walks through the schema and replaces $ref values.
func replaceRefs(data any, targetRef string, replacement map[string]any) {
	switch v := data.(type) {
	case map[string]any:
		// Check if this is a $ref that needs to be replaced.
		if ref, ok := v["$ref"].(string); ok && ref == targetRef {
			// Remove the $ref and add the replacement properties.
			delete(v, "$ref")
			for key, val := range replacement {
				v[key] = val
			}
			return
		}
		// Recursively process all map values.
		for _, val := range v {
			replaceRefs(val, targetRef, replacement)
		}
	case []any:
		// Recursively process all array elements.
		for _, val := range v {
			replaceRefs(val, targetRef, replacement)
		}
	}
}

// fixCircularOneOfRefs finds and fixes circular "$ref": "#" references in oneOf schemas.
// These occur when CommandSpec.JSONSchemaOneOf() creates a reference to the root schema
// instead of inlining the object schema.
func fixCircularOneOfRefs(data any) {
	switch v := data.(type) {
	case map[string]any:
		// Check if this object has a oneOf with a circular reference.
		if oneOf, ok := v["oneOf"].([]any); ok {
			// Also check if there are properties defined (CommandSpec case).
			if properties, hasProps := v["properties"].(map[string]any); hasProps {
				// Look for "$ref": "#" in the oneOf array.
				for i, item := range oneOf {
					if itemMap, ok := item.(map[string]any); ok {
						if ref, ok := itemMap["$ref"].(string); ok && ref == "#" {
							// Replace the circular reference with an inline object schema.
							oneOf[i] = map[string]any{
								"type":       "object",
								"properties": properties,
							}
						}
					}
				}
			}
		}

		// Recursively process all map values.
		for _, val := range v {
			fixCircularOneOfRefs(val)
		}
	case []any:
		// Recursively process all array elements.
		for _, val := range v {
			fixCircularOneOfRefs(val)
		}
	}
}
