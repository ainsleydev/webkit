# @ainsleydev/sveltekit-helper - Context & Design Decisions

This document captures the context, design decisions, and architectural choices made during the creation of this package.

## Overview

**Purpose**: SvelteKit utilities, components and helpers for ainsley.dev builds.

**Philosophy**: Keep it simple, lightweight, and let consumers own their styling.

## Key Design Decisions

### 1. Svelte 5 with Runes

All components use the latest Svelte 5 runes API:
- `$state()` for reactive state
- `$derived` for computed values
- `$props()` for component props
- `$restProps` for spreading additional attributes

**Reasoning**: Modern API, better type safety, more explicit reactivity.

### 2. SCSS with BEM Naming Convention

All components use SCSS with BEM (Block Element Modifier) naming:
```scss
.component {
  $self: &;

  &__element {
    // Element styles
  }

  &--modifier {
    // Modifier styles
  }
}
```

**Reasoning**: Provides clear component structure while allowing easy overrides via CSS variables or class targeting.

### 3. CSS Variables for Customization

Every component exposes CSS variables for customization:
- `--container-max-width`
- `--row-gap`
- `--col-gap`
- `--form-input-padding`
- `--form-button-bg`
- etc.

**Reasoning**: Allows consumers to customize without touching component code or scaffolding.

### 4. Minimal Form Utilities

Originally considered complex `clientForm` and `serverForm` utilities with stores, validation orchestration, and SvelteKit integration.

**Decision**: Removed complexity. Kept only:
- `generateFormSchema()` - Generates Zod schema from Payload form fields
- `flattenZodErrors()` - Converts Zod errors to simple key-value object

**Reasoning**: Consumers have different form needs. Better to provide simple helpers than complex abstractions. PayloadForm component handles its own state with simple Svelte 5 runes.

### 5. Consumer-Defined Grid Classes

Originally included all grid classes (.col-1 through .col-12, responsive variants) in Column component.

**Decision**: Removed all grid classes. Column now only provides base structure with `--col-gap` CSS variable.

**Reasoning**:
- Reduces package size
- Avoids repetitive code
- Gives consumers full control over their grid system
- Different projects have different breakpoint needs

Example consumer implementation:
```css
.col-12 { width: 100%; }
.col-6 { width: 50%; }

@media (min-width: 768px) {
  .col-tab-6 { width: 50%; }
}
```

### 6. Two Distribution Methods (Phase 1 & 2)

**Phase 1 (Current)**: Direct npm import
- Install via `pnpm add @ainsleydev/sveltekit-helper`
- Import stable components directly
- Works for Grid, form utilities, Payload components

**Phase 2 (Future)**: CLI scaffolding
- Command: `webkit svelte scaffold button`
- For highly customizable components
- See SCAFFOLD.md for implementation details

**Reasoning**: Some components (Grid, PayloadMedia) are stable and don't need customization. Others (Button, Alert, Form inputs) benefit from being scaffolded so consumers can modify them directly.

## Architecture

### Package Structure

```
packages/sveltekit-helper/
├── src/
│   ├── components/
│   │   ├── grid/           # Container, Row, Column
│   │   └── payload/        # PayloadForm, PayloadMedia
│   ├── utils/
│   │   └── forms/          # generateFormSchema, flattenZodErrors
│   └── index.ts
├── tests/
├── package.json
├── README.md
├── SCAFFOLD.md             # Future CLI scaffolding plans
└── CONTEXT.md              # This file
```

### Component Exports

Clean import paths without `/dist`:
```typescript
import { Container, Row, Column } from '@ainsleydev/sveltekit-helper/components/grid'
import { PayloadForm, PayloadMedia } from '@ainsleydev/sveltekit-helper/components/payload'
import { generateFormSchema, flattenZodErrors } from '@ainsleydev/sveltekit-helper/utils/forms'
```

### Build Configuration

