# WebKit

A webkit framework and SDK for ainsley.dev

## Releasing

WebKit uses [GoReleaser](https://goreleaser.com/) for automated releases. The release process is
triggered by creating and pushing a git tag.

### Quick Release

Use the interactive tag tool:

```bash
pnpm release:tag
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

### Version Numbering

WebKit follows [Semantic Versioning](https://semver.org/):

- **Patch** (v1.0.1): Bug fixes and minor changes
- **Minor** (v1.1.0): New features, backwards compatible
- **Major** (v2.0.0): Breaking changes
