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

[![Build Status](https://github.com/ainsleydev/webkit/actions/workflows/test.yaml/badge.svg)](https://github.com/ainsleydev/webkit/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ainsleydev/webkit)](https://goreportcard.com/report/github.com/ainsleydev/webkit)
[![Coverage](https://codecov.io/gh/ainsleydev/webkit/branch/main/graph/badge.svg)](https://codecov.io/gh/ainsleydev/webkit)
[![Maintainability](https://api.codeclimate.com/v1/badges/f5912a1dec11b8003850/maintainability)](https://codeclimate.com/github/ainsleydev/webkit/maintainability)
[![Latest Release](https://img.shields.io/github/v/release/ainsleydev/webkit)](https://github.com/ainsleydev/webkit/releases)
[![Twitter Handle](https://img.shields.io/twitter/follow/ainsleydev)](https://twitter.com/ainsleydev)

</div>

## WebKit

WebKit is a CLI tool that transforms a single `app.json` manifest into production-ready
infrastructure and CI/CD pipelines. It generates Terraform configurations, GitHub Actions workflows,
and project files and manages secrets in different enviornments.

**Key Features:**

- **Single source of truth**: Define apps, resources, and environments in one manifest
- **Infrastructure as code**: Automated Terraform generation and management
- **Secrets management**: Built-in SOPS/Age encryption with environment resolution
- **CI/CD automation**: GitHub Actions workflows for builds, deploys, and backups
- **Developer experience**: Idempotent updates, validation, and zero-config defaults

> **Note:** For user documentation and guides,
> visit [webkit.ainsley.dev](https://webkit.ainsley.dev)


## Installation

**Quick install:**
```bash
curl -sSL https://raw.githubusercontent.com/ainsleydev/webkit/main/bin/install.sh | sh
```

Or download binaries from the [latest release](https://github.com/ainsleydev/webkit/releases/latest).

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
│   │   ├── files/       # Project file generation
│   │   ├── infra/       # Terraform infrastructure commands
│   │   └── secrets/     # SOPS encryption/decryption
│   ├── infra/           # Terraform wrapper and state management
│   ├── manifest/        # File tracking and manifest source tagging
│   ├── scaffold/        # Template scaffolding system
│   ├── secrets/         # SOPS/Age integration
│   ├── templates/       # Embedded project templates
│   └── util/            # Shared utilities
├── platform/            # Terraform modules (infrastructure definitions)
│   ├── digitalocean/    # DigitalOcean provider modules
│   └── ...              # Additional providers
└── docs/                # Documentation source (whitepaper, specs)
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

**Note:** Make sure you have [act](https://github.com/nektos/act) installed:

```bash
brew install act
```

## Releasing

WebKit uses [GoReleaser](https://goreleaser.com/) for automated releases. The release process is
triggered by creating and pushing a git tag.

### Quick Release

Use the interactive tag tool:

```bash
pnpm tag
```

This will guide you through:

1. Choosing between creating or deleting a tag
2. Selecting the version bump type (patch, minor, or major)
3. Confirming the version
4. Creating and pushing the tag

### What Happens Next

When a tag is pushed:

1. GitHub Actions automatically triggers the GoReleaser workflow/
2. Binaries are built for each platform.
3. A GitHub release is created with the binaries attached
4. Release notes can be edited on the GitHub releases page

### Semantic Versioning

WebKit follows [Semantic Versioning](https://semver.org/).

- **Patch** (v1.0.1): Bug fixes and minor changes
- **Minor** (v1.1.0): New features, backwards compatible
- **Major** (v2.0.0): Breaking changes

---

## Copyright

You may not, except with our express written permission, distribute or commercially exploit the
content found within this repository or any written text within this repository. Nor may you transmit
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

