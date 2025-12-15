# Apps

Apps are the core building blocks of your WebKit project. Each app represents a deployable service with its own build configuration, infrastructure settings, and environment variables.

## Attributes

| Key | Description | Required |
|-----|-------------|----------|
| `name` | Machine-readable name (kebab-case) | Yes |
| `type` | Application type | Yes |
| `path` | Relative path to the app directory | Yes |
| `description` | Human-readable description | No |
| `build` | Build configuration | No |
| `infrastructure` | Deployment settings | No |
| `domains` | Domain configuration | No |
| `environment` | Per-environment variables | No |
| `commands` | Custom build/test/lint commands | No |
| `monitoring` | Uptime monitoring settings | No |

## App types

WebKit supports these application types:

| Type | Description | Default commands |
|------|-------------|------------------|
| `svelte-kit` | SvelteKit applications | `pnpm build`, `pnpm lint`, `pnpm test` |
| `payload` | Payload CMS applications | `pnpm build`, `pnpm lint` |
| `golang` | Go applications | `go build`, `golangci-lint run`, `go test` |

Each type comes with sensible defaults for commands and configuration.

## Basic example

A minimal app definition:

```json
{
  "apps": [
    {
      "name": "web",
      "type": "svelte-kit",
      "path": "./apps/web"
    }
  ]
}
```

## Build configuration

Configure how your app is built:

```json
{
  "apps": [
    {
      "name": "api",
      "type": "golang",
      "path": "./apps/api",
      "build": {
        "dockerfile": true,
        "port": 8080,
        "health_check_path": "/api/health"
      }
    }
  ]
}
```

| Field | Description | Default |
|-------|-------------|---------|
| `dockerfile` | Generate Dockerfile | `true` |
| `port` | Exposed port | `3000` |
| `health_check_path` | Path for health check endpoint during deployment | `/` |

## Infrastructure

Define where and how your app is deployed:

```json
{
  "apps": [
    {
      "name": "web",
      "type": "svelte-kit",
      "path": "./apps/web",
      "infrastructure": {
        "provider": "digital_ocean",
        "type": "app",
        "config": {
          "instance_size": "basic-xs",
          "instance_count": 2,
          "region": "lon"
        }
      }
    }
  ]
}
```

### Provider options

| Provider | Description |
|----------|-------------|
| `digital_ocean` | DigitalOcean (App Platform or Droplets) |
| `hetzner` | Hetzner Cloud VMs |

### Infrastructure types

| Type | Description |
|------|-------------|
| `app` | Managed container platform (DigitalOcean App Platform) |
| `vm` | Virtual machine (Droplet or Hetzner server) |

### Config options

Config varies by provider and type. See the [infrastructure providers](/infrastructure/overview) documentation for detailed options.

## Domains

Configure domains for your app:

```json
{
  "apps": [
    {
      "name": "web",
      "domains": {
        "primary": "example.com",
        "aliases": ["www.example.com", "app.example.com"],
        "managed": true
      }
    }
  ]
}
```

| Field | Description |
|-------|-------------|
| `primary` | Main domain for the app |
| `aliases` | Additional domains pointing to the app |
| `managed` | Whether WebKit manages DNS records |
| `unmanaged` | Domains managed externally |

### Unmanaged domains

For domains not managed by your infrastructure provider:

```json
{
  "domains": {
    "primary": "example.com",
    "unmanaged": ["legacy.example.com"]
  }
}
```

Unmanaged domains are configured in the app but DNS must be set up manually.

## Environment variables

Define per-environment variables using the object format:

```json
{
  "apps": [
    {
      "name": "cms",
      "environment": {
        "dev": {
          "DATABASE_URL": {
            "source": "value",
            "value": "postgres://localhost:5432/cms_dev"
          }
        },
        "staging": {
          "DATABASE_URL": {
            "source": "resource",
            "value": "postgres-staging.connection_url"
          }
        },
        "production": {
          "DATABASE_URL": {
            "source": "resource",
            "value": "postgres.connection_url"
          },
          "PAYLOAD_SECRET": {
            "source": "sops",
            "value": "payload_secret"
          }
        }
      }
    }
  ]
}
```

### Variable sources

| Source | Description | Example value |
|--------|-------------|---------------|
| `value` | Static string | `"https://api.example.com"` |
| `resource` | Terraform output | `"postgres.connection_url"` |
| `sops` | Encrypted secret | `"api_key"` |

See [Environment variables](/manifest/environment-variables) for detailed documentation.

## Commands

Customise commands for build, test, lint, and format. Commands can be specified in three formats:

### Boolean format

Enable or disable a command:

```json
{
  "commands": {
    "test": false,
    "build": true
  }
}
```

### String format

Override the command with a simple string:

