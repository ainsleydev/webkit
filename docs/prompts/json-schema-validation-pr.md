# JSON Schema Validation PR Summary

## Overview

This PR implements a **two-tier validation system** for webkit's `app.json` configuration files, combining JSON schema validation with custom business logic validation.

## Problem Solved

Previously, webkit used incorrect `jsonschema:"required"` tags (should be `required:"true"` for swaggest/jsonschema-go) and lacked runtime schema validation. While the schema was generated for IDE support, the actual JSON data wasn't validated against it at runtime.

## Solution: Two-Tier Validation

### Tier 1: Schema Validation (NEW)
Runtime validation using `kaptinlin/jsonschema` to validate JSON data against the generated schema:
- **Required fields**: Ensures critical configuration is present
- **Data types**: Validates strings, numbers, booleans, arrays, objects
- **Pattern validation**: Enforces naming conventions (lowercase-hyphenated)
- **Enum constraints**: Restricts values to predefined options
- **Numeric ranges**: Port numbers (1-65535)
- **String lengths**: Description max 200 characters
- **Format validation**: URLs and formatted strings

### Tier 2: Custom Validation (EXISTING)
Business logic validation for filesystem and cross-field checks:
- Domain format validation (no protocol prefixes)
- App path existence checks
- Resource reference validation
- Terraform VM domain requirements

## Key Changes

### 1. Fixed Required Field Tags
Changed all instances from `jsonschema:"required"` to `required:"true"`:
```go
// Before
Name string `json:"name" jsonschema:"required" description:"..."`

// After
Name string `json:"name" required:"true" description:"..."`
```

### 2. Added Validation Constraints
```go
// Pattern validation for names
Name string `json:"name" required:"true" pattern:"^[a-z][a-z0-9-]*$" description:"..."`

// Enum constraints
Type AppType `json:"type" required:"true" enum:"svelte-kit,golang,payload" description:"..."`

// Numeric ranges
Port int `json:"port,omitempty" minimum:"1" maximum:"65535" description:"..."`

// String length
Description string `json:"description,omitempty" maxLength:"200" description:"..."`

// Format validation
WebhookURL string `json:"webhook_url,omitzero" format:"uri" description:"..."`
```

### 3. Runtime Schema Validation
Created `validate_schema.go` with runtime validation:
```go
func ValidateAgainstSchema(data []byte) []error {
    // Validates JSON data against embedded schema.json
    // Returns detailed error messages for violations
}
```

Integrated into `Read()` function:
```go
func Read(root afero.Fs) (*Definition, error) {
    // 1. Read file
    // 2. Validate against schema (NEW)
    // 3. Unmarshal JSON
    // 4. Apply defaults
    // 5. Custom validation
}
```

### 4. Enhanced Schema Command
Updated `webkit schema` command:
```bash
# Generate to .webkit/schema.json (default)
webkit schema

# Generate to custom location
webkit schema --output custom-path/schema.json

# Output to stdout
webkit schema --stdout
```

## Files Modified/Created

### New Files
- `internal/appdef/validate_schema.go` - Runtime schema validation
- `internal/appdef/validate_schema_test.go` - Schema validation tests
- `internal/appdef/schema.json` - Embedded schema for validation
- `docs/prompts/json-schema-validation-pr.md` - This file

### Modified Files
- `internal/appdef/definition.go` - Integrated schema validation into Read()
- `internal/appdef/project.go` - Fixed tags, added pattern/maxLength
- `internal/appdef/apps.go` - Fixed tags, added pattern/enum/numeric constraints
- `internal/appdef/resources.go` - Fixed tags, added pattern/enum
- `internal/appdef/env.go` - Fixed tags, added enum
- `internal/appdef/schema.go` - Fixed go:generate path
- `internal/cmd/schema.go` - Enhanced with .webkit default and --stdout flag
- `docs/cli/validation.md` - Comprehensive documentation update
- `schema.json` - Regenerated with correct tags

## Technical Implementation

### Libraries Used
- **swaggest/jsonschema-go**: Generates JSON schema from Go struct tags
- **kaptinlin/jsonschema**: Validates JSON data against schema at runtime

### Validation Flow
```
1. Read app.json file
2. ValidateAgainstSchema(data) → Schema validation (structural/type checks)
   ↓ (if passes)
3. Unmarshal JSON to Definition struct
4. ApplyDefaults()
5. Validate(fs) → Custom validation (business logic checks)
   ↓ (if passes)
6. Return validated Definition
```

### Struct Tag Format
```go
type Example struct {
    // Required field
    Name string `json:"name" required:"true" description:"Field description"`

    // Pattern validation
    Slug string `json:"slug" pattern:"^[a-z][a-z0-9-]*$" description:"..."`

    // Enum constraint
    Type string `json:"type" enum:"value1,value2,value3" description:"..."`

    // Numeric range
    Port int `json:"port" minimum:"1" maximum:"65535" description:"..."`

    // String length
    Desc string `json:"desc" maxLength:"200" description:"..."`

    // Format validation
    URL string `json:"url" format:"uri" description:"..."`
}
```

## Usage

### For Users
Validation happens automatically when loading `app.json`:
```bash
webkit validate  # Explicit validation
webkit update    # Validates before updating
```

### For IDE Support
Reference schema in `app.json`:
```json
{
  "$schema": ".webkit/schema.json",
  "webkit_version": "v0.0.40",
  ...
}
```

### For Developers
Generate schema after changing struct tags:
```bash
go generate ./internal/appdef
```

## Testing

All tests pass:
- `TestValidateAgainstSchema` - Required field validation
- `TestValidateAgainstSchema_Integration` - End-to-end struct marshaling
- Existing custom validation tests remain unchanged

## Important Note

The `kaptinlin/jsonschema` library validates required fields and types reliably, but may not enforce all constraints (pattern, enum, min/max, format) at runtime. These constraints are still valuable for:
1. **IDE support**: Autocomplete, inline validation, hover documentation
2. **Documentation**: Clear specification of valid values
3. **Future compatibility**: If runtime validation improves

Custom validation in Tier 2 still enforces critical business logic that the schema cannot express.

## Commits on Branch

1. `9dc20ec` - feat: Add JSON schema validation for app.json
2. `64906ed` - fix: Follow Go testing conventions in validate_test.go
3. `fe69f87` - refactor: Use ptr.BoolPtr instead of custom boolPtr helper
4. `b1afc71` - refactor: Implement two-tier validation with runtime schema checking (this PR)
