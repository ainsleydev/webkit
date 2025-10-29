# @ainsleydev/email-templates

Composable email template building blocks with theme system for React Email.

## Features

- **Theme system** - Customise colours, branding, and styles
- **BaseEmail component** - Reusable layout with logo and footer
- **Runtime rendering** - No build step required
- **Type-safe** - Full TypeScript support
- **Flexible** - Create your own templates using React Email components
- **Lightweight** - Minimal dependencies

## Installation

```bash
pnpm add @ainsleydev/email-templates
```

## Quick start

### 1. Create your email template

```typescript
// emails/ForgotPassword.tsx
import { BaseEmail, Heading, Text, Section, Button, generateStyles } from '@ainsleydev/email-templates'
import type { EmailTheme } from '@ainsleydev/email-templates'

interface ForgotPasswordProps {
  theme: EmailTheme
  user: { firstName: string }
  resetUrl: string
}

export const ForgotPasswordEmail = ({ theme, user, resetUrl }: ForgotPasswordProps) => {
  const styles = generateStyles(theme)

  return (
    <BaseEmail theme={theme} previewText="Reset your password">
      <Heading style={styles.heading}>
        Hello, {user.firstName}!
      </Heading>
      <Text style={styles.text}>
        We received a request to reset your password. Click the button below to continue.
      </Text>
      <Section style={{ textAlign: 'center' }}>
        <Button href={resetUrl} style={styles.button}>
          Reset Password
        </Button>
      </Section>
    </BaseEmail>
  )
}
```

### 2. Render the template

```typescript
import { renderEmail } from '@ainsleydev/email-templates'
import { ForgotPasswordEmail } from './emails/ForgotPassword'

const html = await renderEmail({
  component: ForgotPasswordEmail,
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

## Theme customisation

### Using partial overrides

```typescript
const html = await renderEmail({
  component: ForgotPasswordEmail,
  props: {
    user: { firstName: 'Jane' },
    resetUrl: 'https://example.com/reset/abc',
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

### Creating a complete theme

```typescript
import { defaultTheme, type EmailTheme } from '@ainsleydev/email-templates'

const customTheme: EmailTheme = {
  ...defaultTheme,
  branding: {
    companyName: 'Custom Corp',
    logoUrl: 'https://custom.com/logo.png',
    logoWidth: 200,
    footerText: 'All rights reserved.',
    websiteUrl: 'https://custom.com',
  },
  colours: {
    ...defaultTheme.colours,
    text: {
      ...defaultTheme.colours.text,
      action: '#ff0000',
    },
  },
}

// Use in all your emails.
const html = await renderEmail({
  component: MyEmail,
  props: myProps,
  theme: customTheme,
})
```

## Usage with Payload CMS

Payload CMS allows you to customise email templates for authentication emails.

```typescript
import { buildConfig } from 'payload'
import { renderEmail } from '@ainsleydev/email-templates'
import { ForgotPasswordEmail } from './emails/ForgotPassword'
import { VerifyAccountEmail } from './emails/VerifyAccount'

export default buildConfig({
  email: {
    fromName: 'My App',
    fromAddress: 'noreply@myapp.com',
    // Configure your email transport.
  },
  collections: [
    {
      slug: 'users',
      auth: {
        forgotPassword: {
          generateEmailHTML: async ({ token, user }) => {
            return renderEmail({
              component: ForgotPasswordEmail,
              props: {
                user: { firstName: user.firstName },
                resetUrl: `${process.env.FRONTEND_URL}/reset-password?token=${token}`,
              },
              theme: {
                branding: {
                  companyName: 'My App',
                  logoUrl: `${process.env.FRONTEND_URL}/logo.png`,
                },
              },
            })
          },
        },
        verify: {
          generateEmailHTML: async ({ token, user }) => {
            return renderEmail({
              component: VerifyAccountEmail,
              props: {
                user: { firstName: user.firstName },
                verifyUrl: `${process.env.FRONTEND_URL}/verify?token=${token}`,
              },
            })
          },
        },
      },
    },
  ],
})
```

## React Email components

All React Email components are re-exported for convenience:

```typescript
import {
  Html, Head, Preview, Body, Container, Section, Row, Column,
  Heading, Text, Button, Link, Img, Hr,
  // ...and more
} from '@ainsleydev/email-templates'

// Use them directly in your templates:
<Heading>Welcome</Heading>
<Text>Hello world</Text>
<Button href="...">Click here</Button>
```

See [React Email documentation](https://react.email/docs/components/html) for full component API.

## API

### `renderEmail(options)`

Renders an email template component to HTML string.

**Options:**
- `component` - Your React component that accepts `theme` prop
- `props` - Props to pass to your component (excluding theme)
- `theme` - Optional theme overrides
- `plainText` - Render as plain text instead of HTML (default: `false`)

**Returns:** `Promise<string>` - HTML or plain text string

### `BaseEmail`

Base layout component with logo, content area, and footer.

**Props:**
- `theme: EmailTheme` - Theme configuration
- `previewText?: string` - Email preview text
- `children: React.ReactNode` - Email content

### `generateStyles(theme)`

Generates common style objects from theme configuration.

**Returns:** Object with `heading`, `text`, `button`, `hr`, etc. styles

### `defaultTheme`

Default theme configuration based on ainsley.dev design system.

### `mergeTheme(partial)`

Merges partial theme with default theme.

## Examples

See the [test file](src/renderer.test.ts) for more examples of creating custom templates.

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

## TODO

- **Go CLI support** - Add CLI command for rendering templates from Go via `exec.Command`. This would allow Go applications to use the same email templates without needing a Node.js runtime dependency.

## Licence

MIT © [ainsley.dev](https://ainsley.dev)
