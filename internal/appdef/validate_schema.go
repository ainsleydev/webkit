package appdef

import (
	_ "embed"
	"fmt"

	"github.com/kaptinlin/jsonschema"
)

//go:embed schema.json
var schemaJSON []byte

// ValidateAgainstSchema validates the provided JSON data against the embedded JSON schema.
// It returns a slice of validation errors if the data does not conform to the schema,
// or nil if validation passes.
//
// This function provides compile-time schema validation by:
//   - Using the embedded schema.json generated from the Go struct tags
//   - Validating required fields, types, patterns, and constraints
//   - Returning detailed error messages for schema violations
//
// Example usage:
//
//	data, _ := json.Marshal(definition)
//	if errs := ValidateAgainstSchema(data); errs != nil {
//	    for _, err := range errs {
//	        fmt.Println(err)
//	    }
//	}
func ValidateAgainstSchema(data []byte) []error {
	compiler := jsonschema.NewCompiler()

	// Load the embedded schema.
	schema, err := compiler.Compile(schemaJSON)
	if err != nil {
		return []error{fmt.Errorf("compiling schema: %w", err)}
	}

	// Validate the data against the schema.
	result := schema.ValidateJSON(data)
	if result.IsValid() {
		return nil
	}

	// Collect all validation errors.
	details := result.GetDetailedErrors()
	if len(details) == 0 {
		return []error{fmt.Errorf("schema validation failed with no details")}
	}

	var errs []error
	for _, detail := range details {
		errs = append(errs, fmt.Errorf("%s", detail))
	}

	return errs
}
