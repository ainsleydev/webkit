# Validation

WebKit provides comprehensive validation for `app.json` configuration files to catch errors early and ensure your project is properly configured.

## Overview

WebKit uses a **two-tier validation approach** that combines JSON schema validation with custom business logic validation:

### Tier 1: Schema Validation
Validates structural and type constraints using JSON schema:

- **Required fields**: Ensures all critical configuration is present
- **Data types**: Validates strings, numbers, booleans, arrays, and objects
- **Pattern validation**: Enforces naming conventions (e.g., lowercase-hyphenated names)
- **Enum constraints**: Restricts values to predefined options (e.g., app types, providers)
- **Numeric ranges**: Validates port numbers (1-65535)
- **String length**: Enforces maximum lengths for descriptions (200 characters)
- **Format validation**: Validates URLs and other formatted strings

### Tier 2: Custom Validation
Validates business logic and cross-field constraints:

- **Domain formats**: Prevents common mistakes like including protocol prefixes
- **File paths**: Verifies that app directories actually exist on the filesystem
- **Resource references**: Checks that environment variables reference valid resources and outputs
- **Infrastructure requirements**: Ensures terraform-managed VMs have required domains configured

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
# Generate to .webkit/schema.json (default)
webkit schema

# Generate to custom location
webkit schema --output custom-path/schema.json

# Output to stdout
webkit schema --stdout
```

The generated schema file can be referenced in your `app.json`:

```json
{
  "$schema": ".webkit/schema.json",
  "webkit_version": "v0.0.40",
  "project": { ... }
}
```

**Note:** The schema is also available on GitHub at the project URL and can be referenced directly for IDE support without local generation.

## Automatic Validation

Validation runs automatically during:

- `webkit update` - Validates before updating any generated files
- `appdef.Read()` - Any command that loads the app definition
- CI/CD pipelines - Catches configuration errors before deployment

## Validation Rules

### Required Fields (Schema Validation)

**Project:**
- `name` - Unique project identifier (lowercase-hyphenated pattern)
- `title` - Human-readable project name
- `description` - Project description (max 200 characters)
- `repo.owner` - GitHub repository owner
- `repo.name` - GitHub repository name

**App:**
- `name` - Unique app identifier (lowercase-hyphenated pattern)
- `title` - Human-readable app name
- `type` - Application type (must be: `svelte-kit`, `golang`, or `payload`)
- `path` - Relative path to app source code
- `infra.provider` - Cloud provider (must be: `digitalocean` or `backblaze`)
- `infra.type` - Infrastructure type (vm, app, container, function)

**Resource:**
- `name` - Unique resource identifier (lowercase-hyphenated pattern)
- `type` - Resource type (must be: `postgres` or `s3`)
- `provider` - Cloud provider (must be: `digitalocean` or `backblaze`)

**Environment Variables:**
- `source` - Source type (must be: `value`, `resource`, or `sops`)

### Pattern Validation (Schema Validation)

Names must follow the lowercase-hyphenated pattern:

❌ **Invalid:**
```json
{
  "name": "MyApp",        // Uppercase
  "name": "my_app",       // Underscores
  "name": "my.app",       // Dots
  "name": "123-app"       // Starts with number
}
```

✅ **Valid:**
```json
{
  "name": "my-app",
  "name": "myapp",
  "name": "my-app-123"
}
```

### Enum Validation (Schema Validation)

Certain fields only accept predefined values:

**App Types:**
- `svelte-kit` - SvelteKit application
- `golang` - Go application
- `payload` - Payload CMS application

**Resource Types:**
- `postgres` - PostgreSQL database
- `s3` - S3-compatible object storage

**Providers:**
- `digitalocean` - DigitalOcean cloud provider
- `backblaze` - Backblaze cloud provider

**Environment Sources:**
- `value` - Static string value
- `resource` - Terraform resource reference
- `sops` - Encrypted secret from SOPS

**Domain Types:**
- `primary` - Primary domain for the app
- `alias` - Alias/redirect domain
- `unmanaged` - Domain not managed by webkit

### Numeric Validation (Schema Validation)

**Port numbers** must be between 1 and 65535:

❌ **Invalid:**
```json
{
  "build": {
    "port": 0        // Too low
  }
}
```

```json
{
  "build": {
    "port": 70000    // Too high
  }
}
```

✅ **Valid:**
```json
{
  "build": {
    "port": 3000
  }
}
```

### String Length Validation (Schema Validation)

**Descriptions** have a maximum length of 200 characters:

❌ **Invalid:**
```json
{
  "description": "This is a very long description that exceeds the maximum allowed length of 200 characters. It contains way too much information and should be condensed to be more concise and fit within the character limit imposed by the schema validation rules."
}
```

✅ **Valid:**
```json
{
  "description": "A concise description of the project or app."
}
```

### Format Validation (Schema Validation)

**Webhook URLs** must be valid URIs:

❌ **Invalid:**
```json
{
  "notifications": {
    "webhook_url": "not-a-valid-url"
  }
}
```

✅ **Valid:**
```json
{
  "notifications": {
    "webhook_url": "https://hooks.slack.com/services/XXX/YYY/ZZZ"
  }
}
```

### Domain Validation (Custom Validation)

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

### Path Validation (Custom Validation)

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

### Terraform-Managed VM Validation (Custom Validation)

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

### Environment Variable Validation (Custom Validation)

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

The validation system uses a **two-tier architecture** combining JSON schema validation with custom business logic:

### Tier 1: Schema Validation
- **swaggest/jsonschema-go** - Generates JSON schema from Go struct tags
- **kaptinlin/jsonschema** - Validates JSON data against the schema at runtime
- **Struct tags** - Define validation constraints (`required`, `pattern`, `enum`, `minimum`, `maximum`, `maxLength`, `format`)

### Tier 2: Custom Validation
- **Filesystem checks** - Verifies app paths exist
- **Cross-field validation** - Ensures resource references are valid
- **Business logic** - Enforces webkit-specific rules (e.g., VMs must have domains)

### Validation Flow

1. **Schema Validation** (`ValidateAgainstSchema()`)
   - Runs first to catch structural and type errors
   - Validates against embedded `schema.json`
   - Fast and comprehensive structural checks

2. **Custom Validation** (`Validate()`)
   - Runs after schema validation passes
   - Performs filesystem and cross-reference checks
   - Enforces business logic constraints

### Implementation Files

**Schema Validation:**
- `internal/appdef/schema.go` - Schema generation from Go structs
- `internal/appdef/validate_schema.go` - Runtime schema validation
- `internal/appdef/schema.json` - Embedded JSON schema

**Custom Validation:**
- `internal/appdef/validate.go` - Business logic and cross-field validation

**Commands:**
- `internal/cmd/validate.go` - Explicit validation command
- `internal/cmd/schema.go` - Schema generation command

**Integration:**
- `internal/appdef/definition.go` - Integrates both validation tiers in `Read()`

## Testing

Comprehensive validation tests cover both tiers:

**Schema Validation Tests** (`internal/appdef/validate_schema_test.go`):
- Required field validation
- Pattern matching for names
- Enum constraint enforcement
- Numeric range validation
- String length limits
- Format validation (URLs)
- Integration with struct marshaling

**Custom Validation Tests** (`internal/appdef/validate_test.go`):
- Domain format validation
- File path existence checks
- Resource reference validation
- Terraform VM domain requirements
- Error message formatting
- Edge cases and error conditions
