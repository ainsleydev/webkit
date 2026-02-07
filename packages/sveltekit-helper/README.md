# @ainsleydev/sveltekit-helper

SvelteKit utilities, components and helpers for ainsley.dev builds.

## Installation

```bash
pnpm add @ainsleydev/sveltekit-helper
```

## Features

- **Grid System**: Responsive Container, Row, and Column components with CSS variables
- **Navigation Components**: Mobile-first Sidebar and Hamburger menu components
- **Form Utilities**: Schema generation and error helpers for Zod validation
- **Payload CMS Integration**: Ready-to-use components for Payload CMS forms and media
- **SCSS with BEM**: All components use SCSS with BEM naming convention

## Grid Components

### CSS Variable Customization

All Grid components use CSS variables with fallback values, allowing flexible customization:

**Override Priority (highest to lowest):**
1. Inline styles: `<Container style="--container-padding: 2rem">`
2. Page/component-scoped: `.pricing-page { --container-padding: 3rem; }`
3. Global: `:root { --container-padding: 2rem; }`
4. Component defaults: Defined in each component's `<style>` block

**Responsive Variables:**

Row and Column components include mobile-specific overrides (< 568px). You can customise responsive behaviour using:

```css
:root {
	/* Override both desktop and mobile */
	--row-gap: 1.5rem;

	/* Override mobile only */
	--row-gap-mobile: 0.75rem;
	--col-gap-mobile: 0.75rem;
}
```

**Fallback chain on mobile:**
1. `--row-gap-mobile` (if set)
2. `--row-gap` (if set)
3. `0.5rem` (component default)

### Container

Center content horizontally with predefined max-width and support for breakout layouts.

```svelte
<script>
	import { Container, Row, Column } from '@ainsleydev/sveltekit-helper/components/grid'
</script>

<Container>
	<Row>
		<Column class="col-12 col-desk-6">
			Content
		</Column>
	</Row>
</Container>
```

#### Customisation

Override CSS variables globally from `:root`:

```css
/* Global override for ALL containers */
:root {
	--container-padding: 2rem;
	--container-max-width: 1400px;
	--container-breakout-max-width: 1600px;
}

/* Page-specific override */
.pricing-page {
	--container-padding: 3rem;
}
```

Or use inline styles for single instances:

```svelte
<Container style="--container-padding: 2rem">
	<Row>...</Row>
</Container>
```

### Row

Flexbox row container with gap management.

```svelte
<Row>
	<Column class="col-12 col-tab-6">
		Column 1
	</Column>
	<Column class="col-12 col-tab-6">
		Column 2
	</Column>
</Row>

<!-- No gaps -->
<Row noGaps>
	<Column class="col-6">Content</Column>
</Row>
```

#### Customisation

```css
/* Global override */
:root {
	--row-gap: 1.5rem;
	--row-gap-mobile: 0.75rem; /* Optional: mobile-specific gap (< 568px) */
}
```

Or use inline styles:

```svelte
<Row style="--row-gap: 0.5rem">
	<Column>...</Column>
</Row>
```

### Column

Base column component with customisable gap. Consumers should define their own grid classes in global styles.

```svelte
<Column class="col-12 col-tab-6 col-desk-4">
	Content
</Column>
```

#### Customisation

```css
/* Global column gap */
:root {
	--col-gap: 1.5rem;
	--col-gap-mobile: 0.75rem; /* Optional: mobile-specific gap (< 568px) */
}

/* Define your own grid classes */
.col-12 { width: 100%; }
.col-6 { width: 50%; }

@media (min-width: 768px) {
	.col-tab-6 { width: 50%; }
}
```

## Navigation Components

### Sidebar

Mobile-first sidebar navigation component with toggle and hamburger display modes. Automatically collapses on mobile and remains visible on desktop.

```svelte
<script>
	import { Sidebar } from '@ainsleydev/sveltekit-helper/components'
</script>

<Sidebar bind:isOpen>
	<nav>
		<a href="/">Home</a>
		<a href="/about">About</a>
		<a href="/contact">Contact</a>
	</nav>
</Sidebar>
```

#### Props

- `menuLabel?: string` - Label for toggle button (default: 'Menu')
- `isOpen?: boolean` - Bindable open/closed state
- `position?: 'left' | 'right'` - Sidebar position (default: 'left')
- `width?: string` - Sidebar width on mobile (default: '50vw')
- `top?: number` - Sticky position offset on desktop (default: 160)
- `closeOnOverlayClick?: boolean` - Close when overlay is clicked (default: true)
- `overlayOpacity?: number` - Overlay opacity when open (default: 0.3)
- `toggleStyle?: 'toggle' | 'hamburger'` - Toggle display mode (default: 'toggle')
- `class?: string` - Additional CSS classes
- `onOpen?: () => void` - Callback when sidebar opens
- `onClose?: () => void` - Callback when sidebar closes
- `onToggle?: (isOpen: boolean) => void` - Callback when sidebar toggles

#### Examples

With hamburger menu:

```svelte
<Sidebar toggleStyle="hamburger" bind:isOpen>
	<nav>...</nav>
</Sidebar>
```

Right-side with custom width:

```svelte
<Sidebar position="right" width="300px">
	<nav>...</nav>
</Sidebar>
```

#### Customisation

