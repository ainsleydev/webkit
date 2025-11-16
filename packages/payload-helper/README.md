## Payload Helper

Payload Helper is a project designed to extend and enhance functionality for Payload CMS. This
project includes custom configurations, scripts, and utilities to streamline development and content
management processes.

## Installation

```bash
npm install @ainsleydev/payload-helper
# or
pnpm add @ainsleydev/payload-helper
# or
yarn add @ainsleydev/payload-helper
```

## Usage

Add the plugin to your Payload configuration:

```typescript
import { payloadHelper } from '@ainsleydev/payload-helper'

export default buildConfig({
	plugins: [
		payloadHelper({
			siteName: 'My Site',
			admin: {
				logo: {
					path: '/images/logo-light.svg',
					darkModePath: '/images/logo-dark.svg',
					width: 150,
					height: 40,
					alt: 'My Site Logo',
					className: 'custom-logo-class',
				},
				icon: {
					path: '/images/icon-light.svg',
					darkModePath: '/images/icon-dark.svg',
					width: 120,
					height: 120,
					alt: 'My Site Icon',
					className: 'custom-icon-class',
				},
			},
		}),
	],
})
```

## Available Exports

The payload-helper package exports a variety of utilities, collections, globals, and components that you can use directly in your Payload projects.

### Collections

Pre-built Payload collections ready to use:

```typescript
import { Media, Redirects, imageSizes, imageSizesWithAvif } from '@ainsleydev/payload-helper'

// Use in your Payload config
export default buildConfig({
	collections: [
		Media({ includeAvif: true, additionalFields: [...] }),
		Redirects(),
	],
})
```

**Subpath imports** (optional):
```typescript
import { Media, Redirects } from '@ainsleydev/payload-helper/collections'
```

### Globals

Pre-configured global types for common website needs:

```typescript
import { Settings, Navigation, countries, languages } from '@ainsleydev/payload-helper'

export default buildConfig({
	globals: [
		Settings({ additionalTabs: [...] }),
		Navigation({ includeFooter: true }),
	],
})
```

**Subpath imports** (optional):
```typescript
import { Settings, Navigation } from '@ainsleydev/payload-helper/globals'
```

### Utilities

Helpful utility functions for validation, field operations, and Lexical conversion:

```typescript
import {
	env,
	fieldHasName,
	validateURL,
	validatePostcode,
	htmlToLexical,
	lexicalToHtml,
} from '@ainsleydev/payload-helper'

// Use in field validation
{
	name: 'website',
	type: 'text',
	validate: validateURL,
}

// Convert HTML to Lexical format
const lexicalData = await htmlToLexical('<p>Hello world</p>')

// Convert Lexical to HTML
const html = await lexicalToHtml(lexicalData)
```

**Subpath imports** (optional):
```typescript
import { validateURL, htmlToLexical } from '@ainsleydev/payload-helper/util'
```

### Common/Reusable Fields

Reusable field definitions:

```typescript
import { SEOFields } from '@ainsleydev/payload-helper'
import { seoPlugin } from '@payloadcms/plugin-seo'

export default buildConfig({
	plugins: [
		seoPlugin({
			collections: ['pages'],
			fields: SEOFields,
		}),
	],
})
```

**Subpath imports** (optional):
```typescript
import { SEOFields } from '@ainsleydev/payload-helper/common'
```

### Endpoints

Custom API endpoints:

```typescript
import { findBySlug } from '@ainsleydev/payload-helper'

// Use in your collection config
export const Pages: CollectionConfig = {
	slug: 'pages',
	endpoints: [findBySlug],
}
```

**Subpath imports** (optional):
```typescript
import { findBySlug } from '@ainsleydev/payload-helper/endpoints'
```

### Email Components

Customizable email templates for authentication flows:

```typescript
import { ForgotPasswordEmail, VerifyAccountEmail } from '@ainsleydev/payload-helper'
import type { ForgotPasswordEmailProps, VerifyAccountEmailProps } from '@ainsleydev/payload-helper'

// Use directly in custom email handlers
const emailHtml = ForgotPasswordEmail({
	resetPasswordToken: 'token123',
	frontEndUrl: 'https://yoursite.com',
	// ...other props
})
```

