# npm packages

This directory contains npm packages that are published to the npm registry. WebKit
uses [Changesets](https://github.com/changesets/changesets) to manage versioning and publishing of
these packages automatically.

## How it works

When you make changes to packages in this directory and merge them to the `main` branch, a GitHub
Action workflow automatically handles versioning and publishing to npm. Changesets provides a
structured way to track which packages have changed and what type of version bump they require (
patch, minor, or major).

### The workflow

1. **Make changes** to one or more packages in a pull request.
2. **Create a changeset** using `pnpm changeset` to document your changes.
3. **Merge to main** - the GitHub Action detects the changeset.
4. **Version PR created** - a pull request is automatically created with updated package versions
   and changelogs.
5. **Merge version PR** - once merged, packages are automatically published to npm.

After publishing, the applied changeset files are deleted automatically, keeping the repository
clean for the next release cycle.

## Creating a changeset

When you make changes to a package, you need to create a changeset file that describes what changed
and the version bump type.

### Using the CLI

Run the following command from the repository root:

```bash
pnpm changeset
```

This interactive CLI will:

1. Ask which packages have changed (select using space bar, confirm with enter).
2. Ask whether the change is a `patch`, `minor`, or `major` bump for each package.
3. Prompt you to write a summary of the changes.

The CLI creates a new markdown file in `.changeset/` with a randomly generated name. This file
should be committed with your pull request.

### Version bump types

Choose the appropriate version bump based on [Semantic Versioning](https://semver.org/):

- **Patch** (0.0.X) - Bug fixes and minor changes that don't affect the API.
- **Minor** (0.X.0) - New features that are backwards compatible.
- **Major** (X.0.0) - Breaking changes that are not backwards compatible.

### Example

```bash
$ pnpm changeset
🦋  Which packages would you like to include?
◉ @ainsleydev/payload-helper

🦋  Which packages should have a major bump?
◯ @ainsleydev/payload-helper

🦋  Which packages should have a minor bump?
◯ @ainsleydev/payload-helper

🦋  Which packages should have a patch bump?
◉ @ainsleydev/payload-helper

🦋  Please enter a summary for this change (this will be in the changelogs).
Summary › Fixed bug in focal point field generation

✅ Changeset added! - commit it now!
```

## Publishing workflow

### Automatic publishing (recommended)

When you merge a pull request with changesets to `main`:

1. The `publish.yaml` workflow runs automatically.
2. If changesets exist, a "Version Packages" pull request is created/updated.
3. The version PR updates `package.json` versions and `CHANGELOG.md` files.
4. When you merge the version PR, packages are automatically published to npm.

## Available commands

From the repository root, you can run:

- `pnpm changeset` - Create a new changeset.
- `pnpm changeset:version` - Update package versions based on changesets.
- `pnpm changeset:publish` - Build and publish packages to npm.

## NPM authentication

The GitHub Actions workflow uses the `ORG_NPM_TOKEN` secret for authentication with npm. This token
is configured at the organisation level and has permission to publish packages under the
`@ainsleydev` scope.

For local publishing, you'll need to authenticate with npm:

```bash
npm login
```

## Multiple changesets

You can create multiple changesets in a single pull request if you're making multiple distinct
changes. Each changeset is processed independently, and the version bumps are combined
intelligently (e.g., if one changeset specifies `patch` and another specifies `minor`, the final
bump will be `minor`).
