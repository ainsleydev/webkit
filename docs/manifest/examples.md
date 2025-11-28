# Examples

This page provides complete `app.json` examples for common project configurations. Use these as starting points for your own projects.

## SvelteKit static site

A simple static site deployed to DigitalOcean App Platform:

```json
{
  "$schema": "https://raw.githubusercontent.com/ainsleydev/webkit/main/schema.json",
  "project": {
    "name": "portfolio",
    "title": "My Portfolio",
    "description": "Personal portfolio website",
    "repo": "github.com/username/portfolio"
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
          "instance_size": "basic-xxs",
          "instance_count": 1,
          "region": "lon"
        }
      },
      "domains": {
        "primary": "portfolio.example.com"
      },
      "environment": {
        "production": {
          "PUBLIC_SITE_URL": {
            "source": "value",
            "value": "https://portfolio.example.com"
          }
        }
      }
    }
  ]
}
```

## Payload CMS with Postgres

A Payload CMS application with a managed PostgreSQL database:

```json
{
  "$schema": "https://raw.githubusercontent.com/ainsleydev/webkit/main/schema.json",
  "project": {
    "name": "content-cms",
    "title": "Content CMS",
    "description": "Headless CMS for content management",
    "repo": "github.com/username/content-cms"
  },
  "apps": [
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
            "value": "dev-secret-change-in-production"
          }
        },
        "production": {
          "DATABASE_URL": {
            "source": "resource",
            "value": "postgres.connection_url"
          },
          "PAYLOAD_SECRET": {
            "source": "sops",
            "path": "payload_secret"
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
  ]
}
```

## Full-stack application

A complete full-stack setup with SvelteKit frontend and Payload CMS backend:

```json
{
  "$schema": "https://raw.githubusercontent.com/ainsleydev/webkit/main/schema.json",
  "project": {
    "name": "my-saas",
    "title": "My SaaS Application",
    "description": "Full-stack SaaS application",
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
          },
          "PUBLIC_SITE_URL": {
            "source": "value",
            "value": "https://app.example.com"
          }
        },
        "staging": {
          "PUBLIC_API_URL": {
            "source": "value",
            "value": "https://staging-api.example.com"
          },
          "PUBLIC_SITE_URL": {
            "source": "value",
            "value": "https://staging.example.com"
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
          "instance_count": 1,
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
            "path": "payload_secret"
          },
          "S3_BUCKET": {
            "source": "resource",
            "value": "media.bucket_name"
          },
          "S3_ACCESS_KEY": {
            "source": "sops",
            "path": "s3_access_key"
          },
          "S3_SECRET_KEY": {
            "source": "sops",
            "path": "s3_secret_key"
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
        "size": "db-s-1vcpu-2gb",
        "region": "lon1",
        "version": "15"
      },
      "backup": {
        "enabled": true,
        "schedule": "0 2 * * *",
        "retention_days": 14
      }
    },
    {
      "name": "media",
      "type": "s3",
      "provider": "digital_ocean",
      "config": {
        "region": "ams3",
        "acl": "public-read",
        "cdn": {
          "enabled": true,
          "ttl": 86400
        }
      }
    },
    {
      "name": "backups",
      "type": "s3",
      "provider": "b2",
      "config": {
        "bucket_type": "allPrivate"
      }
    }
  ],
  "monitoring": {
    "status_page": {
      "name": "My SaaS Status",
      "slug": "my-saas-status"
    }
  }
}
```

## Go API on Hetzner

A Go API deployed to a Hetzner VM with attached storage:

```json
{
  "$schema": "https://raw.githubusercontent.com/ainsleydev/webkit/main/schema.json",
  "project": {
    "name": "api-service",
    "title": "API Service",
    "description": "Backend API service",
    "repo": "github.com/username/api-service"
  },
  "apps": [
    {
      "name": "api",
      "type": "golang",
      "path": "./apps/api",
      "build": {
        "dockerfile": true,
        "port": 8080
      },
      "infrastructure": {
        "provider": "hetzner",
        "type": "vm",
        "config": {
          "server_type": "cx22",
          "location": "nbg1",
          "image": "ubuntu-22.04"
        }
      },
      "domains": {
        "primary": "api.example.com"
      },
      "commands": {
        "build": {
          "run": "go build -o bin/api ./cmd/api"
        },
        "test": {
          "run": "go test ./..."
        },
        "lint": {
          "run": "golangci-lint run"
        }
      },
      "environment": {
        "production": {
          "DATABASE_URL": {
            "source": "sops",
            "path": "database_url"
          },
          "STORAGE_PATH": {
            "source": "value",
            "value": "/data/uploads"
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
      "name": "data",
      "type": "volume",
      "provider": "hetzner",
      "config": {
        "size": 100,
        "location": "nbg1",
        "mount_point": "/data"
      }
    }
  ]
}
```

