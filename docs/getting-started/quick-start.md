# Quick start

Get up and running with WebKit in under 5 minutes. This guide walks you through creating a minimal configuration and generating your first project files.

## Create your manifest

WebKit uses a single `app.json` file to define your entire project. Create one in your project root:

```json
{
  "$schema": "https://raw.githubusercontent.com/ainsleydev/webkit/main/schema.json",
  "project": {
    "name": "my-project",
    "title": "My Project",
    "repo": "github.com/myorg/my-project"
  },
  "apps": [
    {
      "name": "web",
      "type": "svelte-kit",
      "path": "./apps/web"
    }
  ]
}
```

The `$schema` property enables IDE autocompletion and validation.

## Run WebKit update

Generate all project files by running:

```bash
webkit update
```

WebKit analyses your manifest and generates:

- **GitHub Actions workflows** - CI/CD pipelines for testing, building, and deploying
- **Docker configuration** - `.dockerignore` and build settings
- **Project files** - `package.json`, `pnpm-workspace.yaml`, `turbo.json`
- **Git configuration** - `.gitignore`, `.editorconfig`
- **Manifest tracking** - `.webkit/manifest.json` for change detection

## What was generated?

After running `webkit update`, your project structure looks like this:

```
my-project/
├── .github/
│   ├── actions/           # Reusable workflow actions
│   └── workflows/         # CI/CD pipelines
├── .webkit/
│   └── manifest.json      # Tracks generated files
├── apps/
│   └── web/               # Your SvelteKit app
├── app.json               # WebKit manifest
├── package.json           # Root package configuration
├── pnpm-workspace.yaml    # Workspace configuration
└── turbo.json             # Turborepo configuration
```

## Validate your configuration

Before making changes, validate your manifest:

```bash
webkit validate
```

This checks for:
- Required fields
- Valid app types and paths
- Resource configuration
- Environment variable references

## Idempotent updates

WebKit is designed to be run repeatedly. Each time you run `webkit update`:

1. It reads the current manifest from `.webkit/manifest.json`
2. Generates files based on your `app.json`
3. Writes a new manifest
4. Cleans up orphaned files (files that were generated but are no longer needed)

This means you can safely run `webkit update` whenever you change your `app.json`.

## Detect drift

Check if any generated files have been manually modified:

```bash
webkit drift
```

WebKit tracks content hashes of all generated files. If you've modified a generated file, drift detection warns you before overwriting your changes.

## Next steps

Now that you understand the basics:

- Build a complete [portfolio site](/getting-started/your-first-project) with deployment
- Learn about [core concepts](/getting-started/core-concepts) like manifest tracking and drift detection
- Explore the [manifest reference](/manifest/overview) for all configuration options
