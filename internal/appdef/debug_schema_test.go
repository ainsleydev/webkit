package appdef

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// TestSchemaItself validates that schema.json is valid JSON Schema (WRONG WAY)
func TestSchemaItself(t *testing.T) {
	compiler := jsonschema.NewCompiler()

	// This is WRONG - passing raw bytes
	if err := compiler.AddResource("schema.json", schemaJSON); err != nil {
		t.Fatalf("Failed to add schema resource: %v", err)
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		t.Fatalf("Failed to compile schema: %v", err)
	}

	fmt.Printf("Schema compiled successfully: %+v\n", schema)
}

// TestSchemaWithExplicitDraft tests if specifying draft-07 helps (WRONG WAY)
func TestSchemaWithExplicitDraft(t *testing.T) {
	compiler := jsonschema.NewCompiler()
	compiler.DefaultDraft(jsonschema.Draft7)

	// This is WRONG - passing raw bytes
	if err := compiler.AddResource("schema.json", schemaJSON); err != nil {
		t.Fatalf("Failed to add schema resource: %v", err)
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		t.Fatalf("Failed to compile schema with draft-7: %v", err)
	}

	fmt.Printf("Schema compiled successfully with draft-7: %+v\n", schema)
}

// TestSchemaWithUnmarshal tests if unmarshaling the JSON first helps (CORRECT WAY)
func TestSchemaWithUnmarshal(t *testing.T) {
	compiler := jsonschema.NewCompiler()

	// Unmarshal the JSON first as recommended by the library
	doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaJSON))
	if err != nil {
		t.Fatalf("Failed to unmarshal schema JSON: %v", err)
	}

	t.Logf("Unmarshaled schema type: %T", doc)

	// Add the unmarshaled document
	if err := compiler.AddResource("schema.json", doc); err != nil {
		t.Fatalf("Failed to add schema resource: %v", err)
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		t.Fatalf("Failed to compile schema: %v", err)
	}

	t.Logf("Schema compiled successfully!")
	fmt.Printf("Schema: %+v\n", schema)
}

// TestValidationWithCorrectApproach validates JSON data with the correct approach
func TestValidationWithCorrectApproach(t *testing.T) {
	testJSON := `{
		"webkit_version": "1.0.0",
		"project": {
			"name": "test-project",
			"title": "Test Project",
			"description": "A comprehensive test project with all fields",
			"repo": {
				"owner": "testowner",
				"name": "testrepo"
			}
		},
		"notifications": {
			"webhook_url": "https://hooks.slack.com/services/XXX/YYY/ZZZ"
		},
		"shared": {
			"env": {
				"default": {
					"API_URL": {
						"source": "value",
						"value": "https://api.example.com"
					}
				}
			}
		},
		"resources": [{
			"name": "test-db",
			"type": "postgres",
			"provider": "digitalocean"
		}],
		"apps": [{
			"name": "test-app",
			"title": "Test App",
			"type": "golang",
			"description": "A test application",
			"path": "./apps/test",
			"build": {
				"dockerfile": "Dockerfile",
				"port": 8080
			},
			"infra": {
				"provider": "digitalocean",
				"type": "vm"
			},
			"domains": [{
				"name": "test.example.com",
				"type": "primary"
			}]
		}]
	}`

	compiler := jsonschema.NewCompiler()

	// CORRECT WAY: Unmarshal schema first
	schemaDoc, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaJSON))
	if err != nil {
		t.Fatalf("Failed to unmarshal schema: %v", err)
	}

	if err := compiler.AddResource("schema.json", schemaDoc); err != nil {
		t.Fatalf("Failed to add schema resource: %v", err)
	}

	schema, err := compiler.Compile("schema.json")
	if err != nil {
		t.Fatalf("Failed to compile schema: %v", err)
	}

	// Unmarshal the test data
	var testData any
	if err := json.Unmarshal([]byte(testJSON), &testData); err != nil {
		t.Fatalf("Failed to unmarshal test data: %v", err)
	}

	// Validate
	if err := schema.Validate(testData); err != nil {
		t.Fatalf("Validation failed: %v", err)
	}

	t.Log("Validation passed!")
}

// TestPrintSchemaBytes prints the schema to see what we're working with
func TestPrintSchemaBytes(t *testing.T) {
	t.Logf("Schema length: %d bytes", len(schemaJSON))
	t.Logf("Schema type: %T", schemaJSON)
	t.Logf("First 100 chars: %s", string(schemaJSON[:min(100, len(schemaJSON))]))

	// Check if it's valid JSON
	if schemaJSON[0] != '{' {
		t.Errorf("Schema doesn't start with '{', starts with: %c (%d)", schemaJSON[0], schemaJSON[0])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
