<p align="center">
  <a href="https://webkit.ainsley.dev">
    <img src="./resources/symbol.png" height="96">
    <h3 align="center">WebKit</h3>
  </a>
</p>

<p align="center">
  Infrastructure-as-code framework for full-stack web applications
</p>

<p align="center">
  <a href="https://webkit.ainsley.dev"><strong>Documentation</strong></a> ·
  <a href="https://webkit.ainsley.dev/getting-started"><strong>Getting Started</strong></a> ·
  <a href="https://webkit.ainsley.dev/examples"><strong>Examples</strong></a> ·
</p>

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/ainsleydev/webkit)](https://goreportcard.com/report/github.com/ainsleydev/webkit)
[![Release](https://img.shields.io/github/v/release/ainsleydev/webkit?color=brightgreen&label=Release)](https://github.com/ainsleydev/webkit/releases)
[![Maintainability](https://api.codeclimate.com/v1/badges/f5912a1dec11b8003850/maintainability)](https://codeclimate.com/github/ainsleydev/webkit/maintainability)
[![Coverage](https://codecov.io/gh/ainsleydev/webkit/branch/main/graph/badge.svg)](https://codecov.io/gh/ainsleydev/webkit)
![Made with Go](https://img.shields.io/badge/Made%20with-Go-00ADD8.svg?logo=go)
[![Go Reference](https://pkg.go.dev/badge/github.com/ainsleydev/webkit.svg)](https://pkg.go.dev/github.com/ainsleydev/webkit)
[![Twitter Handle](https://img.shields.io/twitter/follow/ainsleydev)](https://twitter.com/ainsleydev)

</div>

## WebKit

WebKit is a CLI tool that transforms a single `app.json` manifest into production-ready
infrastructure and CI/CD pipelines. It generates Terraform configurations, GitHub Actions workflows,
and project files - all without cluttering your project repository with infrastructure code.

**Key Features:**

- **Single source of truth**: Define apps, resources, and environments in `app.json`
- **Clean repositories**: No `infra/` folder needed - workflows contain everything
- **Infrastructure as code**: Automated Terraform generation with centralized modules
- **Secrets management**: Built-in SOPS/Age encryption with environment-specific decryption
- **CI/CD automation**: GitHub Actions workflows for plans, deploys, and drift detection
- **Developer experience**: Idempotent updates, local testing, and zero-config defaults

> **Note:** For user documentation and guides,
> visit [webkit.ainsley.dev](https://webkit.ainsley.dev)

## Installation

**Quick install:**

```bash
curl -sSL https://raw.githubusercontent.com/ainsleydev/webkit/main/bin/install.sh | sh
```

Or download binaries from
the [latest release](https://github.com/ainsleydev/webkit/releases/latest).

**Verify installation:**

```bash
webkit version
```

## Development Setup

Run the following to get setup with `webkit`.

```bash
make setup
```

### Prerequisites

- **Go** 1.23 or higher
- **pnpm** (for task runners and local workflow testing)
- **act** (for testing GitHub Actions locally)
- **age** (for secrets encryption)

### Project Structure

```
webkit/
├── cmd/webkit/          # CLI entry point
├── internal/
│   ├── appdef/          # App manifest parsing and validation
│   ├── cmd/             # CLI command implementations
│   │   ├── cicd/        # GitHub Actions workflow generation
│   │   ├── env/         # Environment variable management
│   │   ├── docs/        # Documentation generation
│   │   ├── files/       # Project file generation
│   │   ├── infra/       # Terraform infrastructure commands
│   │   └── secrets/     # SOPS encryption/decryption
│   ├── infra/           # Terraform wrapper and state management
│   ├── manifest/        # File tracking and manifest source tagging
│   ├── scaffold/        # Template scaffolding system
│   ├── secrets/         # SOPS/Age integration
│   ├── templates/       # Embedded project templates
│   └── util/            # Shared utilities
└── platform/            # Terraform modules (separate infra repository)
    ├── providers/       # Provider-specific modules (DO, B2, etc.)
    └── modules/         # Orchestration modules (apps, resources)
```

### Architecture Overview

**Key Components:**

- **appdef**: Defines the structure of `app.json` and handles unmarshaling with validation
- **scaffold**: Template rendering engine with file tracking and idempotent updates
- **manifest**: Tracks which files are generated and their sources (app/resource/project)
- **secrets**: SOPS integration for encrypting/decrypting environment variables
- **infra**: Terraform wrapper that manages state and tfvars generation

## Local Workflow Testing

You can simulate GitHub Actions workflows locally using [act](https://github.com/nektos/act). `act`
runs from your local computer, whatever files are currently in your working directory, including
any uncommitted changes.

```bash
# Test lint workflow
pnpm act:lint

# Test test workflow
pnpm act:test

# Dry-run release workflow (shows what would run without executing)
pnpm act:release
```

## Releasing

*For maintainers only*

WebKit uses [GoReleaser](https://goreleaser.com/) for automated releases. The release process is
triggered by creating and pushing a git tag.

### Quick Release

Use the interactive tag tool:

```bash
pnpm tag
```

This will guide you through:

1. Choosing between creating or deleting a tag.
2. Selecting the version bump type (`patch`, `minor`, or `major`).
3. Confirming the version.
4. Creating and pushing the tag.

When a tag is pushed, the version will be injected then GitHub Actions automatically triggers the
GoReleaser workflow, builds binaries for each platform and creates a GitHub release with the
binaries attached.

### Semantic Versioning

WebKit follows [Semantic Versioning](https://semver.org/).

- **Patch** (v1.0.1): Bug fixes and minor changes
- **Minor** (v1.1.0): New features, backwards compatible
- **Major** (v2.0.0): Breaking changes

## Copyright

You may not, except with our express written permission, distribute or commercially exploit the
content found within this repository or any written text within this repository. Nor may you
transmit
it or store it in any other website or other form of electronic retrieval system.

Any redistribution or reproduction of part or all of the contents in any form is prohibited other
than the following:

- You may print or download to a local hard disk extracts for your personal and non-commercial use
  only.
- You may copy the content to individual third parties for their personal use, but only if you
  acknowledge the website
  as the source of the material,

## License

Code Copyright 2023 ainsley.dev. Code released under the [BSD-3 Clause](LICENSE).

