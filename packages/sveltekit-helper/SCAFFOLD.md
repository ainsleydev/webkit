# Scaffolding CLI Tool

This document outlines the plan for implementing scaffolding support for SvelteKit components via the WebKit CLI.

## Overview

Add scaffolding support to allow users to generate customisable Svelte components in their project, similar to shadcn's approach. This allows users to own and modify the components rather than importing them from the package.

## Implementation Plan

### Phase 1: CLI Command Structure

**Location**: `internal/cmd/svelte/`

```
internal/cmd/svelte/
├── cmd.go           # Command registration
├── scaffold.go      # Scaffold logic
└── scaffold_test.go # Tests
```

**Command Definition**:
```go
// internal/cmd/svelte/cmd.go
package svelte

import "github.com/urfave/cli/v3"

var Command = &cli.Command{
    Name:        "svelte",
    Usage:       "Manage SvelteKit projects",
    Description: "Commands for working with SvelteKit applications",
    Commands: []*cli.Command{
        ScaffoldCmd,
    },
}
```

**Scaffold Command**:
```go
// internal/cmd/svelte/scaffold.go
var ScaffoldCmd = &cli.Command{
    Name:  "scaffold",
    Usage: "Scaffold SvelteKit components into your project",
    Commands: []*cli.Command{
        {
            Name:  "button",
            Usage: "Scaffold a Button component",
            Action: scaffoldButton,
        },
        {
            Name:  "alert",
            Usage: "Scaffold an Alert component",
            Action: scaffoldAlert,
        },
        // Add more components as needed
    },
}
```

**Register in CLI**:
```go
// internal/cmd/cli.go
import "github.com/ainsleydev/webkit/internal/cmd/svelte"

Commands: []*cli.Command{
    // ... existing commands
    svelte.Command,
}
```

### Phase 2: Template Storage

**Location**: `internal/templates/svelte/`

```
internal/templates/svelte/
├── components/
│   ├── Button.svelte.tmpl
│   ├── Alert.svelte.tmpl
│   ├── FormGroup.svelte.tmpl
│   ├── FormInput.svelte.tmpl
│   └── FormLabel.svelte.tmpl
└── embed.go
```

**Embed Templates**:
```go
// internal/templates/svelte/embed.go
package svelte

import "embed"

//go:embed components/*.svelte.tmpl
var ComponentTemplates embed.FS
```

**Example Template** (`Button.svelte.tmpl`):
```svelte
<script lang="ts">
	export let variant: 'primary' | 'secondary' | 'ghost' = 'primary'
	export let size: 'sm' | 'md' | 'lg' = 'md'
	export let loading = false
	export let disabled = false
</script>

<button
	class="btn btn--{variant} btn--{size}"
	class:btn--loading={loading}
	disabled={disabled || loading}
	{...$$restProps}
>
	{#if loading}
		<span class="btn__loader"></span>
	{/if}
	<slot />
</button>

<style lang="scss">
	.btn {
		$self: &;

		--btn-padding-x: 1rem;
		--btn-padding-y: 0.5rem;
		--btn-bg: var(--colour-blue-500);
		--btn-text: var(--colour-white);
		--btn-border-radius: 0.375rem;

		position: relative;
		display: inline-flex;
		align-items: centre;
		justify-content: centre;
		gap: 0.5rem;
		padding: var(--btn-padding-y) var(--btn-padding-x);
		background: var(--btn-bg);
		colour: var(--btn-text);
		border: none;
		border-radius: var(--btn-border-radius);
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s;

		&:hover:not(:disabled) {
			opacity: 0.9;
		}

		&:disabled {
			opacity: 0.5;
			cursor: not-allowed;
		}

		&--primary {
			--btn-bg: var(--colour-blue-500);
			--btn-text: var(--colour-white);
		}

		&--secondary {
			--btn-bg: var(--colour-grey-200);
			--btn-text: var(--colour-grey-900);
		}

		&--ghost {
			--btn-bg: transparent;
			--btn-text: var(--colour-blue-500);
		}

		&--sm {
			--btn-padding-x: 0.75rem;
			--btn-padding-y: 0.375rem;
			font-size: 0.875rem;
		}

		&--lg {
			--btn-padding-x: 1.5rem;
			--btn-padding-y: 0.75rem;
			font-size: 1.125rem;
		}

		&--loading {
			pointer-events: none;
		}

		&__loader {
			width: 1rem;
			height: 1rem;
			border: 2px solid currentColour;
			border-top-colour: transparent;
			border-radius: 50%;
			animation: spin 0.6s linear infinite;
		}
	}

	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
```

### Phase 3: Scaffold Implementation

**Scaffolding Logic**:
```go
func scaffoldButton(ctx context.Context, cmd *cli.Command) error {
    // 1. Read template from embedded FS
    tmpl, err := template.ParseFS(ComponentTemplates, "components/Button.svelte.tmpl")
    if err != nil {
        return errors.Wrap(err, "parsing button template")
    }

    // 2. Determine output path
    outputPath := "src/lib/components/Button.svelte"

    // 3. Use scaffold generator
    gen := scaffold.NewGenerator(afero.NewOsFs(), manifest, printer)

    // 4. Scaffold file (won't overwrite if exists)
    return gen.Template(outputPath, tmpl, nil,
        scaffold.WithTracking(manifest.SourceScaffold(), "component:button"),
        scaffold.WithScaffoldMode(), // Won't overwrite existing
    )
}
```

### Phase 4: Template Validation

**Strategy**: Validate templates are valid Svelte before embedding.

**Options**:
1. **Pre-commit validation**: Run Svelte compiler on templates during development
2. **Test validation**: Write Go tests that validate templates
3. **CI validation**: GitHub Actions workflow to validate templates

