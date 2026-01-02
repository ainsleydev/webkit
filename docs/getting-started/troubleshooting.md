# Troubleshooting

Common issues and solutions when working with WebKit projects.

## Payload CMS migration checks

### Migration check failing in CI but passing locally

**Symptom**: Migration check passes locally but fails in CI with "Schema changes detected"

**Cause**: Your local `node_modules` is out of sync with `pnpm-lock.yaml`. This happens when you run `git pull` but forget to run `pnpm install`.

**Solution**: Always run `pnpm install` after pulling changes:

```bash
git pull
pnpm install
```

**Prevention**: WebKit automatically creates a `scripts/check-deps.js` file in Payload apps that validates dependencies are in sync before running migrations. The `migrate:create` script runs this check automatically:

```bash
# This command now checks dependencies first
pnpm migrate:create
```

If dependencies are out of sync, you'll see:

```
❌ Dependencies out of sync!
   pnpm-lock.yaml is newer than node_modules
   Run: pnpm install
```

### Understanding the migration workflow

When working with Payload CMS migrations:

1. **Pull changes**: `git pull`
2. **Install dependencies**: `pnpm install`
3. **Check for migrations**: `pnpm migrate:create`
4. **Review generated migration** (if any)
5. **Commit migration file** (if generated)
6. **Push changes**: `git push`

The dependency check ensures step 2 isn't skipped, preventing CI failures.

## Environment variables

### Missing environment variables in production

**Symptom**: Application fails to start with "Missing required environment variable"

**Cause**: Environment variables from SOPS or resources haven't been generated or synced.

**Solution**:

```bash
# For SOPS variables
webkit env generate

# For resource variables (after Terraform apply)
webkit infra apply
```

## Drift detection

### Drift detected for files I didn't modify

**Symptom**: `webkit drift` reports changes to generated files you haven't touched

**Cause**: WebKit templates have been updated since you last ran `webkit update`

**Solution**: Run `webkit update` to regenerate files with latest templates:

```bash
webkit update
```

### Preserving manual changes to generated files

**Symptom**: Your changes to a generated file are overwritten by `webkit update`

**Cause**: The file is marked as "generated" rather than "scaffolded"

**Solution**: Generated files (like GitHub workflows) are intentionally overwritten to stay in sync with your manifest. If you need custom behaviour:

1. Modify your `app.json` manifest instead, or
2. Copy the file to a new location and reference it from your custom code

## Infrastructure issues

### Terraform state lock errors

**Symptom**: "Error acquiring state lock"

**Cause**: Another process is running Terraform, or a previous run didn't clean up properly

**Solution**:

1. Wait for other operations to complete
2. If lock is stuck: `terraform force-unlock <lock-id>`
3. Check CI/CD pipelines for running infrastructure jobs

### Resource already exists errors

**Symptom**: "Resource already exists" during `webkit infra apply`

**Cause**: The resource exists in your cloud provider but not in Terraform state

**Solution**: Import the existing resource:

```bash
webkit infra import
```

This brings existing cloud resources under WebKit management without recreating them.

## Build failures

### Docker build fails with missing dependencies

**Symptom**: Docker build fails during `pnpm install`

**Cause**: Dependencies not properly specified in `package.json`

**Solution**:

1. Ensure `pnpm-lock.yaml` is committed to git
2. Run `pnpm install` locally to verify lockfile is correct
3. Check Docker build logs for specific missing packages

### Turbo build cache issues

**Symptom**: Builds succeed locally but fail in CI, or stale code is deployed

**Cause**: Turbo cache is stale or corrupted

**Solution**:

```bash
# Clear Turbo cache
pnpm turbo clean

# Rebuild everything
pnpm turbo build --force
```

## Next steps

If you're still experiencing issues:

1. Check the [CLI reference](/cli/overview) for command details
2. Review the [manifest reference](/manifest/overview) for configuration options
3. Open an issue on [GitHub](https://github.com/ainsleydev/webkit/issues)
