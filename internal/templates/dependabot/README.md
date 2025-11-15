# Dependabot + pnpm Monorepo Solution

This directory contains the solution for handling Dependabot PRs in pnpm monorepos.

## The Problem

When Dependabot updates `package.json` files in a pnpm workspace, it doesn't update the `pnpm-lock.yaml` file. This causes CI failures with the error:

```
ERR_PNPM_OUTDATED_LOCKFILE  Cannot install with "frozen-lockfile" because pnpm-lock.yaml is not up to date
```

## The Solution

We've implemented a **GitHub Actions workflow** that automatically:

1. Detects Dependabot PRs (when `package.json` files change)
2. Runs `pnpm install --no-frozen-lockfile` to update the lockfile
3. Commits the updated `pnpm-lock.yaml` back to the PR

### Workflow Location

The auto-fix workflow is located at:
```
.github/workflows/dependabot-pnpm-lockfile.yaml
```

### How It Works

- **Trigger**: Runs on PRs when any `package.json` file changes
- **Condition**: Only executes if the actor is `dependabot[bot]`
- **Action**: Updates and commits the lockfile automatically
- **Permissions**: Requires `contents: write` and `pull-requests: write`

### Benefits

✅ **Automatic** - No manual intervention needed
✅ **Simple** - Single workflow file, no complex configuration
✅ **Safe** - Only runs on Dependabot PRs
✅ **Fast** - Uses pnpm caching for speed

## Additional Recommendations

### Fix "pnpm" Field Warning

If you see this warning:
```
WARN  The field "pnpm" was found in /path/to/package.json.
This will not take effect. You should configure "pnpm" at the root of the workspace instead.
```

**Solution**: Remove the `"pnpm"` field from workspace package.json files. Keep it only in the root `package.json`.

### Dependabot Configuration Best Practices

For pnpm monorepos, use these settings in `.github/dependabot.yaml`:

```yaml
- package-ecosystem: "npm"  # Use "npm" for pnpm
  directory: "/path/to/workspace"
  schedule:
    interval: "weekly"
  groups:
    dependencies:
      patterns:
        - "*"
```

## Alternative Solutions

If this approach doesn't work for your use case:

1. **Renovate Bot** - Better pnpm monorepo support out of the box
2. **Manual commits** - Manually run `pnpm install` on Dependabot PRs
3. **Disable frozen-lockfile** - Not recommended for production

## Troubleshooting

### Workflow not running?

- Check that the workflow file exists in `.github/workflows/`
- Ensure the repository has Actions enabled
- Verify permissions are set correctly

### Commits not appearing?

- Check GitHub Actions logs for errors
- Ensure `GITHUB_TOKEN` has write permissions
- Verify branch protection rules allow bot commits

### Still getting lockfile errors?

- Manually run `pnpm install` locally
- Check for pnpm version mismatches
- Ensure `pnpm-workspace.yaml` is correctly configured