**Recommended Approach**:
```go
// internal/templates/svelte/validate_test.go
func TestTemplatesAreValidSvelte(t *testing.T) {
    t.Parallel()

    entries, err := ComponentTemplates.ReadDir("components")
    require.NoError(t, err)

    for _, entry := range entries {
        t.Run(entry.Name(), func(t *testing.T) {
            t.Parallel()

            content, err := ComponentTemplates.ReadFile("components/" + entry.Name())
            require.NoError(t, err)

            // Validate template syntax
            _, err = template.New(entry.Name()).Parse(string(content))
            assert.NoError(t, err, "template should be valid")

            // Could also call out to svelte compiler for deeper validation
        })
    }
}
```

### Phase 5: Component Library

**Initial Components to Scaffold**:

1. **Button** (`webkit svelte scaffold button`)
   - Variants: primary, secondary, ghost
   - Sizes: sm, md, lg
   - Loading state
   - Disabled state

2. **Alert** (`webkit svelte scaffold alert`)
   - Types: info, success, warning, error
   - Dismissible option
   - Icon slot
   - Title and description

3. **Form Components** (`webkit svelte scaffold form`)
   - FormGroup - Field wrapper with label/error
   - FormInput - Text input with validation states
   - FormLabel - Label with required indicator
   - FormTextarea - Textarea component
   - FormCheckbox - Checkbox with label

4. **Card** (`webkit svelte scaffold card`)
   - Header, body, footer slots
   - Variants: default, outlined, elevated

### Phase 6: Demo/Testing Strategy

**Challenge**: How to demo/test scaffolded components without a SvelteKit app in WebKit repo.

**Solutions**:

1. **Separate Demo Repo**:
   - Create `webkit-sveltekit-demo` repository
   - Use `webkit svelte scaffold` commands to generate components
   - Visual regression testing with Playwright
   - Component playground with Storybook

2. **Test in `internal/playground`**:
   - Add a minimal SvelteKit app to playground
   - Scaffold components during CI
   - Run `pnpm build` to validate

3. **Template Unit Tests**:
   - Test that templates render without errors
   - Validate TypeScript types are correct
   - Ensure SCSS compiles

**Recommended**: Combination of #2 (playground) and #3 (unit tests).

### Phase 7: Documentation

**Update README.md**:
```markdown
## Scaffolding Components

Generate customisable components in your SvelteKit project:

```bash
# Scaffold a button component
webkit svelte scaffold button

# Scaffold form components
webkit svelte scaffold form

# Scaffold an alert component
webkit svelte scaffold alert
```

Components are created in `src/lib/components/` and can be customised to match your design system.
```

**Add to Component READMEs**:
Each component template should include usage documentation:
```svelte
<!--
	@component

	Customisable button component with variants and loading states.

	@example
	```svelte
	<Button variant="primary" size="md">
		Click me
	</Button>
	```

	@example
	```svelte
	<Button variant="secondary" loading>
		Loading...
	</Button>
	```
-->
```

## Implementation Checklist

### Phase 1: Foundation
- [ ] Create `internal/cmd/svelte/` directory
- [ ] Implement `Command` registration
- [ ] Implement `ScaffoldCmd` with button subcommand
- [ ] Register in `internal/cmd/cli.go`
- [ ] Write basic tests

### Phase 2: Templates
- [ ] Create `internal/templates/svelte/` directory
- [ ] Create `Button.svelte.tmpl`
- [ ] Create `embed.go` with `//go:embed`
- [ ] Add template validation tests

### Phase 3: Scaffolding Logic
- [ ] Implement `scaffoldButton()` function
- [ ] Use `scaffold.Generator` with `ModeScaffold`
- [ ] Add manifest tracking
- [ ] Test scaffolding workflow

### Phase 4: Additional Components
- [ ] Add `Alert.svelte.tmpl`
- [ ] Add form component templates
- [ ] Add card template
- [ ] Update command with all subcommands

### Phase 5: Testing & Validation
- [ ] Add template validation tests
- [ ] Add integration tests
- [ ] Test in playground app
- [ ] Validate TypeScript compilation
- [ ] Validate SCSS compilation

### Phase 6: Documentation
- [ ] Update main README with scaffolding section
- [ ] Add component documentation to templates
- [ ] Create SCAFFOLDING.md guide
- [ ] Add examples to docs

## Usage Examples

```bash
# Scaffold a button component
webkit svelte scaffold button
# Creates: src/lib/components/Button.svelte

# Scaffold form components
webkit svelte scaffold form
# Creates:
#   src/lib/components/Form/FormGroup.svelte
#   src/lib/components/Form/FormInput.svelte
#   src/lib/components/Form/FormLabel.svelte

# Scaffold an alert
webkit svelte scaffold alert
# Creates: src/lib/components/Alert.svelte
```

## Design Decisions

1. **Go CLI over JS CLI**: Better integration with existing WebKit architecture
2. **Template validation via tests**: Ensures templates are valid before embedding
3. **Scaffold mode**: Uses `ModeScaffold` so files won't be overwritten
4. **Manifest tracking**: All scaffolded files tracked in `.webkit/manifest.json`
5. **SCSS with BEM**: All templates use SCSS with BEM naming convention
6. **CSS variables**: Maximum customisation flexibility
7. **Headless components**: No opinionated styling, just structure

## Future Enhancements

1. **Interactive scaffolding**: Prompt for component options (variants, features)
2. **Batch scaffolding**: `webkit svelte scaffold all` to generate all components
3. **Custom templates**: Allow users to define their own component templates
4. **Theme configuration**: Generate theme config file with colour tokens
5. **Component variants**: `webkit svelte scaffold button --variant=outline`
