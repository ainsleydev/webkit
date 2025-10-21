# webkit update

The `webkit update` command regenerates all project files based on your `app.json` manifest. This is the primary command you'll use when working with WebKit.

## Synopsis

```bash
webkit update
```

## Description

`webkit update` reads your `app.json` manifest and ensures all generated files are in sync. It:

1. **Loads the previous manifest** (`.webkit-manifest.json`) to track what was previously generated
2. **Generates new files** based on the current `app.json`:
   - Code style configurations (EditorConfig, Prettier, Biome)
   - Git settings (`.gitignore`, GitHub repository settings)
   - Environment files (`.env`, `.env.production`)
   - Secrets scaffolding (SOPS configuration, encrypted files)
   - CI/CD workflows (GitHub Actions)
   - Project tooling (package.json, pnpm-workspace.yaml, turbo.json)
3. **Compares manifests** to detect changes
4. **Updates the `webkit_version`** field in `app.json` to match the CLI version
5. **Cleans up orphaned files** that no longer have a source in `app.json`
6. **Preserves customisations** you've made to generated files

## When to use

Run `webkit update` whenever you:

- Create a new `app.json` manifest
- Add or remove apps or resources
- Modify environment variables
- Change infrastructure configuration
- Update to a new version of WebKit

**Rule of thumb:** If you edit `app.json`, run `webkit update`.

## Examples

### Basic usage

```bash
webkit update
```

Output:
```
Updating project dependencies...

🏃 Manifest: Scaffold manifest files
🏃 Definition: Update webkit_version in app.json
🏃 Env: Scaffold .env files
🏃 Secrets: Scaffold secret files
🏃 Files: Create code style files
🏃 Files: Create git settings
🏃 Files: Create package.json
🏃 CICD: Create app PR workflows
🏃 CICD: Creates drift detection workflow
🏃 Env: Sync .env files
🏃 Secrets: Sync secret files

✓ Successfully updated project dependencies!
```

### After adding a new app

Edit `app.json` to add an app:

```json
{
  "apps": [
    {
      "name": "api",
      "type": "go",
      "path": "services/api"
    }
  ]
}
```

Run the update:

```bash
webkit update
```

WebKit will:
- Generate `.env` files in `services/api/`
- Create GitHub Actions workflows for building the API
- Add the app to monorepo configurations (if applicable)
- Create Docker-related files

### After removing an app

Remove an app from `app.json`, then:

```bash
webkit update
```

WebKit automatically:
- Deletes orphaned workflow files
- Removes the app from monorepo configs
- Cleans up tracked files that no longer have a source

## What gets generated

### Project-level files

- `.editorconfig` - Editor configuration
- `.prettierrc` - Prettier configuration
- `biome.json` - Biome linter/formatter config
- `.gitignore` - Git ignore patterns
- `.github/settings.yml` - Repository settings
- `package.json` - Monorepo root manifest
- `pnpm-workspace.yaml` - pnpm workspace configuration
- `turbo.json` - Turborepo configuration
- `.webkit-manifest.json` - File tracking manifest

### Per-app files

For each app defined in `apps`:

- `{app.path}/.env` - Development environment variables
- `{app.path}/.env.production` - Production environment variables
- `{app.path}/.dockerignore` - Docker ignore patterns
- `.github/workflows/pr-{app.name}.yml` - Pull request workflow

### Secrets files

- `.sops.yaml` - SOPS configuration
- `secrets/shared.yaml` - Shared secrets (encrypted)
- `secrets/{app.name}.yaml` - Per-app secrets (encrypted)

### Infrastructure workflows

- `.github/workflows/drift-detection.yml` - Infrastructure drift monitoring
- `.github/workflows/backup-{resource.name}.yml` - Database backup workflows

## How file tracking works

WebKit uses `.webkit-manifest.json` to track generated files. This manifest contains:

```json
{
  "version": "1.0.0",
  "files": {
    ".github/workflows/pr-web.yml": {
      "source": "app:web",
      "hash": "abc123...",
      "modified": false
    }
  }
}
```

**Key points:**

- `source`: Which app, resource, or project setting generated this file
- `hash`: Checksum of the file content
- `modified`: Whether you've customised the file

When you run `webkit update`:
1. WebKit generates files based on `app.json`
2. It compares file hashes with the manifest
3. If a file was modified (`modified: true`), WebKit preserves your changes
4. If a file's source no longer exists in `app.json`, it's deleted
5. The manifest is updated to reflect the new state

::: warning
Do not edit `.webkit-manifest.json` directly. WebKit manages this file automatically.
:::

## Handling customisations

WebKit detects when you customise generated files and preserves those customisations on subsequent updates.

**Example:**

1. Run `webkit update` to generate `.github/workflows/pr-web.yml`
2. You add a custom step to the workflow
3. Run `webkit update` again
4. WebKit detects the file was modified and preserves your custom step

**Limitations:**

If the underlying template changes significantly (e.g., adding a required step), WebKit may regenerate the file and lose customisations. In this case, WebKit warns you and you'll need to reapply your changes.

To avoid conflicts:
- Keep customisations minimal
- Document custom changes in comments
- Consider using separate custom workflows for complex additions

## Idempotency

`webkit update` is idempotent—running it multiple times with the same `app.json` produces the same result:

```bash
webkit update
# No changes

webkit update
# Still no changes
```

This makes it safe to run frequently. You can run it after every `app.json` edit without worrying about breaking things.

## Updating WebKit version

When you upgrade WebKit, run `webkit update` to:

- Update the `webkit_version` field in `app.json`
- Regenerate files with new templates
- Apply any new defaults or features

Example:

```bash
# Upgrade WebKit
brew upgrade webkit

# Update project files
webkit update
```

## Troubleshooting

### Update fails with validation error

```
Error: validation failed: app "web" references non-existent resource "db"
```

Fix the issue in `app.json`, then run `webkit update` again.

### Files aren't updating

If generated files aren't reflecting changes in `app.json`:

1. Check that you edited `app.json` correctly (valid JSON)
2. Ensure the `webkit_version` matches your installed CLI version
3. Try deleting `.webkit-manifest.json` to force a full regeneration (⚠️ this will overwrite any customisations)

### Customisations were lost

If WebKit overwrote your customisations:

1. Check `git diff` to see what changed
2. Reapply your customisations
3. Consider moving custom logic to separate files that WebKit doesn't manage

## See also

- **[webkit drift](/cli/webkit-drift)** - Validate your manifest
- **[webkit scaffold](/cli/webkit-scaffold)** - Generate individual files
- **[Manifest reference](/manifest/overview)** - Learn about `app.json`
- **[Core concepts - Idempotent updates](/core-concepts/overview#idempotent-updates)**