## Edge application with Turso

A SvelteKit application using Turso for edge-distributed SQLite:

```json
{
  "$schema": "https://raw.githubusercontent.com/ainsleydev/webkit/main/schema.json",
  "project": {
    "name": "edge-app",
    "title": "Edge Application",
    "description": "Globally distributed edge application",
    "repo": "github.com/username/edge-app"
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
          "region": "lon"
        }
      },
      "domains": {
        "primary": "edge-app.example.com"
      },
      "environment": {
        "production": {
          "DATABASE_URL": {
            "source": "resource",
            "value": "db.url"
          },
          "DATABASE_AUTH_TOKEN": {
            "source": "resource",
            "value": "db.auth_token"
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
      "name": "db",
      "type": "sqlite",
      "provider": "turso",
      "config": {
        "group": "production",
        "primary_location": "lhr"
      }
    }
  ]
}
```

## Multi-environment configuration

Demonstrating environment-specific configuration:

```json
{
  "$schema": "https://raw.githubusercontent.com/ainsleydev/webkit/main/schema.json",
  "project": {
    "name": "multi-env",
    "title": "Multi-Environment App",
    "repo": "github.com/username/multi-env"
  },
  "apps": [
    {
      "name": "web",
      "type": "svelte-kit",
      "path": "./apps/web",
      "infrastructure": {
        "provider": "digital_ocean",
        "type": "app"
      },
      "environment": {
        "dev": {
          "PUBLIC_API_URL": {
            "source": "value",
            "value": "http://localhost:3001"
          },
          "DATABASE_URL": {
            "source": "value",
            "value": "postgres://localhost:5432/app_dev"
          },
          "LOG_LEVEL": {
            "source": "value",
            "value": "debug"
          }
        },
        "staging": {
          "PUBLIC_API_URL": {
            "source": "value",
            "value": "https://staging-api.example.com"
          },
          "DATABASE_URL": {
            "source": "resource",
            "value": "postgres-staging.connection_url"
          },
          "LOG_LEVEL": {
            "source": "value",
            "value": "info"
          }
        },
        "production": {
          "PUBLIC_API_URL": {
            "source": "value",
            "value": "https://api.example.com"
          },
          "DATABASE_URL": {
            "source": "resource",
            "value": "postgres-prod.connection_url"
          },
          "LOG_LEVEL": {
            "source": "value",
            "value": "warn"
          },
          "SENTRY_DSN": {
            "source": "sops",
            "path": "sentry_dsn"
          }
        }
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
      },
      "staging": {
        "NODE_ENV": {
          "source": "value",
          "value": "staging"
        }
      },
      "dev": {
        "NODE_ENV": {
          "source": "value",
          "value": "development"
        }
      }
    }
  }
}
```

## With monitoring and notifications

Complete setup with uptime monitoring and Slack notifications:

```json
{
  "$schema": "https://raw.githubusercontent.com/ainsleydev/webkit/main/schema.json",
  "project": {
    "name": "monitored-app",
    "title": "Monitored Application",
    "repo": "github.com/username/monitored-app"
  },
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
          "region": "lon"
        }
      },
      "domains": {
        "primary": "monitored-app.example.com"
      },
      "monitoring": {
        "http": true,
        "custom": [
          {
            "type": "http_keyword",
            "name": "API Health Check",
            "url": "https://monitored-app.example.com/api/health",
            "keyword": "\"status\":\"ok\"",
            "interval": 60
          }
        ]
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
        "region": "lon1"
      },
      "backup": {
        "enabled": true,
        "schedule": "0 3 * * *"
      }
    }
  ],
  "monitoring": {
    "status_page": {
      "name": "Monitored App Status",
      "slug": "monitored-app-status",
      "custom_domain": "status.example.com"
    },
    "notifications": {
      "slack": true
    },
    "custom": [
      {
        "type": "dns",
        "name": "DNS Check",
        "hostname": "monitored-app.example.com",
        "dns_server": "8.8.8.8"
      }
    ]
  }
}
```

## Next steps

These examples cover the most common configurations. For detailed documentation on each section:

- [Project configuration](/manifest/project)
- [Apps configuration](/manifest/apps)
- [Resources configuration](/manifest/resources)
- [Environment variables](/manifest/environment-variables)
- [Monitoring configuration](/manifest/monitoring)
