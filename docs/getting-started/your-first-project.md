# Your first project

This guide walks you through building a complete full-stack web application using WebKit. We'll create a SvelteKit frontend with a Payload CMS backend, backed by PostgreSQL and S3 storage.

By the end, you'll have:
- A working local development environment
- Automated CI/CD pipelines
- Infrastructure ready to deploy to DigitalOcean
- Secrets management with SOPS

## Project overview

We're building **ContentHub**, a content management platform with:
- **Web app**: SvelteKit frontend for displaying content
- **CMS**: Payload CMS for content management
- **Database**: PostgreSQL for data storage
- **Storage**: S3-compatible object storage for media

## Step 1: Initialise the project

Create your project directory and initialise Git:

```bash
mkdir contenthub
cd contenthub
git init
```

Create the basic directory structure for your apps:

```bash
mkdir -p apps/web services/cms
```

## Step 2: Create the manifest

Create `app.json` in your project root:

```json
{
  "project": {
    "name": "contenthub",
    "title": "ContentHub",
    "description": "A modern content management platform",
    "repo": "https://github.com/yourusername/contenthub"
  },
  "shared": {
    "env": {
      "production": [
        {
          "key": "NODE_ENV",
          "type": "value",
          "value": "production"
        },
        {
          "key": "SENTRY_DSN",
          "type": "secret",
          "from": "sops:secrets/shared.yaml:/SENTRY_DSN"
        }
      ]
    }
  },
  "resources": [
    {
      "name": "db",
      "type": "postgres",
      "provider": "digitalocean",
      "config": {
        "size": "db-s-2vcpu-4gb",
        "engine_version": "17",
        "region": "lon1"
      },
      "outputs": [
        "connection_url",
        "host",
        "port",
        "database",
        "user"
      ]
    },
    {
      "name": "storage",
      "type": "s3",
      "provider": "digitalocean",
      "config": {
        "region": "ams3",
        "acl": "public-read"
      },
      "outputs": [
        "bucket_name",
        "endpoint",
        "region",
        "access_key_id",
        "secret_access_key"
      ]
    }
  ],
  "apps": [
    {
      "name": "cms",
      "type": "payload",
      "description": "Payload CMS for content management",
      "path": "services/cms",
      "build": {
        "dockerfile": "Dockerfile"
      },
      "infra": {
        "provider": "digitalocean",
        "type": "vm",
        "config": {
          "size": "s-2vcpu-4gb",
          "region": "lon1",
          "domain": "cms.contenthub.com"
        }
      },
      "env": {
        "dev": [
          {
            "key": "DATABASE_URL",
            "type": "from_resource",
            "from": "resource:db:connection_url"
          },
          {
            "key": "PAYLOAD_SECRET",
            "type": "value",
            "value": "dev-secret-change-in-production"
          },
          {
            "key": "S3_BUCKET",
            "type": "from_resource",
            "from": "resource:storage:bucket_name"
          }
        ],
        "production": [
          {
            "key": "DATABASE_URL",
            "type": "from_resource",
            "from": "resource:db:connection_url"
          },
          {
            "key": "PAYLOAD_SECRET",
            "type": "secret",
            "from": "sops:secrets/cms.yaml:/PAYLOAD_SECRET"
          },
          {
            "key": "S3_BUCKET",
            "type": "from_resource",
            "from": "resource:storage:bucket_name"
          },
          {
            "key": "S3_ENDPOINT",
            "type": "from_resource",
            "from": "resource:storage:endpoint"
          }
        ]
      },
      "depends_on": ["db", "storage"]
    },
    {
      "name": "web",
      "type": "sveltekit",
      "description": "Frontend web application",
      "path": "apps/web",
      "build": {
        "dockerfile": "Dockerfile"
      },
      "infra": {
        "provider": "digitalocean",
        "type": "container",
        "config": {
          "region": "lon1",
          "domain": "contenthub.com",
          "instance_count": 2
        }
      },
      "env": {
        "dev": [
          {
            "key": "PUBLIC_API_URL",
            "type": "value",
            "value": "http://localhost:3000"
          },
          {
            "key": "PRIVATE_CMS_URL",
            "type": "value",
            "value": "http://localhost:3001"
          }
        ],
        "production": [
          {
            "key": "PUBLIC_API_URL",
            "type": "value",
            "value": "https://api.contenthub.com"
          },
          {
            "key": "PRIVATE_CMS_URL",
            "type": "value",
            "value": "https://cms.contenthub.com"
          },
          {
            "key": "CDN_URL",
            "type": "from_resource",
            "from": "resource:storage:endpoint"
          }
        ]
      }
    }
  ]
}
```

