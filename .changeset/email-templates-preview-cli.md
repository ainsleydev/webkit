---
"@ainsleydev/email-templates": patch
---

Added CLI preview command that allows consumers to spin up a preview server for their custom email templates. The command automatically discovers `.tsx` and `.jsx` files in a specified directory and serves them with the default theme.

Usage:
- `npx @ainsleydev/email-templates preview` - Preview templates in current directory
- `npx @ainsleydev/email-templates preview ./emails` - Preview templates in specific directory
- `npx @ainsleydev/email-templates preview ./emails --port=3001` - Use custom port
