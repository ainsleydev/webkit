# @ainsleydev/sveltekit-helper

SvelteKit utilities, components and helpers for ainsley.dev builds.

## Installation

```bash
pnpm add @ainsleydev/sveltekit-helper
```

## Features

- **Grid System**: Responsive Container, Row, and Column components with CSS variables
- **Form Utilities**: Client-side form management with Zod validation
- **Payload CMS Integration**: Ready-to-use components for Payload CMS forms and media

## Grid Components

### Container

Center content horizontally with predefined max-width and support for breakout layouts.

```svelte
<script>
	import { Container, Row, Column } from '@ainsleydev/sveltekit-helper/components/Grid'
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

Override CSS variables to customise the container:

```css
.container {
	--container-padding: 2rem;
	--container-max-width: 1400px;
	--container-breakout-max-width: 1600px;
}
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
.row {
	--row-gap: 1.5rem;
}
```

### Column

Responsive column with 12-column grid system.

```svelte
<!-- Full width on mobile, half on tablet, third on desktop -->
<Column class="col-12 col-tab-6 col-desk-4">
	Content
</Column>

<!-- With offset -->
<Column class="col-8 offset-desk-2">
	Centred content
</Column>
```

#### Customisation

```css
[class*='col-'] {
	--col-gap: 1.5rem;
}
```

#### Responsive Classes

- **Base**: `col-1` to `col-12`, `col-auto`
- **Tablet (768px+)**: `col-tab-1` to `col-tab-12`, `col-tab-auto`
- **Desktop (1024px+)**: `col-desk-1` to `col-desk-12`, `col-desk-auto`
- **Offsets**: `offset-tab-1` to `offset-tab-11`, `offset-desk-1` to `offset-desk-11`

## Form Utilities

### clientForm

Creates a reactive client-side form store with validation and submission handling.

```svelte
<script>
	import { z } from 'zod'
	import { clientForm } from '@ainsleydev/sveltekit-helper/utils/forms'

	const schema = z.object({
		email: z.string().email(),
		password: z.string().min(8)
	})

	const { fields, errors, validate, submitting, enhance } = clientForm(
		schema,
		{ submissionDelay: 300 },
		async (data) => {
			const response = await fetch('/api/login', {
				method: 'POST',
				body: JSON.stringify(data)
			})
		}
	)
</script>

<form use:enhance method="POST">
	<input
		type="email"
		name="email"
		bind:value={$fields.email}
		on:blur={() => validate({ field: 'email' })}
	/>
	{#if $errors.email}
		<span class="error">{$errors.email}</span>
	{/if}

	<button type="submit" disabled={$submitting}>
		{$submitting ? 'Submitting...' : 'Submit'}
	</button>
</form>
```

### serverForm

Validates form data on the server using Zod schema.

```typescript
// src/routes/login/+page.server.ts
import { serverForm } from '@ainsleydev/sveltekit-helper/utils/forms'
import { z } from 'zod'

export const actions = {
	default: async ({ request }) => {
		const schema = z.object({
			email: z.string().email(),
			password: z.string().min(8)
		})

		const { valid, data, errors } = await serverForm(request, schema)

		if (!valid) {
			return { errors }
		}

		// Process form...
	}
}
```

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

Style the form using CSS variables:

```css
.payload-form {
	--form-gap: 1.5rem;
	--form-input-padding: 1rem;
	--form-input-border: 1px solid #e5e7eb;
	--form-input-border-radius: 0.5rem;
	--form-input-bg: #ffffff;
	--form-input-text: #111827;
	--form-error-color: #ef4444;
	--form-success-color: #10b981;
	--form-button-bg: #3b82f6;
	--form-button-text: #ffffff;
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