Override CSS variables globally from `:root`:

```css
:root {
	--sidebar-width: 400px;
	--sidebar-min-width: 300px;
	--sidebar-bg: #1a1a1a;
	--sidebar-border-colour: rgba(255, 255, 255, 0.2);
	--sidebar-overlay-colour: #000;
	--sidebar-overlay-opacity: 0.5;

	/* Toggle button */
	--sidebar-toggle-bg: #2a2a2a;
	--sidebar-toggle-colour: #fff;
	--sidebar-toggle-padding: 0.5rem 1.5rem;
	--sidebar-toggle-radius: 8px;
	--sidebar-toggle-font-size: 1rem;

	/* Inner spacing */
	--sidebar-inner-padding: 2rem 2rem 0 2rem;
}
```

Or use inline styles:

```svelte
<Sidebar style="--sidebar-bg: #2a2a2a; --sidebar-width: 400px">
	<nav>...</nav>
</Sidebar>
```

### Hamburger

Hamburger menu icon with animation for mobile navigation. Uses `svelte-hamburgers` under the hood.

```svelte
<script>
	import { Hamburger } from '@ainsleydev/sveltekit-helper/components'

	let isOpen = $state(false)
</script>

<Hamburger bind:isOpen />
```

#### Props

- `isOpen?: boolean` - Bindable open/closed state
- `gap?: string` - Distance from top/right edges (default: '0.8rem')
- `class?: string` - Additional CSS classes
- `ariaLabel?: string` - Accessibility label (default: 'Toggle menu')
- `onChange?: (isOpen: boolean) => void` - Callback when state changes

#### Customisation

```css
:root {
	--hamburger-gap: 1rem;
	--hamburger-z-index: 10000;
	--hamburger-colour: #fff;
	--hamburger-layer-width: 28px;
	--hamburger-layer-height: 3px;
	--hamburger-layer-spacing: 6px;
	--hamburger-border-radius: 3px;
}
```

## Form Utilities

### generateFormSchema

Generates a Zod schema from Payload CMS form fields.

```typescript
import { generateFormSchema } from '@ainsleydev/sveltekit-helper/utils/forms'

const fields = [
	{ blockType: 'text', name: 'name', label: 'Name', required: true },
	{ blockType: 'email', name: 'email', label: 'Email', required: true },
	{ blockType: 'textarea', name: 'message', label: 'Message', required: false }
]

const schema = generateFormSchema(fields)
// Returns Zod schema with appropriate validation
```

### flattenZodErrors

Converts Zod validation errors into a simple key-value object.

```typescript
import { flattenZodErrors } from '@ainsleydev/sveltekit-helper/utils/forms'
import { z } from 'zod'

const schema = z.object({ email: z.string().email() })
const result = schema.safeParse({ email: 'invalid' })

if (!result.success) {
	const errors = flattenZodErrors(result.error)
	// { email: 'Invalid email' }
}
```

## Payload CMS Components

### PayloadForm

Renders a form dynamically from Payload CMS form builder fields.

```svelte
<script>
	import { PayloadForm } from '@ainsleydev/sveltekit-helper/components/payload'

	export let data
</script>

<PayloadForm
	form={data.form}
	apiEndpoint="/api/forms"
/>
```

#### Custom Submission

```svelte
<PayloadForm
	form={data.form}
	onSubmit={async (formData) => {
		// Custom submission logic
		await customAPI.submit(formData)
	}}
/>
```

#### Customisation

Override CSS variables globally:

```css
/* Global form styling */
:root {
	--form-gap: 1.5rem;
	--form-input-padding: 1rem;
	--form-input-border: 1px solid #e5e7eb;
	--form-input-border-radius: 0.5rem;
	--form-input-bg: #ffffff;
	--form-input-colour: #111827;
	--form-error-colour: #ef4444;
	--form-error-bg: #fee2e2;
	--form-success-colour: #10b981;
	--form-success-bg: #d1fae5;
	--form-button-bg: #3b82f6;
	--form-button-colour: #ffffff;
	--form-button-hover-bg: #2563eb;
	--form-button-disabled-bg: #9ca3af;
}
```

### PayloadMedia

Renders responsive images and videos from Payload CMS media fields with automatic format prioritisation (AVIF → WebP → JPEG/PNG).

```svelte
<script>
	import { PayloadMedia } from '@ainsleydev/sveltekit-helper/components/payload'

	export let data
</script>

<PayloadMedia
	data={data.image}
	loading="lazy"
	maxWidth={1200}
/>
```

#### Props

- `data`: Payload media object with `url`, `sizes`, `mimeType`, etc.
- `loading`: Optional `'lazy'` or `'eager'` loading strategy
- `maxWidth`: Optional maximum width to limit responsive sources
- `breakpointBuffer`: Pixels to add to breakpoint media queries (default: 50)
- `className`: Optional CSS class name
- `onload`: Optional load event handler

## Peer Dependencies

- `svelte@^5.0.0`
- `@sveltejs/kit@^2.0.0`

### Optional Dependencies

- `payload@^3.0.0` - For Payload CMS components
- `zod@^3.0.0` - For form validation

## Development

```bash
# Install dependencies
pnpm install

# Build the package
pnpm build

# Run tests
pnpm test

# Lint and format
pnpm lint
pnpm format
```

## License

MIT