### Schema Utilities

For projects that need JSON schema generation (e.g., Go type generation):

```typescript
import { fieldMapper, schemas, addGoJSONSchema } from '@ainsleydev/payload-helper'
import type { SchemaOptions } from '@ainsleydev/payload-helper'

// Apply Go schema mappings
const config = fieldMapper(sanitizedConfig, {
	useWebKitMedia: true,
	assignRelationships: true,
})
```

## Configuration

### Admin configuration

The `admin` object contains configuration for admin UI customisation:

#### Logo configuration

The logo appears in the navigation area of the Payload admin dashboard.

- `path` (required): Path to the logo image file
- `darkModePath` (optional): Path to the logo for dark mode
- `width` (optional): Logo width in pixels (default: 150)
- `height` (optional): Logo height in pixels (default: 40)
- `alt` (optional): Alt text for the logo (defaults to siteName)
- `className` (optional): Custom CSS class name

#### Icon configuration

The icon appears in the top left corner of the Payload admin dashboard.

- `path` (required): Path to the icon image file
- `darkModePath` (optional): Path to the icon for dark mode
- `width` (optional): Icon width in pixels (default: 120)
- `height` (optional): Icon height in pixels (default: 120)
- `alt` (optional): Alt text for the icon (defaults to siteName)
- `className` (optional): Custom CSS class name

### Web server configuration

Configure cache invalidation for your web server:

```typescript
payloadHelper({
	siteName: 'My Site',
	webServer: {
		apiKey: 'your-api-key',
		baseURL: 'https://your-site.com',
		cacheEndpoint: '/api/cache/purge',
	},
})
```

### Email configuration

Configure branded email templates for Payload authentication flows. Automatically applies to all collections with `auth` enabled.

```typescript
payloadHelper({
	siteName: 'My Site',
	email: {
		frontEndUrl: 'https://your-site.com', // Optional, defaults to Payload's serverURL
		theme: {
			branding: {
				companyName: 'My Company',
				logoUrl: 'https://your-site.com/logo.png',
			},
			colours: {
				background: {
					accent: '#ff5043',
				},
			},
		},
		forgotPassword: {
			heading: 'Reset your password',
			bodyText: 'Click the button below to reset your password.',
			buttonText: 'Reset Password',
		},
		verifyAccount: {
			heading: 'Welcome aboard',
			bodyText: 'Please verify your email address.',
			buttonText: 'Verify Email',
		},
	},
})
```

#### Previewing emails

Preview your emails with your actual branding directly from your Payload configuration:

```bash
npx payload-helper preview-emails
```

This command will:
- Read your `payload.config.ts` to extract your email theme configuration
- Generate preview templates with your actual branding
- Launch a preview server at http://localhost:3000

You can optionally specify a custom port:

```bash
npx payload-helper preview-emails --port 3001
```

The preview will show both ForgotPassword and VerifyAccount emails using your configured theme, frontEndUrl, and content overrides.

## Utilities

### Environment Variables

The package exports an `env` utility object that provides access to environment variables required by Payload Helper.

```typescript
import { env } from '@ainsleydev/payload-helper'

// Access environment variables
const spaces = {
	key: env.DO_SPACES_KEY,
	secret: env.DO_SPACES_SECRET,
	endpoint: env.DO_SPACES_ENDPOINT,
	region: env.DO_SPACES_REGION,
	bucket: env.DO_SPACES_BUCKET,
}
```

**Note:** When using ESM modules, always import utilities from the main package entry point as shown above. Direct imports to subpaths (e.g., `@ainsleydev/payload-helper/dist/util/env`) are not supported.

## Open Source

ainsley.dev permits the use of any HTML, SCSS and Javascript found within the repository for use
with external projects.

## Trademark

ainsley.dev and the ainsley.dev logo are either registered trademarks or trademarks of ainsley.dev
LTD in the United Kingdom and/or other countries. All other trademarks are the property of their
respective owners.
