package appdef

import (
	"github.com/goccy/go-json"
	"github.com/pkg/errors"
	"github.com/swaggest/jsonschema-go"
)

//go:generate go run ../../cmd/schema/main.go --output=../templates/schema.json

// GenerateSchema generates a JSON schema for the Definition type.
func GenerateSchema() ([]byte, error) {
	reflector := jsonschema.Reflector{}

	schema, err := reflector.Reflect(&Definition{})
	if err != nil {
		return nil, errors.Wrap(err, "reflecting schema")
	}

	// Add metadata
	schema.WithTitle("WebKit App Definition")
	schema.WithDescription("Schema for webkit app.json configuration files")
	schema.WithExtraPropertiesItem("$schema", "http://json-schema.org/draft-07/schema#")

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return nil, errors.Wrap(err, "marshalling schema to JSON")
	}

	return data, nil
}
