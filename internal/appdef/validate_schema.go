package appdef

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v6"
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

	// Unmarshal the schema JSON into a document.
	// The library requires a parsed JSON value, not raw bytes.
	schemaDoc, err := jsonschema.UnmarshalJSON(bytes.NewReader(schemaJSON))
	if err != nil {
		return []error{fmt.Errorf("unmarshaling schema: %w", err)}
	}

	// Add the schema to the compiler.
	if err := compiler.AddResource("schema.json", schemaDoc); err != nil {
		return []error{fmt.Errorf("adding schema resource: %w", err)}
	}

	// Compile the schema.
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return []error{fmt.Errorf("compiling schema: %w", err)}
	}

	// Unmarshal the JSON data into a generic interface.
	var v any
	if err = json.Unmarshal(data, &v); err != nil {
		return []error{fmt.Errorf("unmarshaling JSON: %w", err)}
	}

	// Validate the data against the schema.
	if err = schema.Validate(v); err != nil {
		// Convert validation error to slice.
		if ve, ok := err.(*jsonschema.ValidationError); ok {
			var errs []error
			collectErrors(ve, &errs)
			return errs
		}
		return []error{err}
	}

	return nil
}

// collectErrors recursively collects all validation errors from the error tree.
func collectErrors(ve *jsonschema.ValidationError, errs *[]error) {
	if len(ve.Causes) == 0 {
		*errs = append(*errs, fmt.Errorf("%s: %s", ve.InstanceLocation, ve.Error()))
		return
	}
	for _, cause := range ve.Causes {
		collectErrors(cause, errs)
	}
}
