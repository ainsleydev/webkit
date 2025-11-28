# Your first project

In this tutorial, you'll build and deploy a SvelteKit portfolio site using WebKit. By the end, you'll have:

- A fully configured monorepo structure
- Automated CI/CD pipelines with GitHub Actions
- Infrastructure deployed to DigitalOcean App Platform
- Environment variable management

## Prerequisites

Before starting, ensure you have:

- [WebKit installed](/getting-started/installation)
- A GitHub account and repository
- A DigitalOcean account (for deployment)
- Node.js 18+ and pnpm installed

## Step 1: Create your project

Start with a new directory and initialise your SvelteKit app:

```bash
mkdir portfolio && cd portfolio
mkdir -p apps/web
cd apps/web
pnpm create svelte@latest .
cd ../..
```

Choose your preferred SvelteKit options. For a portfolio site, we recommend:
- Skeleton project
- TypeScript
- ESLint and Prettier

## Step 2: Create your manifest

Create `app.json` in the project root:

```json
{
  "$schema": "https://raw.githubusercontent.com/ainsleydev/webkit/main/schema.json",
  "project": {
    "name": "portfolio",
    "title": "My Portfolio",
    "description": "Personal portfolio website",
    "repo": "github.com/yourusername/portfolio"
  },
  "apps": [
    {
      "name": "web",
      "type": "svelte-kit",
      "path": "./apps/web",
      "build": {
        "dockerfile": true,
        "port": 3000
      },
      "infrastructure": {
        "provider": "digital_ocean",
        "type": "app"
      },
      "domains": {
        "primary": "portfolio.example.com"
      }
    }
  ]
}
```

This configuration tells WebKit:
- Your app is a SvelteKit project in `./apps/web`
- It should generate a Dockerfile exposing port 3000
- Deploy to DigitalOcean App Platform
- Use `portfolio.example.com` as the primary domain

## Step 3: Generate project files

Run WebKit to generate all configuration:

```bash
webkit update
```

WebKit creates:
- `.github/workflows/` - CI/CD pipelines
- `.github/actions/` - Reusable workflow components
- `package.json` - Root package with workspace scripts
- `pnpm-workspace.yaml` - Workspace configuration
- `turbo.json` - Turborepo build configuration
- `.gitignore`, `.editorconfig` - Project settings

## Step 4: Set up environment variables

Create environment-specific configuration in your manifest:

```json
{
  "apps": [
    {
      "name": "web",
      "environment": {
        "production": {
          "PUBLIC_SITE_URL": {
            "source": "value",
            "value": "https://portfolio.example.com"
          }
        },
        "staging": {
          "PUBLIC_SITE_URL": {
            "source": "value",
            "value": "https://staging.portfolio.example.com"
          }
        }
      }
    }
  ]
}
```

Generate environment files:

```bash
webkit env scaffold
```

This creates `.env.example` files that you can copy and customise locally.

## Step 5: Configure GitHub secrets

For CI/CD to work, add these secrets to your GitHub repository:

| Secret | Description |
|--------|-------------|
| `DIGITALOCEAN_ACCESS_TOKEN` | DigitalOcean API token |
| `TF_STATE_BUCKET` | S3-compatible bucket for Terraform state |
| `TF_STATE_ACCESS_KEY` | Access key for state bucket |
| `TF_STATE_SECRET_KEY` | Secret key for state bucket |

Navigate to your repository's Settings → Secrets and variables → Actions to add these.

## Step 6: Deploy infrastructure

First, generate Terraform configuration:

```bash
webkit infra plan
```

Review the planned changes, then apply:

```bash
webkit infra apply
```

This provisions:
- DigitalOcean App Platform application
- Domain configuration
- Any additional resources defined in your manifest

## Step 7: Push and deploy

Commit your changes and push to GitHub:

```bash
git add .
git commit -m "feat: Initial WebKit setup"
git push origin main
```

GitHub Actions automatically:
1. Runs linting and tests
2. Builds your SvelteKit app
3. Deploys to DigitalOcean App Platform

Monitor the deployment in your repository's Actions tab.

## Adding monitoring (optional)

WebKit can configure uptime monitoring for your site. Add monitoring configuration:

```json
{
  "monitoring": {
    "status_page": {
      "name": "Portfolio Status",
      "slug": "portfolio-status"
    }
  },
  "apps": [
    {
      "name": "web",
      "monitoring": {
        "http": true
      }
    }
  ]
}
```

Run `webkit update` to generate monitoring resources, then `webkit infra apply` to provision them.

## Project structure

Your final project structure:

```
portfolio/
├── .github/
│   ├── actions/
│   │   ├── setup/
│   │   └── notify/
│   └── workflows/
│       ├── pr.yaml
│       └── release.yaml
├── .webkit/
│   └── manifest.json
├── apps/
│   └── web/
│       ├── src/
│       ├── package.json
│       └── svelte.config.js
├── app.json
├── package.json
├── pnpm-workspace.yaml
└── turbo.json
```

## Next steps

You've successfully deployed your first WebKit project! Continue learning:

- [Core concepts](/getting-started/core-concepts) - Understand WebKit's architecture
- [Apps configuration](/manifest/apps) - Advanced app settings
- [Resources](/manifest/resources) - Add databases and storage
- [Infrastructure providers](/infrastructure/providers/digital-ocean) - DigitalOcean configuration options
