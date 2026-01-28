# @ainsleydev/email-templates

## 0.0.4

### Patch Changes

- 3ea3383: Fix invalid peer dependency version range for @react-email/components and add React 19 support

## 0.0.3

### Patch Changes

- 8ccfa4d: Added CLI preview command that allows consumers to spin up a preview server for their custom email templates. The command automatically discovers `.tsx` and `.jsx` files in a specified directory and serves them with the default theme.

  Usage:
  - `npx @ainsleydev/email-templates preview` - Preview templates in current directory
  - `npx @ainsleydev/email-templates preview ./emails` - Preview templates in specific directory
  - `npx @ainsleydev/email-templates preview ./emails --port=3001` - Use custom port

## 0.0.2

### Patch Changes

- e5aa225: Initial Release

## 0.0.1

### Patch Changes

- Initial release with email template building blocks.
- Theme system with customisable colours and branding.
- BaseEmail component for consistent layout.
- Generic renderEmail() function accepting any React component.
- Direct re-export of React Email components (no namespace).
- Full TypeScript support with comprehensive types.
- Generic colour naming (dark, darker, highlight, accent vs company-specific names).
