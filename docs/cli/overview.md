# Command reference

WebKit provides a comprehensive CLI for managing your project lifecycle. This page documents all available commands.

## Global options

These options are available for all commands:

| Flag | Description |
|------|-------------|
| `--help`, `-h` | Show help for any command |
| `--version`, `-v` | Show WebKit version |

## Commands

### webkit update

Regenerate all project files from your `app.json` manifest.

```bash
webkit update
```

This is the primary command you'll use. It:
- Reads your `app.json` manifest
- Generates GitHub Actions workflows
- Creates Docker configuration
- Updates project files (package.json, turbo.json, etc.)
- Tracks all generated files in `.webkit/manifest.json`
- Cleans up orphaned files

### webkit validate

Validate your `app.json` manifest without generating files.

```bash
webkit validate
```

Checks for:
- Required fields
- Valid values and types
- Resource references
- Domain configuration

### webkit drift

Detect manual modifications to generated files.

```bash
webkit drift
```

Compares current file contents against stored hashes to identify changes made outside of WebKit.

### webkit scaffold

Generate individual components without running a full update.

```bash
webkit scaffold <component>
```

Available components vary based on your configuration.

### webkit version

Display the installed WebKit version.

```bash
webkit version
```

## Subcommand groups

### webkit infra

Infrastructure management commands using Terraform.

| Command | Description |
|---------|-------------|
| `webkit infra plan` | Preview infrastructure changes |
| `webkit infra apply` | Apply infrastructure changes |
| `webkit infra destroy` | Destroy all infrastructure |
| `webkit infra output` | Display Terraform outputs |
| `webkit infra import` | Import existing resources |
| `webkit infra exec -- <cmd>` | Run arbitrary Terraform commands |

See [Infrastructure overview](/infrastructure/overview) for detailed documentation.

### webkit secrets

SOPS-encrypted secrets management.

| Command | Description |
|---------|-------------|
| `webkit secrets scaffold` | Create secret file templates |
| `webkit secrets sync` | Sync secrets with SOPS files |
| `webkit secrets encrypt` | Encrypt a secrets file |
| `webkit secrets decrypt` | Decrypt a secrets file |
| `webkit secrets get <key>` | Get a specific secret value |
| `webkit secrets validate` | Validate secrets configuration |

### webkit env

Environment variable management.

| Command | Description |
|---------|-------------|
| `webkit env scaffold` | Create .env.example files |
| `webkit env sync` | Sync environment variables |
| `webkit env generate` | Generate .env files for an environment |

### webkit cicd

CI/CD workflow generation.

| Command | Description |
|---------|-------------|
| `webkit cicd actions` | Copy reusable GitHub Actions |
| `webkit cicd backup` | Generate backup workflow |
| `webkit cicd pr` | Generate PR workflow |
| `webkit cicd release` | Generate release workflow |

### webkit docs

Documentation generation.

| Command | Description |
|---------|-------------|
| `webkit docs agents` | Generate AGENTS.md |
| `webkit docs readme` | Generate README.md |

### webkit payload

Payload CMS management.

| Command | Description |
|---------|-------------|
| `webkit payload bump` | Update Payload CMS dependencies |
| `webkit payload bump --dry-run` | Preview updates without applying |

## Examples

### Initial project setup

```bash
# Create your manifest
vim app.json

# Generate all project files
webkit update

# Validate configuration
webkit validate
```

### Deploy infrastructure

```bash
# Preview changes
webkit infra plan

# Apply changes
webkit infra apply

# View outputs
webkit infra output
```

### Manage secrets

```bash
# Create secret templates
webkit secrets scaffold

# Encrypt secrets for production
webkit secrets encrypt production

# Get a specific secret
webkit secrets get PAYLOAD_SECRET --env production
```

### Update environment files

```bash
# Generate .env for local development
webkit env generate dev

# Sync all environments
webkit env sync
```

## Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | General error |
| `2` | Configuration error |
| `3` | Validation error |

## Further reading

- [Quick start guide](/getting-started/quick-start)
- [Infrastructure overview](/infrastructure/overview)
- [Manifest reference](/manifest/overview)
