# Manifest examples

This page provides complete, realistic `app.json` examples for different project types. Use these as starting points for your own projects.

## Simple static site

A basic static website with S3 storage for assets.

```json
{
  "project": {
    "name": "marketing-site",
    "title": "Marketing Site",
    "description": "Company marketing website",
    "repo": "https://github.com/company/marketing-site"
  },
  "resources": [
    {
      "name": "assets",
      "type": "s3",
      "provider": "digitalocean",
      "config": {
        "region": "ams3",
        "acl": "public-read"
      },
      "outputs": ["bucket_name", "endpoint"]
    }
  ],
  "apps": [
    {
      "name": "web",
      "type": "sveltekit",
      "path": "apps/web",
      "build": {
        "dockerfile": "Dockerfile"
      },
      "infra": {
        "provider": "digitalocean",
        "type": "container",
        "config": {
          "region": "ams3",
          "domain": "example.com",
          "instance_count": 2
        }
      },
      "env": {
        "dev": [
          {
            "key": "PUBLIC_ASSETS_URL",
            "type": "value",
            "value": "http://localhost:9000"
          }
        ],
        "production": [
          {
            "key": "PUBLIC_ASSETS_URL",
            "type": "from_resource",
            "from": "resource:assets:endpoint"
          }
        ]
      }
    }
  ]
}
```

**What this creates:**
- DigitalOcean Spaces bucket for static assets
- SvelteKit app deployed to App Platform
- Environment-specific asset URLs
- CI/CD workflows for building and deploying the app

---

## Full-stack application

A complete application with frontend, backend API, database, and storage.

```json
{
  "project": {
    "name": "task-manager",
    "title": "Task Manager",
    "description": "Collaborative task management platform",
    "repo": "https://github.com/company/task-manager"
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
          "key": "LOG_LEVEL",
          "type": "value",
          "value": "info"
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
      "outputs": ["connection_url", "host", "port", "database"]
    },
    {
      "name": "cache",
      "type": "redis",
      "provider": "digitalocean",
      "config": {
        "size": "db-s-1vcpu-1gb",
        "region": "lon1"
      },
      "outputs": ["connection_url", "host", "port"]
    },
    {
      "name": "storage",
      "type": "s3",
      "provider": "digitalocean",
      "config": {
        "region": "ams3",
        "acl": "private"
      },
      "outputs": ["bucket_name", "endpoint", "access_key_id", "secret_access_key"]
    }
  ],
  "apps": [
    {
      "name": "api",
      "type": "go",
      "description": "Backend API service",
      "path": "services/api",
      "build": {
        "dockerfile": "Dockerfile",
        "args": {
          "GO_VERSION": "1.23"
        }
      },
      "infra": {
        "provider": "digitalocean",
        "type": "vm",
        "config": {
          "size": "s-2vcpu-4gb",
          "region": "lon1",
          "domain": "api.taskmanager.com"
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
            "key": "REDIS_URL",
            "type": "from_resource",
            "from": "resource:cache:connection_url"
          },
          {
            "key": "JWT_SECRET",
            "type": "value",
            "value": "dev-jwt-secret"
          }
        ],
        "production": [
          {
            "key": "DATABASE_URL",
            "type": "from_resource",
            "from": "resource:db:connection_url"
          },
          {
            "key": "REDIS_URL",
            "type": "from_resource",
            "from": "resource:cache:connection_url"
          },
          {
            "key": "JWT_SECRET",
            "type": "secret",
            "from": "sops:secrets/api.yaml:/JWT_SECRET"
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
      "depends_on": ["db", "cache"]
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
          "domain": "taskmanager.com",
          "instance_count": 3
        }
      },
      "env": {
        "dev": [
          {
            "key": "PUBLIC_API_URL",
            "type": "value",
            "value": "http://localhost:8080"
          }
        ],
        "production": [
          {
            "key": "PUBLIC_API_URL",
            "type": "value",
            "value": "https://api.taskmanager.com"
          },
          {
            "key": "PUBLIC_CDN_URL",
            "type": "from_resource",
            "from": "resource:storage:endpoint"
          }
        ]
      }
    }
  ]
}
```

**What this creates:**
- PostgreSQL database for persistent data
- Redis cache for sessions and caching
- S3 bucket for file uploads
- Go API service on a VM
- SvelteKit frontend on App Platform
- Complete CI/CD pipelines for both apps
- Environment-specific configurations
- Secrets management setup

---

## Content management system

A CMS-based project with Payload CMS, database, and media storage.

