# Email Preview Templates

Example preview templates showing how to preview payload-helper emails with your branding.

## Usage in Your Project

Copy these files to your CMS project's `preview-emails/` directory and customize with your branding.

### 1. Copy the templates

Create `preview-emails/ForgotPassword.tsx` and `preview-emails/VerifyAccount.tsx` (see examples in this directory).

### 2. Customize the theme

Update the `theme` object with your actual branding from your `payloadHelper()` config.

### 3. Run the preview

```bash
npx email-templates preview ./preview-emails
```

Open http://localhost:3000 to view your emails.

## Tips

- Use the same theme configuration from your `payloadHelper()` config
- Test different user names and content variations
- Check on mobile by accessing from your phone on the same network
