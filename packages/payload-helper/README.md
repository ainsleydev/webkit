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

## Open Source

ainsley.dev permits the use of any HTML, SCSS and Javascript found within the repository for use
with external projects.

## Trademark

ainsley.dev and the ainsley.dev logo are either registered trademarks or trademarks of ainsley.dev
LTD in the United Kingdom and/or other countries. All other trademarks are the property of their
respective owners.
