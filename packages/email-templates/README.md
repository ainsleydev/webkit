# @ainsleydev/email-templates

Composable, reusable email templates built with React Email for JavaScript and Go projects.

## Features

- 🎨 **Theme system** - Customise colours, branding, and styles
- 🔄 **Runtime rendering** - No build step required
- 🎯 **Type-safe** - Full TypeScript support
- 🧩 **Composable** - Partial theme overrides with defaults
- 🌐 **Cross-language** - JavaScript API + CLI for Go integration
- 📧 **Production-ready** - Based on React Email best practices

## Installation

```bash
pnpm add @ainsleydev/email-templates
```

## Quick start

```typescript
import { renderEmail } from '@ainsleydev/email-templates'

const html = await renderEmail({
	template: 'forgot-password',
	props: {
		user: { firstName: 'John' },
		resetUrl: 'https://example.com/reset/token123',
	},
})

// Send via your email service.
await emailService.send({
	to: 'user@example.com',
	subject: 'Reset your password',
	html,
})
```

## Usage with Payload CMS

```typescript
import { renderEmail } from '@ainsleydev/email-templates'

export default buildConfig({
	email: {
		fromName: 'My App',
		fromAddress: 'noreply@myapp.com',
		transportOptions: { /* ... */ },
	},
	// Custom email handler.
	onInit: async (payload) => {
		payload.email = {
			...payload.email,
			sendEmail: async (message) => {
				const html = await renderEmail({
					template: 'forgot-password',
					props: {
						user: { firstName: message.to },
						resetUrl: message.data.resetUrl,
					},
					theme: {
						branding: {
							companyName: 'My Company',
							logoUrl: 'https://mycompany.com/logo.png',
						},
					},
				})
				// Send email with your provider.
			},
		}
	},
})
```

## Theme customisation

### Partial overrides

```typescript
const html = await renderEmail({
	template: 'verify-account',
	props: {
		user: { firstName: 'Jane' },
		verifyUrl: 'https://example.com/verify/abc',
	},
	theme: {
		branding: {
			companyName: 'My Company',
			logoUrl: 'https://mycompany.com/logo.png',
			logoWidth: 150,
			footerText: 'All rights reserved.',
			websiteUrl: 'https://mycompany.com',
		},
		colours: {
			text: {
				action: '#0066cc', // Custom link colour.
			},
		},
	},
})
```

### Complete theme

```typescript
import { defaultTheme, type EmailTheme } from '@ainsleydev/email-templates'

const customTheme: EmailTheme = {
	...defaultTheme,
	branding: {
		companyName: 'Custom Corp',
		logoUrl: 'https://custom.com/logo.png',
		logoWidth: 200,
	},
}
```

## Available templates

### forgot-password

Password reset email with call-to-action button.

```typescript
await renderEmail({
	template: 'forgot-password',
	props: {
		user: { firstName: 'John' },
		resetUrl: 'https://example.com/reset/token',
	},
})
```

### verify-account

Account verification email for new users.

```typescript
await renderEmail({
	template: 'verify-account',
	props: {
		user: { firstName: 'Jane' },
		verifyUrl: 'https://example.com/verify/token',
	},
})
```

## CLI usage (Go integration)

The package includes a CLI for rendering templates from Go or other languages.

### Command line

```bash
npx @ainsleydev/email-templates render \
  --template forgot-password \
  --props '{"user":{"firstName":"John"},"resetUrl":"https://example.com/reset"}' \
  --theme '{"branding":{"companyName":"My Company"}}'
```

### Go integration

```go
package main

import (
	"os/exec"
	"encoding/json"
)

func renderEmail() (string, error) {
	props := map[string]interface{}{
		"user": map[string]string{"firstName": "John"},
		"resetUrl": "https://example.com/reset/token123",
	}
	propsJSON, _ := json.Marshal(props)

	theme := map[string]interface{}{
		"branding": map[string]interface{}{
			"companyName": "My Company",
		},
	}
	themeJSON, _ := json.Marshal(theme)

	cmd := exec.Command("npx", "@ainsleydev/email-templates", "render",
		"--template", "forgot-password",
		"--props", string(propsJSON),
		"--theme", string(themeJSON))

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
```

## Development

```bash
# Install dependencies.
pnpm install

# Run tests.
pnpm test

# Build package.
pnpm build

# Format code.
pnpm format

# Lint code.
pnpm lint
```

## Licence

MIT © [ainsley.dev](https://ainsley.dev)