```json
{
  "project": {
    "name": "blog-platform",
    "title": "Blog Platform",
    "description": "Multi-author blog with headless CMS",
    "repo": "https://github.com/company/blog-platform"
  },
  "shared": {
    "env": {
      "dev": [
        {
          "key": "NODE_ENV",
          "type": "value",
          "value": "development"
        }
      ],
      "production": [
        {
          "key": "NODE_ENV",
          "type": "value",
          "value": "production"
        },
        {
          "key": "SENDGRID_API_KEY",
          "type": "secret",
          "from": "sops:secrets/shared.yaml:/SENDGRID_API_KEY"
        }
      ]
    }
  },
  "resources": [
    {
      "name": "db",
      "type": "postgres",
      "provider": "digitalocean",
      "description": "Primary database for CMS and app data",
      "config": {
        "size": "db-s-4vcpu-8gb",
        "engine_version": "17",
        "region": "nyc3"
      },
      "outputs": ["connection_url", "host", "port", "database", "user"]
    },
    {
      "name": "media-storage",
      "type": "s3",
      "provider": "digitalocean",
      "description": "Media file storage",
      "config": {
        "region": "nyc3",
        "acl": "public-read"
      },
      "outputs": ["bucket_name", "endpoint", "region"]
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
          "size": "s-4vcpu-8gb",
          "region": "nyc3",
          "domain": "cms.blogplatform.com"
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
            "value": "dev-secret-unsafe"
          },
          {
            "key": "PAYLOAD_PUBLIC_SERVER_URL",
            "type": "value",
            "value": "http://localhost:3000"
          },
          {
            "key": "S3_BUCKET",
            "type": "from_resource",
            "from": "resource:media-storage:bucket_name"
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
            "key": "PAYLOAD_PUBLIC_SERVER_URL",
            "type": "value",
            "value": "https://cms.blogplatform.com"
          },
          {
            "key": "S3_BUCKET",
            "type": "from_resource",
            "from": "resource:media-storage:bucket_name"
          },
          {
            "key": "S3_ENDPOINT",
            "type": "from_resource",
            "from": "resource:media-storage:endpoint"
          },
          {
            "key": "S3_REGION",
            "type": "from_resource",
            "from": "resource:media-storage:region"
          }
        ]
      },
      "depends_on": ["db", "media-storage"]
    },
    {
      "name": "web",
      "type": "sveltekit",
      "description": "Public-facing blog website",
      "path": "apps/web",
      "build": {
        "dockerfile": "Dockerfile"
      },
      "infra": {
        "provider": "digitalocean",
        "type": "container",
        "config": {
          "region": "nyc3",
          "domain": "blogplatform.com",
          "instance_count": 2
        }
      },
      "env": {
        "dev": [
          {
            "key": "PUBLIC_CMS_URL",
            "type": "value",
            "value": "http://localhost:3000"
          },
          {
            "key": "PRIVATE_CMS_API_KEY",
            "type": "value",
            "value": "dev-api-key"
          }
        ],
        "production": [
          {
            "key": "PUBLIC_CMS_URL",
            "type": "value",
            "value": "https://cms.blogplatform.com"
          },
          {
            "key": "PRIVATE_CMS_API_KEY",
            "type": "secret",
            "from": "sops:secrets/web.yaml:/CMS_API_KEY"
          },
          {
            "key": "PUBLIC_CDN_URL",
            "type": "from_resource",
            "from": "resource:media-storage:endpoint"
          }
        ]
      }
    }
  ]
}
```

**What this creates:**
- Payload CMS for content management
- PostgreSQL database
- S3 bucket for media files
- SvelteKit frontend for displaying content
- Automated database backups
- CI/CD workflows for both apps

---

## Microservices architecture

A multi-service application with shared resources.

