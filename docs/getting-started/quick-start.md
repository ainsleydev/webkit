# Quick start

This guide will walk you through creating your first WebKit project in under 5 minutes. By the end, you'll have a working `app.json` manifest and understand the basic workflow.

## Prerequisites

Before you begin, make sure you have:

- WebKit installed ([installation guide](/getting-started/installation))
- Git initialised in your project directory
- A GitHub repository (optional, but recommended)

## Create your first manifest

1. **Navigate to your project directory**:
   ```bash
   mkdir my-website
   cd my-website
   git init
   ```

2. **Create an `app.json` file**:
   ```bash
   cat > app.json << 'MANIFEST'
   {
     "project": {
       "name": "my-website",
       "title": "My Website",
       "description": "A modern full-stack web application",
       "repo": "https://github.com/yourusername/my-website"
     },
     "resources": [
       {
         "name": "db",
         "type": "postgres",
         "provider": "digitalocean",
         "config": {
           "size": "db-s-1vcpu-1gb",
           "region": "lon1"
         }
       }
     ],
     "apps": [
       {
         "name": "web",
         "type": "sveltekit",
         "path": "apps/web",
         "infra": {
           "provider": "digitalocean",
           "type": "container"
         }
       }
     ]
   }
   MANIFEST
   ```

3. **Run your first update**:
   ```bash
   webkit update
   ```

   You should see output like:
   ```
   Updating project dependencies...

   🏃 Manifest: Scaffold manifest files
   🏃 Definition: Update webkit_version in app.json
   🏃 Env: Scaffold .env files
   🏃 Secrets: Scaffold secret files
   🏃 Files: Create code style files
   🏃 Files: Create git settings
   🏃 Files: Create package.json
   🏃 CICD: Create app PR workflows
   🏃 CICD: Creates drift detection workflow
   🏃 Env: Sync .env files
   🏃 Secrets: Sync secret files

   ✓ Successfully updated project dependencies!
   ```

4. **Check what was generated**:
   ```bash
   ls -la
   ```

   You'll see WebKit has created:
   - `.github/workflows/` - GitHub Actions CI/CD pipelines
   - `.github/settings.yml` - Repository settings
   - `apps/web/.env` - Development environment file
   - `secrets/` - SOPS configuration files
   - `package.json` - Monorepo configuration
   - `.webkit-manifest.json` - File tracking manifest (don't edit this)

## Understanding what happened

WebKit read your `app.json` and generated a complete project structure:

### Project files
- **Code style configs**: EditorConfig, Prettier, Biome configurations
- **Git settings**: `.gitignore`, GitHub repository settings
- **Monorepo tooling**: `package.json`, `pnpm-workspace.yaml` (if applicable)

### CI/CD workflows
- **PR workflows**: Build and test pipelines for each app
- **Drift detection**: Monitors for infrastructure changes
- **Backup workflows**: Automated database backups (if resources defined)

### Environment management
- **.env files**: Local development environment variables
- **Secrets scaffolding**: SOPS-encrypted secret files for each environment

### Infrastructure preparation
When you're ready to deploy, WebKit can generate:
- **Terraform modules**: Infrastructure as code for your resources and apps
- **tfvars files**: Environment-specific Terraform variables
- **State management**: Remote state configuration for Backblaze B2

## Making changes

The beauty of WebKit is that your `app.json` is the source of truth. To add a new app or resource:

1. **Edit your `app.json`**:
   ```json
   {
     "apps": [
       {
         "name": "web",
         "type": "sveltekit",
         "path": "apps/web"
       },
       {
         "name": "api",
         "type": "go",
         "path": "services/api"
       }
     ]
   }
   ```

2. **Run `webkit update` again**:
   ```bash
   webkit update
   ```

WebKit will:
- Generate files for the new `api` app
- Create new workflows for building and deploying it
- Preserve all your existing customisations
- Clean up any orphaned files from removed apps

## Validating your manifest

Before committing changes, validate your manifest:

```bash
webkit drift
```

This command checks for:
- Invalid JSON syntax
- Missing required fields
- Incorrect resource types or providers
- Dependency issues between apps and resources

## Next steps

Now that you've created your first WebKit project:

- **[Build your first project](/getting-started/your-first-project)** - Create a complete full-stack application
- **[Learn core concepts](/core-concepts/overview)** - Understand the philosophy behind WebKit
- **[Explore the manifest](/manifest/overview)** - Deep dive into `app.json` structure
- **[Deploy infrastructure](/infrastructure/overview)** - Push your project to production

## Common patterns

### Adding environment variables

Edit your `app.json`:

```json
{
  "apps": [
    {
      "name": "web",
      "env": {
        "dev": [
          {
            "key": "API_URL",
            "type": "value",
            "value": "http://localhost:8080"
          }
        ],
        "production": [
          {
            "key": "API_URL",
            "type": "value",
            "value": "https://api.my-website.com"
          }
        ]
      }
    }
  ]
}
```

Run `webkit update` and your `.env` files will be synchronised.

### Adding secrets

For sensitive values, use SOPS:

1. **Add a secret reference**:
   ```json
   {
     "env": {
       "production": [
         {
           "key": "DATABASE_PASSWORD",
           "type": "secret",
           "from": "sops:secrets/production.yaml:/DATABASE_PASSWORD"
         }
       ]
     }
   }
   ```

2. **Generate secret files**:
   ```bash
   webkit scaffold secrets
   ```

3. **Edit the encrypted file**:
   ```bash
   webkit secrets encrypt secrets/production.yaml
   ```

### Referencing resource outputs

Apps can reference infrastructure outputs:

```json
{
  "apps": [
    {
      "name": "web",
      "env": {
        "production": [
          {
            "key": "DATABASE_URL",
            "type": "from_resource",
            "from": "resource:db:connection_url"
          }
        ]
      }
    }
  ]
}
```

This creates a dependency between the `web` app and the `db` resource, which WebKit uses for:
- Docker Compose dependency ordering
- Terraform resource dependencies
- Environment variable injection

## Troubleshooting

### webkit update fails

If `webkit update` fails, check:
- Your `app.json` is valid JSON (use `jq . app.json`)
- All required fields are present
- Resource and app names are unique
- Paths reference existing directories

### Generated files look wrong

Remember:
- WebKit preserves files you've customised outside of tracked sections
- Delete `.webkit-manifest.json` to force a full regeneration (⚠️ this will overwrite customisations)
- Check that your `webkit_version` matches your installed CLI version

### Environment variables aren't syncing

Make sure you:
- Run `webkit update` after editing `app.json`
- Check the `env` block is under the correct environment (`dev`, `staging`, `production`)
- Verify secret files exist if using `type: "secret"`
