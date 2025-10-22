# WebKit Secrets Management Implementation Plan

## Overview

WebKit resolves environment variables from three sources: static values, Terraform resource outputs, and SOPS-encrypted secrets. The CLI decrypts and resolves all variables before generating Terraform configurations.

## Manifest Structure

### Environment Variable Schema

```json
{
  "apps": [{
    "name": "cms",
    "env": {
      "production": {
        "DATABASE_URL": {
          "source": "resource",
          "value": "db.connection_url"
        },
        "PAYLOAD_SECRET": {
          "source": "sops",
          "path": "secrets/production.yaml:PAYLOAD_SECRET"
        },
        "PUBLIC_API_URL": {
          "source": "value",
          "value": "https://cms.mysite.com"
        }
      },
      "development": {
        "DATABASE_URL": {
          "source": "resource",
          "value": "db.connection_url"
        },
        "PAYLOAD_SECRET": {
          "source": "sops",
          "path": "secrets/development.yaml:PAYLOAD_SECRET"
        },
        "PUBLIC_API_URL": {
          "source": "value",
          "value": "http://localhost:3000"
        }
      }
    }
  }]
}
```

### Source Types

- **value**: Static string value embedded in manifest
- **resource**: Reference to Terraform resource output (e.g., `db.connection_url`)
- **sops**: Encrypted secret stored in SOPS file with path format `file:key`

## SOPS File Structure

### Directory Layout

```
secrets/
├── production.yaml    # Encrypted
├── staging.yaml       # Encrypted
└── development.yaml   # Encrypted or plaintext
```

### Example SOPS File (Before Encryption)

```yaml
# secrets/production.yaml
PAYLOAD_SECRET: prod_secret_abc123
RESEND_API_KEY: re_prod_xyz789
STRIPE_SECRET: sk_live_abc123
```

### SOPS Configuration

`.sops.yaml` in project root:

```yaml
creation_rules:
  - path_regex: secrets/.*\.yaml$
    age: age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p
```

## Age Key Management

### One Key Per Organization

**Setup (One-time):**

```bash
# Generate age key
age-keygen -o ~/.config/webkit/age.key

# Get public key for .sops.yaml
age-keygen -y ~/.config/webkit/age.key
# Output: age1ql3z7hjy54pw3hyww5ayyfg7zqgvc7w3j2elw8zmrj2kg5sfn9aqmcac8p
```

**GitHub Organization Secret:**

- Name: `WEBKIT_SOPS_AGE_KEY`
- Value: `<contents of ~/.config/webkit/age.key>`

All projects in the organization use this single key.

## WebKit CLI Resolution Logic

### High-Level Process

```
1. Parse app.json
   ↓
2. For each environment (dev/staging/production):
   ↓
3. Iterate through environment variables
   ↓
4. Resolve based on source type:
   - value: use literal value
   - resource: keep as reference for Terraform
   - sops: decrypt now using SOPS Go package
   ↓
5. Output resolved variables
```

### Resolution Implementation (Go)

```go
func resolveEnvironmentVariables(env map[string]EnvConfig) ([]EnvVar, error) {
    resolved := []EnvVar{}
    
    for key, config := range env {
        switch config.Source {
        case "value":
            // Static value - use as-is
            resolved = append(resolved, EnvVar{
                Key:   key,
                Value: config.Value,
                Type:  "GENERAL",
            })
            
        case "resource":
            // Resource reference - Terraform will resolve
            resolved = append(resolved, EnvVar{
                Key:   key,
                Value: fmt.Sprintf("resource:%s", config.Value),
                Type:  "GENERAL",
            })
            
        case "sops":
            // Decrypt secret NOW
            decrypted, err := decryptSOPSSecret(config.Path)
            if err != nil {
                return nil, fmt.Errorf("failed to decrypt %s: %w", key, err)
            }
            resolved = append(resolved, EnvVar{
                Key:   key,
                Value: decrypted,
                Type:  "SECRET",
            })
        }
    }
    
    return resolved, nil
}
```

### SOPS Decryption (Using Go Package)

**Package:** `github.com/getsops/sops/v3`

```go
func decryptSOPSSecret(path string) (string, error) {
    // Parse path: "secrets/production.yaml:PAYLOAD_SECRET"
    parts := strings.Split(path, ":")
    if len(parts) != 2 {
        return "", fmt.Errorf("invalid sops path format: %s", path)
    }
    
    filePath := parts[0]  // "secrets/production.yaml"
    key := parts[1]       // "PAYLOAD_SECRET"
    
    // Decrypt file using SOPS Go package
    decryptedData, err := decrypt.File(filePath, "yaml")
    if err != nil {
        return "", err
    }
    
    // Parse YAML
    var data map[string]string
    if err := yaml.Unmarshal(decryptedData, &data); err != nil {
        return "", err
    }
    
    // Extract value
    value, exists := data[key]
    if !exists {
        return "", fmt.Errorf("key %s not found in %s", key, filePath)
    }
    
    return value, nil
}
```