```json
{
  "project": {
    "name": "ecommerce-platform",
    "title": "E-commerce Platform",
    "description": "Microservices-based e-commerce system",
    "repo": "https://github.com/company/ecommerce"
  },
  "shared": {
    "env": {
      "production": [
        {
          "key": "STRIPE_PUBLIC_KEY",
          "type": "value",
          "value": "pk_live_..."
        },
        {
          "key": "STRIPE_SECRET_KEY",
          "type": "secret",
          "from": "sops:secrets/shared.yaml:/STRIPE_SECRET_KEY"
        }
      ]
    }
  },
  "resources": [
    {
      "name": "main-db",
      "type": "postgres",
      "provider": "digitalocean",
      "config": {
        "size": "db-s-8vcpu-16gb",
        "engine_version": "17",
        "region": "fra1"
      },
      "outputs": ["connection_url"]
    },
    {
      "name": "cache",
      "type": "redis",
      "provider": "digitalocean",
      "config": {
        "size": "db-s-2vcpu-4gb",
        "region": "fra1"
      },
      "outputs": ["connection_url"]
    },
    {
      "name": "product-images",
      "type": "s3",
      "provider": "digitalocean",
      "config": {
        "region": "ams3",
        "acl": "public-read"
      },
      "outputs": ["bucket_name", "endpoint"]
    }
  ],
  "apps": [
    {
      "name": "auth-service",
      "type": "go",
      "path": "services/auth",
      "infra": {
        "provider": "digitalocean",
        "type": "vm",
        "config": {
          "size": "s-2vcpu-4gb",
          "region": "fra1",
          "domain": "auth.ecommerce.com"
        }
      },
      "env": {
        "production": [
          {
            "key": "DATABASE_URL",
            "type": "from_resource",
            "from": "resource:main-db:connection_url"
          },
          {
            "key": "REDIS_URL",
            "type": "from_resource",
            "from": "resource:cache:connection_url"
          },
          {
            "key": "JWT_SECRET",
            "type": "secret",
            "from": "sops:secrets/auth.yaml:/JWT_SECRET"
          }
        ]
      }
    },
    {
      "name": "product-service",
      "type": "go",
      "path": "services/product",
      "infra": {
        "provider": "digitalocean",
        "type": "vm",
        "config": {
          "size": "s-2vcpu-4gb",
          "region": "fra1",
          "domain": "products.ecommerce.com"
        }
      },
      "env": {
        "production": [
          {
            "key": "DATABASE_URL",
            "type": "from_resource",
            "from": "resource:main-db:connection_url"
          },
          {
            "key": "REDIS_URL",
            "type": "from_resource",
            "from": "resource:cache:connection_url"
          },
          {
            "key": "IMAGE_BUCKET",
            "type": "from_resource",
            "from": "resource:product-images:bucket_name"
          }
        ]
      }
    },
    {
      "name": "order-service",
      "type": "go",
      "path": "services/order",
      "infra": {
        "provider": "digitalocean",
        "type": "vm",
        "config": {
          "size": "s-4vcpu-8gb",
          "region": "fra1",
          "domain": "orders.ecommerce.com"
        }
      },
      "env": {
        "production": [
          {
            "key": "DATABASE_URL",
            "type": "from_resource",
            "from": "resource:main-db:connection_url"
          },
          {
            "key": "REDIS_URL",
            "type": "from_resource",
            "from": "resource:cache:connection_url"
          }
        ]
      }
    },
    {
      "name": "storefront",
      "type": "sveltekit",
      "path": "apps/storefront",
      "infra": {
        "provider": "digitalocean",
        "type": "container",
        "config": {
          "region": "fra1",
          "domain": "shop.ecommerce.com",
          "instance_count": 4
        }
      },
      "env": {
        "production": [
          {
            "key": "PUBLIC_AUTH_URL",
            "type": "value",
            "value": "https://auth.ecommerce.com"
          },
          {
            "key": "PUBLIC_PRODUCT_URL",
            "type": "value",
            "value": "https://products.ecommerce.com"
          },
          {
            "key": "PUBLIC_ORDER_URL",
            "type": "value",
            "value": "https://orders.ecommerce.com"
          },
          {
            "key": "PUBLIC_CDN_URL",
            "type": "from_resource",
            "from": "resource:product-images:endpoint"
          }
        ]
      }
    }
  ]
}
```

**What this creates:**
- Three Go microservices (auth, product, order)
- Shared PostgreSQL database
- Redis cache for all services
- S3 bucket for product images
- SvelteKit storefront
- CI/CD pipelines for each service
- Service-to-service communication setup

---

## Key patterns

These examples demonstrate common patterns you can use in your own projects:

### Shared environment variables

```json
{
  "shared": {
    "env": {
      "production": [
        {
          "key": "LOG_LEVEL",
          "type": "value",
          "value": "info"
        }
      ]
    }
  }
}
```

All apps inherit shared environment variables. App-specific variables override shared ones.

### Resource outputs

```json
{
  "key": "DATABASE_URL",
  "type": "from_resource",
  "from": "resource:db:connection_url"
}
```

Reference infrastructure outputs in environment variables. WebKit creates Terraform dependencies automatically.

### Secrets management

```json
{
  "key": "API_KEY",
  "type": "secret",
  "from": "sops:secrets/production.yaml:/API_KEY"
}
```

Store sensitive values in SOPS-encrypted files, referenced in `app.json`.

### Build arguments

```json
{
  "build": {
    "dockerfile": "Dockerfile.prod",
    "args": {
      "NODE_VERSION": "20",
      "BUILD_ENV": "production"
    }
  }
}
```

Pass build-time arguments to Docker builds.

### Multi-region deployment

```json
{
  "resources": [
    {
      "name": "db-primary",
      "type": "postgres",
      "config": {
        "region": "nyc3"
      }
    }
  ],
  "apps": [
    {
      "name": "api-us",
      "infra": {
        "config": {
          "region": "nyc3"
        }
      }
    },
    {
      "name": "api-eu",
      "infra": {
        "config": {
          "region": "fra1"
        }
      }
    }
  ]
}
```

Deploy apps to multiple regions by defining separate app entries.

## Next steps

- **[Project configuration](/manifest/project)** - Define project metadata
- **[Apps configuration](/manifest/apps)** - Configure application deployments
- **[Resources configuration](/manifest/resources)** - Provision infrastructure
- **[Environment variables](/manifest/environment-variables)** - Manage configuration and secrets
