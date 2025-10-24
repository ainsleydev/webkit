## WebKit specific guidelines

This section contains WebKit-specific agent guidelines that are injected into the base AGENTS.md template.

### Project structure

- `cmd/` - Standalone Go commands (genversion, tag, docs).
- `internal/` - Internal packages for WebKit core functionality.
- `pkg/` - Public packages that can be imported by other projects.
- `platform/` - Platform-specific implementations.
- `docs/` - VitePress documentation site.
- `ai/` - AI-related files and prompts.

### Building and running

Use the Makefile for common tasks:

- `make build` - Build the WebKit CLI.
- `make test` - Run all tests.
- `make lint` - Run linters.

### Code generation

WebKit uses several code generation tools:

- `cmd/genversion` - Generates version information.
- `cmd/docs` - Generates AGENTS.md from templates.
- Template-based scaffolding via `internal/scaffold`.

### Working with manifests

The manifest tracker (`internal/manifest`) tracks all generated files. Always use `scaffold.WithTracking()` when generating files to ensure they're properly tracked.
