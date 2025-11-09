# TODO

## Documentation

Create and update the `docs` folder with coherent documentation for WebKit.

## README Generation

Create beautiful looking README's from the `app.json` data.

## Testing

### Standardised test setup utilities

Create consistent test helper patterns in `internal/util/testutil/` for common test setup scenarios:

- **CommandInput setup**: Helper for creating `cmdtools.CommandInput` with MemMapFs and all required dependencies.
- **FileGenerator setup**: Helper for creating `scaffold.FileGenerator` instances in tests.
- **AppDef setup**: Helper for creating and validating `appdef.Definition` instances.

**Current inconsistencies:**

- `internal/scaffold/generate_test.go` - Has `setup(t) *FileGenerator`.
- `internal/cmd/payload/cmd_test.go` - Has `setup(t) (afero.Fs, cmdtools.CommandInput)`.
- Many test files create these manually without helpers.

**Benefits:**

- Reduces boilerplate in test files.
- Ensures consistent test setup across the codebase.
- Makes tests more readable by hiding setup complexity.
- Easier to update test patterns globally.

## Code Quality

### Error wrapping consistency

Consider adding a linter rule to enforce consistent error wrapping patterns across the codebase:

- **Prefer**: `errors.Wrap(err, "context")` from `github.com/pkg/errors`.
- **Alternative**: `fmt.Errorf("context: %w", err)` only when needing to format multiple arguments.

**Current state:**

- AGENTS.md documents the preferred pattern.
- Most code follows the `errors.Wrap` pattern.
- Occasional `fmt.Errorf` usage exists.

**Potential solutions:**

- Add `errorlint` to golangci-lint configuration.
- Add custom `forbidigo` pattern to warn about `fmt.Errorf.*%w` when it could be `errors.Wrap`.
- Create a custom linter using `go-ruleguard`.

## Misc

- BetterStack/OneUptime Providers for Infra.
- Improve Coverage.
- Create an infra plan —destroy command. So we can see whats destroyed?
- Seed utilities for Payload, it's a pain having to do it all the time.


