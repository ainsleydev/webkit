
# Build Process

## Image Building

All Docker images are built in GitHub Actions CI. WebKit generates
workflows that:

1. Build on push to main/staging branches
2. Tag images with git commit SHA: `ghcr.io/{org}/{app}:{git-sha}`
3. Push to GitHub Container Registry
4. Deploy to configured infrastructure

### Image Naming Convention

**Format:** `ghcr.io/{project.name}/{app.name}:{tag}`

**Example:**
- Project: `my-website`
- App: `cms`
- Commit: `abc123f`
- **Result:** `ghcr.io/my-website/cms:abc123f`

### Monorepo Build Context

For monorepos, the build context is the app's `path`:
```json
{
  "name": "cms",
  "path": "services/cms"  // Docker build context
}