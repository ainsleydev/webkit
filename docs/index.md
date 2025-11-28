---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "WebKit"
  text: "Streamline your web project lifecycle"
  tagline: One manifest file. Complete infrastructure. Automatic CI/CD.
  actions:
    - theme: brand
      text: Get started
      link: /getting-started/installation
    - theme: alt
      text: View on GitHub
      link: https://github.com/ainsleydev/webkit

features:
  - title: Single source of truth
    details: Define your entire project in one app.json file. Apps, infrastructure, environments, and monitoring - all in one place.
  - title: Automatic CI/CD
    details: WebKit generates GitHub Actions workflows for testing, building, and deploying. Push to main and watch it deploy.
  - title: Infrastructure as code
    details: Terraform configurations generated automatically. Support for DigitalOcean, Hetzner, Backblaze B2, and Turso.
  - title: Secret management
    details: Built-in SOPS integration with Age encryption. Manage secrets securely across environments.
  - title: Uptime monitoring
    details: Configure HTTP, DNS, and custom monitors. Status pages included with Peekaping integration.
  - title: Drift detection
    details: Track generated files and detect manual modifications. Always know what's changed.
---

## Quick start

Install WebKit and generate your first project in minutes:

```bash
# Install WebKit
curl -sSL https://raw.githubusercontent.com/ainsleydev/webkit/main/bin/install.sh | sh

# Create your manifest
cat > app.json << 'EOF'
{
  "$schema": "https://raw.githubusercontent.com/ainsleydev/webkit/main/schema.json",
  "project": {
    "name": "my-site",
    "title": "My Site",
    "repo": "github.com/username/my-site"
  },
  "apps": [
    {
      "name": "web",
      "type": "svelte-kit",
      "path": "./apps/web"
    }
  ]
}
EOF

# Generate project files
webkit update
```

[Read the full guide →](/getting-started/quick-start)
