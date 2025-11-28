# The manifest (app.json)

The manifest is the single source of truth for your WebKit project. All generated code, infrastructure, and configuration derive from this file. Every project should define an `app.json` file in its root directory.

## Format

- **Type**: JSON with schema validation
- **Schema**: The `$schema` property enables IDE autocompletion and validation
- **Version tracking**: WebKit automatically manages the `webkit_version` field

## Schema reference

Always include the schema reference for IDE support:

```json
{
  "$schema": "https://raw.githubusercontent.com/ainsleydev/webkit/main/schema.json"
}
```

This provides:
- Autocompletion in VS Code and JetBrains IDEs
- Inline validation errors
- Documentation on hover

## Structure overview

A manifest consists of these top-level sections:

| Section | Description | Required |
|---------|-------------|----------|
| `$schema` | JSON Schema URL for validation | Recommended |
| `webkit_version` | CLI version that generated the manifest | Auto-managed |
| `project` | Project metadata (name, title, repo) | Yes |
| `apps` | Application definitions | Yes |
| `resources` | Infrastructure resources (databases, storage) | No |
| `shared` | Shared configuration across apps | No |
| `monitoring` | Uptime monitoring and status pages | No |

## Minimal example

The simplest valid manifest:

```json
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
```

## Complete example

A full manifest demonstrating all features:

```json
{
  "$schema": "https://raw.githubusercontent.com/ainsleydev/webkit/main/schema.json",
  "project": {
    "name": "my-saas",
    "title": "My SaaS Application",
    "description": "A full-stack SaaS platform",
    "repo": "github.com/username/my-saas"
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
        "type": "app",
        "config": {
          "instance_size": "basic-xs",
          "instance_count": 2,
          "region": "lon"
        }
      },
      "domains": {
        "primary": "app.example.com",
        "aliases": ["www.example.com"]
      },
      "environment": {
        "production": {
          "PUBLIC_API_URL": {
            "source": "value",
            "value": "https://api.example.com"
          }
        },
        "staging": {
          "PUBLIC_API_URL": {
            "source": "value",
            "value": "https://staging-api.example.com"
          }
        }
      },
      "monitoring": {
        "http": true
      }
    },
    {
      "name": "cms",
      "type": "payload",
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
          "region": "lon"
        }
      },
      "domains": {
        "primary": "api.example.com"
      },
      "environment": {
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
          }
        }
      },
      "monitoring": {
        "http": true
      }
    }
  ],
  "resources": [
    {
      "name": "postgres",
      "type": "postgres",
      "provider": "digital_ocean",
      "config": {
        "size": "db-s-1vcpu-1gb",
        "region": "lon1",
        "version": "15"
      },
      "backup": {
        "enabled": true,
        "schedule": "0 3 * * *"
      }
    },
    {
      "name": "storage",
      "type": "s3",
      "provider": "digital_ocean",
      "config": {
        "region": "ams3",
        "acl": "public-read"
      }
    }
  ],
  "shared": {
    "environment": {
      "production": {
        "NODE_ENV": {
          "source": "value",
          "value": "production"
        }
      }
    }
  },
  "monitoring": {
    "status_page": {
      "name": "My SaaS Status",
      "slug": "my-saas-status"
    }
  }
}
```

## Version tracking

WebKit automatically manages the `webkit_version` field. When you run `webkit update`, it:

1. Updates the version to match your installed CLI
2. Detects version drift between manifest and CLI
3. Applies necessary migrations for breaking changes

You don't need to set this field manually.

## Validation

Validate your manifest before deployment:

```bash
webkit validate
```

This checks:
- Required fields are present
- Values match expected types
- Resource references are valid
- Domain configurations are correct

## Next steps

Learn more about each section:

- [Project configuration](/manifest/project)
- [Apps configuration](/manifest/apps)
- [Resources configuration](/manifest/resources)
- [Environment variables](/manifest/environment-variables)
- [Monitoring configuration](/manifest/monitoring)
- [Complete examples](/manifest/examples)