**Why Go Package Instead of Exec:**

- No external binary dependency
- Cross-platform compatibility
- Better error handling
- Can decrypt in-memory
- Self-contained webkit binary

## Output Format

### Generated Terraform Variables

**Input (from app.json):**

```json
{
  "DATABASE_URL": {
    "source": "resource",
    "value": "db.connection_url"
  },
  "PAYLOAD_SECRET": {
    "source": "sops",
    "path": "secrets/production.yaml:PAYLOAD_SECRET"
  },
  "PUBLIC_API_URL": {
    "source": "value",
    "value": "https://cms.mysite.com"
  }
}
```

**Output (in generated tfvars):**

```json
{
  "env_vars": [
    {
      "key": "DATABASE_URL",
      "value": "resource:db.connection_url",
      "type": "GENERAL"
    },
    {
      "key": "PAYLOAD_SECRET",
      "value": "prod_secret_abc123",
      "type": "SECRET"
    },
    {
      "key": "PUBLIC_API_URL",
      "value": "https://cms.mysite.com",
      "type": "GENERAL"
    }
  ]
}
```

**Key Points:**

- SOPS secrets are decrypted to actual values
- Resource references remain as `resource:` prefix for Terraform to resolve
- Static values pass through unchanged

## Local Development

### Syncing Secrets to .env Files

**Command:** `webkit dev` or `webkit env sync`

**Process:**

1. Read `app.json`
2. Parse `env.development` section
3. Resolve all environment variables:
   - Decrypt SOPS secrets
   - Use local resource values (e.g., `postgres://localhost:5432`)
   - Use static values
4. Write to `services/cms/.env` using `joho/godotenv` package

**Libraries:**

- SOPS decryption: `github.com/getsops/sops/v3`
- Writing .env files: `github.com/joho/godotenv`

### Example .env Output

```bash
# Generated by webkit dev
# Do not commit this file

DATABASE_URL=postgresql://localhost:5432/myapp_dev
PAYLOAD_SECRET=dev_secret_xyz789
PUBLIC_API_URL=http://localhost:3000
NODE_ENV=development
```

**Why godotenv:**

- Handles escaping special characters
- Handles multiline values (e.g., private keys)
- Standard .env format
- Can preserve comments

## CI/CD Environment

### Age Key Access

**GitHub Actions environment variable:**

```yaml
env:
  SOPS_AGE_KEY: ${{ secrets.WEBKIT_SOPS_AGE_KEY }}
```

**WebKit CLI behavior:**

1. Checks for `SOPS_AGE_KEY` environment variable
2. If present, uses it for decryption
3. If not present, looks for `~/.config/webkit/age.key`
4. Fails if neither exists

## Terraform Resolution

**Resource references are resolved by Terraform modules:**

```hcl
# In modules/apps/main.tf
envs = [
  for env in var.env_vars : {
    key = env.key
    value = (
      startswith(env.value, "resource:")
      ? var.resource_outputs[split(".", trimprefix(env.value, "resource:"))[0]][split(".", trimprefix(env.value, "resource:"))[1]]
      : env.value
    )
    type = env.type
  }
]
```

**Flow:**

1. CLI embeds decrypted SOPS values in tfvars
2. CLI keeps resource references as strings with `resource:` prefix
3. Terraform module strips `resource:` prefix and looks up actual value from `resource_outputs` variable

## Security Considerations

### What Gets Encrypted

- Application secrets (API keys, passwords)
- Third-party service tokens
- Encryption keys

### What Doesn't Get Encrypted

- Infrastructure provider credentials (stored in GitHub Secrets)
- SOPS age key itself (stored in GitHub Secrets)
- Public URLs and non-sensitive configuration

### Best Practices

1. Never commit decrypted SOPS files or `.age.key`
2. Rotate keys when team members leave
3. Use separate files per environment (production, staging, development)
4. Audit access via git history on encrypted files
5. Use age over PGP - simpler and more secure

## Implementation Checklist

- [ ] Update `app.json` schema to support new env var structure
- [ ] Implement source type parsing in Go structs
- [ ] Integrate `github.com/getsops/sops/v3` package
- [ ] Build SOPS path parsing (`file:key` format)
- [ ] Implement decryption logic using SOPS Go package
- [ ] Build environment variable resolver for each source type
- [ ] Create Terraform tfvars generator with resolved values
- [ ] Integrate `github.com/joho/godotenv` for local .env generation
- [ ] Implement `webkit dev` command for local environment setup
- [ ] Add age key detection (env var vs file)
- [ ] Create error handling for missing secrets
- [ ] Add validation for SOPS path format
- [ ] Document secret file structure and naming conventions