- **Compiler**: SWC (via `@sveltejs/package`)
- **Output**: TypeScript declarations + compiled JavaScript
- **Testing**: Vitest
- **Linting**: Biome

## Component Details

### Grid System

**Container**: CSS Grid with breakout layout support
- Main content area: `--container-max-width` (default: 1328px)
- Breakout area: `--container-breakout-max-width` (default: 1500px)
- Padding: `--container-padding` (default: 1rem)

**Row**: Flexbox row with gap management
- Gap: `--row-gap` (default: 1rem)
- Optional `noGaps` prop for flush layouts

**Column**: Base column structure only
- Gap: `--col-gap` (default: 1rem)
- Consumers define grid classes

### Payload CMS Integration

**PayloadForm**: Dynamic form rendering
- Renders forms from Payload CMS form builder
- Supports text, email, number, textarea, checkbox fields
- Built-in validation display
- Customizable via CSS variables
- Simple Svelte 5 runes for state management (no complex stores)
- Props: `form`, `apiEndpoint`, `onSubmit`

**PayloadMedia**: Responsive media component
- Handles images (with responsive sizes) and videos
- Automatic format prioritization: AVIF → WebP → JPEG/PNG
- SVG support
- Props: `data`, `loading`, `maxWidth`, `breakpointBuffer`, `className`, `onload`

### Form Utilities

**generateFormSchema**: Creates Zod schema from Payload fields
- Supports all Payload form field types
- Handles required/optional validation
- Email validation for email fields

**flattenZodErrors**: Simplifies Zod error objects
- Converts nested error structure to flat key-value pairs
- Takes first error message per field

## Peer Dependencies

**Required**:
- `svelte@^5.0.0`
- `@sveltejs/kit@^2.0.0`

**Optional**:
- `payload@^3.0.0` - For Payload CMS components
- `zod@^3.0.0` - For form validation utilities

## Development Workflow

```bash
# Install dependencies
pnpm install

# Build the package
pnpm --filter @ainsleydev/sveltekit-helper build

# Run tests
pnpm --filter @ainsleydev/sveltekit-helper test

# Lint and format
pnpm --filter @ainsleydev/sveltekit-helper lint
pnpm --filter @ainsleydev/sveltekit-helper format
```

## Lessons Learned

1. **Start Simple**: Initial implementation was too complex (clientForm/serverForm stores). Simplified version is better.

2. **Let Consumers Control Styling**: Providing too many opinionated classes (grid classes) was unnecessary. CSS variables + BEM naming gives consumers full control.

3. **TypeScript Module Resolution**: With `"moduleResolution": "Node16"`, all relative imports need `.js` extensions even though files are `.ts`.

4. **Type Exports from Svelte**: Can't export types from `.svelte` module blocks. Need separate `types.ts` files.

5. **Props Must Use `let`**: In Svelte 5, component props must use `let` not `const` to allow them to be reactive and overridable.

## Future Plans

See [SCAFFOLD.md](./SCAFFOLD.md) for Phase 2 CLI scaffolding implementation.

**Potential Additions**:
- Scaffold: Button component
- Scaffold: Alert/Toast component
- Scaffold: Form input components (Input, Select, Checkbox, Radio)
- Scaffold: Card component
- Utility: Animation helpers
- Utility: Accessibility helpers

## Related Repositories

- **sveltekit-boilerplate**: `/Users/ainsley.clark/Desktop/ainsley.dev/sveltekit-boilerplate`
  - Source of Grid component inspiration
  - SCSS architecture reference

- **search-spares**: Project using latest Payload CMS patterns
  - PayloadMedia component reference
  - Context API form patterns

## References

- [Svelte 5 Runes Documentation](https://svelte.dev/docs/svelte/what-are-runes)
- [BEM Methodology](https://getbem.com/)
- [Payload CMS](https://payloadcms.com/)
- [Zod](https://zod.dev/)
