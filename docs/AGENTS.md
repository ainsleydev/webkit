## WebKit

WebKit is a Go-based CLI tool designed to streamline the lifecycle of web projects. It centralises
configuration in a single manifest file (`app.json`) and automatically generates surrounding
infrastructure, CI/CD pipelines, and environment scaffolding.

WebKit helps reduce repetitive setup work, improves consistency across deployments, and provides a
reliable foundation for infrastructure management.

**Key features**:

- Single source of truth via `app.json` manifest.
- Automatic generation of Docker configurations, GitHub workflows, and environment files.
- Infrastructure provisioning through Terraform.
- Secret management via SOPS.
- Template-based file generation with tracking.

## Build & Commands

Essential pnpm scripts for development are listed below.

- **Run CLI**: `pnpm webkit <command>`
- **Build**: `pnpm build`
- **Check all**: `pnpm check`
- **Format**: `pnpm format`
- **Format Go**: `pnpm format:go`
- **Lint**: `pnpm lint`
- **Lint Go**: `pnpm lint:go`
- **Test**: `pnpm test`

## Content

### Language and style

Write all content in British English. Use British spellings and punctuation throughout the codebase:

- "colour" not "color"
- "organised" not "organized"
- "centre" not "center"

### Heading style

Headings should use sentence case (only the first word capitalised):

- Correct: "Setting up webhooks"
- Incorrect: "Setting Up Webhooks"
- Incorrect: "Setting up Webhooks"

## Markdown

- Use `-` for lists, not `*`.
- End list points with a full stop.
