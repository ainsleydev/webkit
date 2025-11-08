# Validation

WebKit provides comprehensive validation for `app.json` configuration files to catch errors early and ensure your project is properly configured.

## Overview

Validation is automatically integrated into webkit and runs whenever you load an `app.json` file. It validates:

- **Required fields**: Ensures all critical configuration is present
- **Domain formats**: Prevents common mistakes like including protocol prefixes
- **File paths**: Verifies that app directories actually exist
- **Resource references**: Checks that environment variables reference valid resources
- **Infrastructure requirements**: Ensures terraform-managed VMs have required domains

## Commands

### `webkit validate`

Explicitly validates your `app.json` configuration and displays detailed error messages.

```bash
webkit validate
```

**Success output:**
```
ℹ Validating app.json...
✓ Validation passed! No errors found.
```

**Error output:**
```
ℹ Validating app.json...
✗ Validation failed with 2 error(s):

  1. app "cms": domain "https://example.com" should not contain protocol prefix (e.g., 'https://')
  2. app "api": path "api" does not exist
```

### `webkit schema`

Generates a JSON schema file for IDE autocomplete and validation support.

```bash
webkit schema --output schema.json
```

The generated `schema.json` file can be referenced in your `app.json`:

```json
{
  "$schema": "./schema.json",
  "webkit_version": "v0.0.40",
  "project": { ... }
}
```

## Automatic Validation

Validation runs automatically during:

- `webkit update` - Validates before updating any generated files
- `appdef.Read()` - Any command that loads the app definition
- CI/CD pipelines - Catches configuration errors before deployment

## Validation Rules

### Required Fields

**Project:**
- `name` - Unique project identifier
- `title` - Human-readable project name
- `description` - Project description
- `repo.owner` - GitHub repository owner
- `repo.name` - GitHub repository name

**App:**
- `name` - Unique app identifier
- `title` - Human-readable app name
- `type` - Application type (payload, svelte-kit, golang)
- `path` - Relative path to app source code
- `infra.provider` - Cloud provider (digitalocean, backblaze)
- `infra.type` - Infrastructure type (vm, app, container, function)

**Resource:**
- `name` - Unique resource identifier
- `type` - Resource type (postgres, s3)
- `provider` - Cloud provider (digitalocean, backblaze)

### Domain Validation

Domains must **not** contain protocol prefixes:

❌ **Invalid:**
```json
{
  "domains": [
    { "name": "https://example.com" },
    { "name": "http://api.example.com" }
  ]
}
```

✅ **Valid:**
```json
{
  "domains": [
    { "name": "example.com" },
    { "name": "api.example.com" }
  ]
}
```

### Path Validation

App paths must exist on the filesystem:

```json
{
  "apps": [
    {
      "name": "cms",
      "path": "cms",  // This directory must exist
      ...
    }
  ]
}
```

If the path doesn't exist, you'll get an error:
```
app "cms": path "cms" does not exist
```

### Terraform-Managed VM Validation

Apps with `infra.type` set to `"vm"` or `"app"` that are terraform-managed (default) **must** have at least one domain configured:

❌ **Invalid:**
```json
{
  "name": "api",
  "infra": {
    "type": "vm",
    "provider": "digitalocean"
  },
  "domains": []  // Error: VM apps must have domains
}
```

✅ **Valid:**
```json
{
  "name": "api",
  "infra": {
    "type": "vm",
    "provider": "digitalocean"
  },
  "domains": [
    { "name": "api.example.com", "type": "primary" }
  ]
}
```

### Environment Variable Validation

Environment variables with `source: "resource"` must reference valid resources and outputs:

❌ **Invalid:**
```json
{
  "resources": [
    { "name": "db", "type": "postgres", "provider": "digitalocean" }
  ],
  "apps": [
    {
      "name": "api",
      "env": {
        "production": {
          "DB_URL": {
            "source": "resource",
            "value": "nonexistent.connection_url"  // Error: resource doesn't exist
          }
        }
      }
    }
  ]
}
```

✅ **Valid:**
```json
{
  "resources": [
    { "name": "db", "type": "postgres", "provider": "digitalocean" }
  ],
  "apps": [
    {
      "name": "api",
      "env": {
        "production": {
          "DB_URL": {
            "source": "resource",
            "value": "db.connection_url"  // Valid: references existing resource
          }
        }
      }
    }
  ]
}
```

**Valid outputs per resource type:**

- **postgres**: `id`, `connection_url`, `host`, `port`, `database`, `user`, `password`
- **s3**: `id`, `bucket_name`, `bucket_url`, `region`

## IDE Support

For the best development experience, add the `$schema` field to your `app.json`:

```json
{
  "$schema": "./schema.json",
  "webkit_version": "v0.0.40",
  ...
}
```

This enables:
- **Autocomplete** - Suggestions for field names and values
- **Inline validation** - Red squiggles for configuration errors
- **Hover documentation** - Descriptions for all fields
- **Type checking** - Ensures correct data types

### VS Code

VS Code automatically recognises the `$schema` field and provides full IntelliSense support.

### JetBrains IDEs

JetBrains IDEs (WebStorm, GoLand, IntelliJ IDEA) also support JSON schema validation natively.

## Generating Schema

To regenerate the schema after webkit updates:

```bash
# Manual generation
webkit schema --output schema.json

# Or use go generate
go generate ./internal/appdef/...
```

The schema is automatically generated during webkit development and should be committed to your repository.

## Error Collection

Validation collects **all** errors before reporting them, rather than stopping at the first error. This allows you to fix multiple issues at once:

```
✗ Validation failed with 4 error(s):

  1. app "cms": domain "https://cms.example.com" should not contain protocol prefix
  2. app "api": path "api" does not exist
  3. app "web": terraform-managed VM/app must have at least one domain configured
  4. shared: env var "S3_BUCKET" in production references non-existent resource "storage"
```

## Implementation Details

The validation system is built using:
- **JSON Schema** - For structural validation and IDE support
- **Custom Validators** - For business logic and cross-field validation
- **swaggest/jsonschema-go** - For schema generation from Go structs

Validation is implemented in:
- `internal/appdef/validate.go` - Core validation logic
- `internal/appdef/schema.go` - Schema generation
- `internal/cmd/validate.go` - CLI command
- `internal/cmd/schema.go` - Schema generation command

## Testing

Comprehensive validation tests are located in `internal/appdef/validate_test.go`, covering:
- All validation rules
- Error message formatting
- Edge cases and error conditions
- Integration with the broader webkit system
