## What is WebKit?

WebKit is a CLI tool that transforms a single `app.json` manifest into production-ready
infrastructure and CI/CD pipelines. It generates Terraform configurations, GitHub Actions workflows,
and standardized project files—letting you focus on building features instead of managing DevOps
boilerplate.





---



### Install Dependencies

```bash
# Clone the repository
git clone https://github.com/ainsleydev/webkit.git
cd webkit

# Install Go dependencies
go mod download

# Install act for local workflow testing
brew install act

# Install age for secrets management
brew install age
```

### Build from Source

```bash
# Build the CLI binary
go build -o webkit ./cmd/webkit

# Install globally
go install ./cmd/webkit
```

### Running Tests

```bash

```



### Architecture Overview

**Core Flow:**

1. **Parse**: Read and validate `app.json` against JSON schema
2. **Resolve**: Decrypt secrets and merge environment variables
3. **Transform**: Convert manifest to Terraform variables and workflow configs
4. **Generate**: Write files using Go templates with tracking metadata
5. **Deploy**: Execute Terraform through GitHub Actions or CLI





---

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

1. GitHub Actions automatically triggers the GoReleaser workflow
2. Binaries are built for each platform
3. A GitHub release is created with the binaries attached
4. Release notes can be edited on the GitHub releases page

### Semantic Versioning

WebKit follows [Semantic Versioning](https://semver.org/):

- **Patch** (v1.0.1): Bug fixes and minor changes
- **Minor** (v1.1.0): New features, backwards compatible
- **Major** (v2.0.0): Breaking changes

---

## Contributing

We welcome contributions! Here's how to get started:

### Making Changes

1. **Fork and clone** the repository
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Make your changes** and add tests
4. **Run tests**: `go test ./...`
5. **Run linting**: `go fmt ./... && go vet ./...`
6. **Commit**: Follow [Conventional Commits](https://www.conventionalcommits.org/)
7. **Push** and open a pull request

### Development Guidelines

- **Test coverage**: Maintain or improve test coverage (aim for >80%)
- **Documentation**: Update docs for user-facing changes
- **Breaking changes**: Clearly document in PR description
- **Code style**: Run `go fmt` and `go vet` before committing
- **Integration tests**: Add tests for CLI commands when modifying behavior

### Testing Your Changes

Before submitting a PR:

```bash
# Run full test suite
go test -v -race ./...

# Test CLI commands in playground
cd internal/playground
../../webkit validate
../../webkit update
```

The `internal/playground` directory is used during development to test CLI commands without
modifying your actual project files.

---

## Platform Modules

The `platform/` directory contains Terraform modules for supported infrastructure providers:

- **DigitalOcean**: Droplets, Postgres, Spaces (S3), domain records
- More providers coming soon

To add a new provider or resource type, see
the [Platform Module Guide](https://webkit.ainsley.dev/contributing/platform-modules).

---

## Troubleshooting

### Common Issues

**`webkit validate` fails with schema errors:**

- Ensure your `app.json` matches the schema at the specified `$schema` URL
- Run `webkit validate` to see detailed error messages

**SOPS decryption fails:**

- Check that `SOPS_AGE_KEY` environment variable is set
- Verify age key has correct permissions: `chmod 600 ~/.config/webkit/age.key`

**Terraform state conflicts:**

- Ensure only one CI/CD pipeline is running at a time
- Check state backend configuration in generated workflows

For more help, see the [documentation](https://webkit.ainsley.dev) or open an issue.


---

## Acknowledgments

- Built by [ainsley.dev](https://ainsley.dev)
- Powered by [Terraform](https://www.terraform.io/), [SOPS](https://github.com/getsops/sops),
  and [Age](https://github.com/FiloSottile/age)
- Template system inspired by Go's `text/template` and [Sprig](https://github.com/Masterminds/sprig)
