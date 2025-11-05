# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this
repository.

## Overview

WebKit is a CLI tool written in Go that transforms a single `app.json` manifest into
production-ready infrastructure and CI/CD pipelines. It generates Terraform configurations, GitHub
Actions workflows, and project files without cluttering project repositories with infrastructure
code.

## Core Commands

### Development

```bash
# Initial setup - installs all dev dependencies
make setup

# Run webkit locally in development mode
pnpm webkit <args>

# Build and install webkit binary
pnpm build
```

### Testing

```bash
# Run all tests (Go + JS)
pnpm test

# Run Go tests only
pnpm test:go

# Run JS tests only (npm packages in packages/)
pnpm test:js

# Run single test
go test -run TestName ./path/to/package
```

### Linting and Formatting

```bash
# Format and lint everything
pnpm check

# Format Go code
pnpm format:go

# Format JS/TS code
pnpm format:js

# Lint Go code (uses golangci-lint)
pnpm lint:go

# Lint JS/TS code (uses biome)
pnpm lint:js
```

### Local Workflow Testing

```bash
# Test GitHub Actions workflows locally using act
pnpm act:lint
pnpm act:test
pnpm act:release  # dry-run only
```

### Generate Commands

```bash
# Generate documentation
pnpm generate:docs

# Generate agent prompts
pnpm generate:agents
```

## Architecture

### High-Level Flow

1. **app.json** → User defines apps, resources, and environments
2. **appdef** → Parses and validates the manifest
3. **scaffold** → Generates files from templates with tracking
4. **manifest** → Tracks generated files with hashes for drift detection
5. **infra** → Wraps Terraform to provision infrastructure
6. **cicd** → Generates GitHub Actions workflows
7. **secrets** → Handles SOPS/Age encryption for environment variables

### Key Packages

#### internal/appdef

Defines the structure of `app.json` and handles parsing/validation. Core types:

- `App`: Represents an application (type: payload, svelte-kit, golang)
- `Resource`: Represents infrastructure resources (postgres, s3, etc.)
- `Environment`: Environment-specific variables (dev, staging, production)
- `Project`: Top-level project metadata

Each app has:

- Build configuration (Dockerfile, port)
- Infrastructure configuration (provider, type, config)
- Environment variables (can reference resources, SOPS secrets, or plain values)
- Commands (build, lint, test, format) with sensible defaults per app type

#### internal/scaffold

Template rendering engine that:

- Generates files from Go templates embedded in `internal/templates/`
- Tracks all generated files in `.webkit/manifest.json`
- Supports two modes:
	- `ModeGenerate`: Always overwrites files
	- `ModeScaffold`: Only creates files if they don't exist
- Adds WebKit notices to generated files
- Tracks file hashes for drift detection

#### internal/manifest

Manages `.webkit/manifest.json` which tracks:

- All generated files and their paths
- What generated each file (e.g., "cicd.BackupWorkflow")
- What caused generation (e.g., "resource:postgres-prod")
- Content hashes (SHA256) for drift detection
- Whether files were scaffolded (user-editable) or generated (overwritable)

The manifest enables:

- Cleanup of orphaned files when manifest changes
- Drift detection via `webkit drift` command
- Idempotent updates (only regenerate changed files)

#### internal/infra

Terraform wrapper that:

- Generates `terraform.tfvars.json` from `app.json`
- Manages Terraform state via DigitalOcean Spaces backend
- Handles `plan`, `apply`, `destroy` commands
- Imports existing resources into Terraform state
- Converts app.json structure to Terraform variables

Key insight: Terraform modules live in `platform/terraform/` and are versioned separately. The CLI
references them via GitHub releases.

#### internal/secrets

SOPS/Age integration for secrets management:

- Encrypts/decrypts `.env.{environment}.enc` files
- Resolves environment variables at build time
- Supports three variable sources:
	- `value`: Plain text values
	- `sops`: Encrypted secrets from SOPS files
	- `resource`: Values from Terraform outputs (e.g., database URLs)

#### internal/cmd

CLI command implementations using urfave/cli/v3:

- `update`: Main command - regenerates all files from app.json
- `secrets`: Manage SOPS encryption
- `env`: Generate .env files for environments
- `infra`: Terraform operations (plan, apply, destroy, import)
- `cicd`: Generate GitHub Actions workflows
- `drift`: Detect files that have changed since generation
- `docs`: Generate documentation

### Terraform Module Structure