## Step 3: Generate project files

Run WebKit to generate all project scaffolding:

```bash
webkit update
```

WebKit creates:
- `.github/workflows/` with CI/CD pipelines for each app
- `.env` files in `apps/web/` and `services/cms/`
- `secrets/` directory with SOPS configuration
- `package.json` for monorepo tooling
- Code style configurations (Prettier, EditorConfig, etc.)
- Git settings and `.gitignore`

## Step 4: Set up secrets

WebKit has scaffolded secret files, but they're not encrypted yet. Let's set them up:

1. **Generate an Age key** (if you haven't already):
   ```bash
   age-keygen -o age-key.txt
   cat age-key.txt
   ```
   
   Save the public key (`age1...`) somewhere safe—you'll need it.

2. **Configure SOPS**:
   
   WebKit created `.sops.yaml`. Verify it contains your Age public key.

3. **Edit and encrypt secrets**:
   ```bash
   # Edit the shared secrets file
   webkit secrets decrypt secrets/shared.yaml
   # Add your SENTRY_DSN, then encrypt:
   webkit secrets encrypt secrets/shared.yaml
   
   # Edit the CMS secrets file
   webkit secrets decrypt secrets/cms.yaml
   # Add PAYLOAD_SECRET, then encrypt:
   webkit secrets encrypt secrets/cms.yaml
   ```

::: tip
Store your Age private key securely and add it to GitHub Secrets as `AGE_SECRET_KEY` for CI/CD workflows to decrypt secrets during deployment.
:::

## Step 5: Create application code

Now we'll create minimal application code to demonstrate the full workflow.

### Create the CMS Dockerfile

`services/cms/Dockerfile`:
```dockerfile
FROM node:20-alpine

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

RUN npm run build

EXPOSE 3000

CMD ["npm", "run", "serve"]
```

### Create the web app Dockerfile

`apps/web/Dockerfile`:
```dockerfile
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

FROM node:20-alpine
WORKDIR /app
COPY --from=builder /app/build ./build
COPY --from=builder /app/package*.json ./
RUN npm install --production

EXPOSE 3000
CMD ["node", "build"]
```

::: info
In a real project, you'd scaffold your Payload CMS and SvelteKit apps using their respective CLIs (`npx create-payload-app` and `npm create svelte@latest`). We're keeping this minimal for demonstration.
:::

## Step 6: Review generated files

Let's examine what WebKit created:

### Environment files

**`services/cms/.env`** (for local development):
```env
DATABASE_URL=postgresql://user:pass@localhost:5432/contenthub
PAYLOAD_SECRET=dev-secret-change-in-production
S3_BUCKET=contenthub-storage-dev
```

**`apps/web/.env`**:
```env
PUBLIC_API_URL=http://localhost:3000
PRIVATE_CMS_URL=http://localhost:3001
```

These files are automatically synced when you run `webkit update`.

### CI/CD workflows

WebKit generated `.github/workflows/pr-cms.yml` and `.github/workflows/pr-web.yml`:

```yaml
name: PR - CMS
on:
  pull_request:
    paths:
      - 'services/cms/**'
      - 'app.json'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build Docker image
        run: |
          docker build -t cms:${{ github.sha }} ./services/cms
```

These workflows automatically build and test your apps on every pull request.

### Drift detection

`.github/workflows/drift-detection.yml` monitors for infrastructure changes:

```yaml
name: Drift Detection
on:
  schedule:
    - cron: '0 9 * * *'  # Daily at 9 AM

jobs:
  detect:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run webkit drift
        run: webkit drift
```

## Step 7: Local development

To run your project locally with Docker Compose:

1. **Generate Docker Compose file** (coming in a future WebKit version):
   ```bash
   webkit generate docker-compose
   ```

2. **Start services**:
   ```bash
   docker-compose up
   ```

Your apps will be available at:
- CMS: `http://localhost:3001`
- Web: `http://localhost:3000`

## Step 8: Deploy infrastructure

When you're ready to deploy to production:

1. **Export infrastructure environment variables**:
   ```bash
   export TF_VAR_digitalocean_token="your-do-token"
   export TF_VAR_b2_application_key_id="your-b2-key-id"
   export TF_VAR_b2_application_key="your-b2-key"
   ```

2. **Plan infrastructure changes**:
   ```bash
   webkit infra plan
   ```

   This generates Terraform configurations in a temporary directory and shows you what will be created.

3. **Apply infrastructure**:
   ```bash
   webkit infra apply
   ```

   WebKit creates:
   - PostgreSQL database cluster
   - S3 storage bucket
   - DigitalOcean droplet for the CMS
   - DigitalOcean App Platform instances for the web app

4. **View infrastructure outputs**:
   ```bash
   webkit infra output
   ```

   This displays resource outputs like database connection strings and storage endpoints.

## Step 9: Commit and push

Commit your work:

```bash
git add .
git commit -m "Initial ContentHub setup with WebKit"
git push origin main
```

Your CI/CD workflows will automatically run on pull requests, building and testing your apps.

## Making changes

### Adding a new app

Edit `app.json` to add a new API service:

```json
{
  "apps": [
    {
      "name": "api",
      "type": "go",
      "path": "services/api",
      "infra": {
        "provider": "digitalocean",
        "type": "vm"
      }
    }
  ]
}
```

Run `webkit update` and WebKit will:
- Generate `.env` files for the API
- Create CI/CD workflows for building and testing the Go app
- Update infrastructure configurations
- Preserve all existing files and customisations

### Modifying resources

Change database size:

```json
{
  "resources": [
    {
      "name": "db",
      "config": {
        "size": "db-s-4vcpu-8gb"
      }
    }
  ]
}
```

Run `webkit infra plan` to preview the change, then `webkit infra apply` to update the database.

### Updating environment variables

Add a new environment variable:

```json
{
  "apps": [
    {
      "name": "web",
      "env": {
        "production": [
          {
            "key": "ANALYTICS_ID",
            "type": "value",
            "value": "UA-12345"
          }
        ]
      }
    }
  ]
}
```

Run `webkit update` and the variable is automatically added to your environment files and deployment configurations.

## Next steps

You've built a complete full-stack application with WebKit! Here's what to explore next:

- **[Core concepts](/core-concepts/overview)** - Understand WebKit's philosophy and architecture
- **[Manifest reference](/manifest/overview)** - Deep dive into all `app.json` options
- **[CLI commands](/cli/overview)** - Explore all WebKit commands
- **[Infrastructure](/infrastructure/overview)** - Learn about Terraform integration and deployment strategies

## Troubleshooting

### Secrets won't decrypt

Ensure:
- Your Age private key is in `~/.config/sops/age/keys.txt`
- The public key in `.sops.yaml` matches your private key
- Secret files were encrypted with `webkit secrets encrypt`

### Docker build fails

Check that:
- Dockerfiles exist in app paths (`{app.path}/Dockerfile`)
- App code is present in the specified paths
- Dependencies are correctly defined in `package.json` or `go.mod`

### Infrastructure apply fails

Verify:
- Cloud provider credentials are exported as environment variables
- Provider regions are valid (check cloud provider documentation)
- Resource names are unique within the project
