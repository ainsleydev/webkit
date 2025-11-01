# Style Conventions

## Go Tests

### Test Naming Convention

All test function names should use PascalCase after the underscore, regardless of whether the method being tested is exported or unexported.

**Examples:**
- `TestApp_DefaultPort` (tests private method `defaultPort()`)
- `TestApp_OrderedCommands` (tests public method `OrderedCommands()`)
- `TestApp_ShouldUseNPM` (tests public method `ShouldUseNPM()`)

This maintains consistency and improves readability across the test suite.
