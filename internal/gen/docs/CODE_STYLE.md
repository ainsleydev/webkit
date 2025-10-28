## HTML

### Indentation

Use tabs instead of spaces for markup. Do not mix tabs with spaces, ensure it is properly formatted.

## SCSS

### Naming

Use BEM (Block Element Modifier) notation for CSS class names.

## Go

### Formatting

Use `gofmt` for standard Go formatting.

### Error handling

Always check and handle errors appropriately.

## JS

### Style

Use `camelCase` for all field names and variables.

## Git

### Commit Messages

Follow a conventional commit format with a type prefix and present tense gerund (doing words):

#### Format

```
<type>: <description>
```

#### Types

- `feat:` - Adding new features or functionality.
- `fix:` - Fixing bugs or issues.
- `chore:` - Updating dependencies, linting, or other maintenance tasks.
- `style:` - Refactoring code or improving code style (no functional changes).
- `test:` - Adding or updating tests.
- `docs:` - Updating documentation.

### Pre-Commit Checklist

Before submitting changes, agents should verify the following:

- [ ] All tests pass locally.
- [ ] Code is properly formatted.
- [ ] New exported types, functions, and constants have documentation comments.