The `platform/terraform/` directory contains:

- `base/`: Base Terraform configuration
- `modules/`: Orchestration modules
	- `apps/`: Module for deploying applications
	- `resources/`: Module for provisioning resources
- `providers/`: Provider-specific implementations
	- `digital_ocean/`: DigitalOcean resources (app, droplet, postgres, bucket, domain_record)
	- `b2/`: Backblaze B2 storage

The CLI generates `terraform.tfvars.json` from `app.json`, which is consumed by these modules.

### Template System

Templates live in `internal/templates/` with the following structure:

- `.github/`: GitHub Actions workflow templates
- `docs/`: Documentation templates
- `terraform/`: Terraform configuration templates
- Various project file templates (.gitignore, package.json, etc.)

Templates use Go's `text/template` with custom functions defined in `internal/templates/funcs.go`.
Common template helpers:

- `toJSON`: Convert to JSON
- `kebabCase`: Convert string to kebab-case
- `snakeCase`: Convert string to snake_case
- `camelCase`: Convert string to camelCase

## Important Patterns

### Idempotent Updates

The `webkit update` command is designed to be run repeatedly. It:

1. Reads the old manifest from `.webkit/manifest.json`
2. Generates all files based on current `app.json`
3. Writes new manifest
4. Cleans up orphaned files (files in old manifest but not new manifest)

### File Tracking

Every generated file is tracked with:

- Generator name (for debugging)
- Source (what in app.json caused this)
- Content hash (for drift detection)
- Scaffold mode flag (determines if file should be overwritten)

Use `WithTracking()` option when calling scaffold methods:

```go
gen.Template(path, tpl, data,
    scaffold.WithTracking("cicd.DeployWorkflow", "app:web"),
    scaffold.WithMode(scaffold.ModeGenerate),
)
```

### Environment Variable Resolution

Variables can have three sources:

1. `value`: Static string (e.g., "http://localhost:3000")
2. `sops`: Encrypted secret key (e.g., "PAYLOAD_SECRET")
3. `resource`: Terraform output path (e.g., "db.connection_url")

The resolution happens at different times:

- `value`: Resolved immediately
- `sops`: Resolved during `webkit env generate`
- `resource`: Resolved during Terraform apply and exposed as outputs

## Testing

Tests use standard Go testing with:

- `testify/assert` for assertions
- `spf13/afero` for filesystem mocking
- `go.uber.org/mock` for interface mocking (generated via `mockgen`)

When testing scaffold/manifest functionality, use `afero.NewMemMapFs()` to avoid touching real
filesystem.

## Release Process

1. Create and push a git tag using `pnpm tag` (interactive)
2. GitHub Actions triggers GoReleaser workflow
3. GoReleaser builds binaries for all platforms
4. Binaries are attached to GitHub release

For npm packages in `packages/`:

1. Use Changesets: `pnpm changeset`
2. Create PR with changeset
3. When merged, Changesets action creates release PR
4. Merge release PR to publish packages

## Key Files

- `app.json`: User's manifest file (not tracked in this repo, but lives in user projects)
- `.webkit/manifest.json`: Tracks all generated files
- `.env.{env}.enc`: SOPS-encrypted environment variables
- `go.mod`: Go 1.25.3, uses urfave/cli/v3, sprig for templates, terraform-exec
- `package.json`: pnpm workspace with scripts for dev workflow
- `Makefile`: Setup script for installing dependencies

## Working with Terraform

When modifying Terraform generation:

1. Update structs in `internal/infra/tf_vars.go`
2. Update corresponding Terraform variables in `platform/terraform/`
3. Test with playground app: `cd internal/playground && webkit update`
4. Verify generated `terraform.tfvars.json`

## Common Patterns

### Adding a New Command

1. Create package in `internal/cmd/yourcommand/`
2. Define `var Command = &cli.Command{...}`
3. Add to commands list in `internal/cmd/cli.go`
4. Add tests in `yourcommand_test.go`

### Adding a New Template

1. Add template file to `internal/templates/`
2. Update `internal/templates/embed.go` if needed
3. Create generator function in relevant cmd package
4. Use scaffold.WithTracking() to track generated file

### Adding a New Resource Type

1. Define type in `internal/appdef/resource.go`
2. Add provider implementation in `platform/terraform/providers/{provider}/{type}/`
3. Update `internal/infra/tf_vars.go` to generate correct tfvars
4. Add to resources module in `platform/terraform/modules/resources/`