```json
{
  "commands": {
    "build": "go build -o bin/api ./cmd/api",
    "test": "go test -race ./..."
  }
}
```

### Object format

Full control with additional options:

```json
{
  "apps": [
    {
      "name": "api",
      "type": "golang",
      "commands": {
        "build": {
          "command": "go build -o bin/api ./cmd/api",
          "working_directory": "./cmd/api",
          "timeout": "10m"
        },
        "test": {
          "command": "go test -race ./...",
          "skip_ci": false
        },
        "lint": {
          "command": "golangci-lint run"
        }
      }
    }
  ]
}
```

| Field | Description | Default |
|-------|-------------|---------|
| `command` | Shell command to execute | Required |
| `working_directory` | Directory to run the command in | App's `path` |
| `skip_ci` | Skip this command in CI/CD workflows | `false` |
| `timeout` | Maximum execution time (e.g., `5m`, `1h`) | None |

### Working directory

By default, commands run in the app's `path` directory. Use `working_directory` to run commands in a different location:

```json
{
  "name": "web",
  "path": "apps/web",
  "commands": {
    "build": "pnpm build",
    "test": {
      "command": "pnpm test",
      "working_directory": "apps/web/src"
    }
  }
}
```

Commands are used in generated CI/CD workflows.

## Monitoring

Enable uptime monitoring for your app:

```json
{
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

### Basic HTTP monitoring

Setting `http: true` creates an HTTP monitor for your primary domain.

### Custom monitors

Add additional monitors:

```json
{
  "monitoring": {
    "http": true,
    "custom": [
      {
        "type": "http_keyword",
        "name": "API Health",
        "url": "https://api.example.com/health",
        "keyword": "\"status\":\"ok\"",
        "interval": 60
      },
      {
        "type": "dns",
        "name": "DNS Check",
        "hostname": "example.com",
        "dns_server": "8.8.8.8"
      }
    ]
  }
}
```

See [Monitoring](/manifest/monitoring) for all monitor types.

## Complete example

A fully configured app:

```json
{
  "apps": [
    {
      "name": "cms",
      "type": "payload",
      "description": "Headless CMS for content management",
      "path": "./apps/cms",
      "build": {
        "dockerfile": true,
        "port": 3000
      },
      "infrastructure": {
        "provider": "digital_ocean",
        "type": "app",
        "config": {
          "instance_size": "basic-s",
          "instance_count": 1,
          "region": "lon"
        }
      },
      "domains": {
        "primary": "cms.example.com"
      },
      "environment": {
        "dev": {
          "DATABASE_URL": {
            "source": "value",
            "value": "postgres://localhost:5432/cms_dev"
          },
          "PAYLOAD_SECRET": {
            "source": "value",
            "value": "dev-secret-change-me"
          }
        },
        "production": {
          "DATABASE_URL": {
            "source": "resource",
            "value": "postgres.connection_url"
          },
          "PAYLOAD_SECRET": {
            "source": "sops",
            "value": "payload_secret"
          },
          "S3_BUCKET": {
            "source": "resource",
            "value": "storage.bucket_name"
          },
          "S3_ENDPOINT": {
            "source": "resource",
            "value": "storage.endpoint"
          }
        }
      },
      "commands": {
        "build": "pnpm build",
        "lint": "pnpm lint"
      },
      "monitoring": {
        "http": true,
        "custom": [
          {
            "type": "http_keyword",
            "name": "CMS Health Check",
            "url": "https://cms.example.com/api/health",
            "keyword": "ok",
            "interval": 60
          }
        ]
      }
    }
  ]
}
```

## Multiple apps

WebKit supports multiple apps in a monorepo:

```json
{
  "apps": [
    {
      "name": "web",
      "type": "svelte-kit",
      "path": "./apps/web",
      "infrastructure": {
        "provider": "digital_ocean",
        "type": "app"
      },
      "domains": {
        "primary": "example.com"
      }
    },
    {
      "name": "cms",
      "type": "payload",
      "path": "./apps/cms",
      "infrastructure": {
        "provider": "digital_ocean",
        "type": "app"
      },
      "domains": {
        "primary": "cms.example.com"
      }
    },
    {
      "name": "api",
      "type": "golang",
      "path": "./apps/api",
      "infrastructure": {
        "provider": "hetzner",
        "type": "vm"
      },
      "domains": {
        "primary": "api.example.com"
      }
    }
  ]
}
```

Each app:
- Has its own build and deployment pipeline
- Can use different infrastructure providers
- Shares resources defined in the manifest

## Next steps

- Configure [resources](/manifest/resources) for databases and storage
- Set up [environment variables](/manifest/environment-variables) properly
- Add [monitoring](/manifest/monitoring) for uptime checks
- See complete [examples](/manifest/examples)
