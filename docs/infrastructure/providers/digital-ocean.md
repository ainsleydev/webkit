# DigitalOcean

DigitalOcean is WebKit's primary infrastructure provider, offering App Platform for containerised deployments, Droplets for VMs, managed Postgres databases, and Spaces for object storage.

## Authentication

Set your DigitalOcean API token as an environment variable:

```bash
export DIGITALOCEAN_ACCESS_TOKEN="your-token"
```

Generate a token at [cloud.digitalocean.com/account/api/tokens](https://cloud.digitalocean.com/account/api/tokens) with read/write permissions.

## App Platform

DigitalOcean App Platform provides managed container hosting with automatic scaling, SSL, and continuous deployment.

### Configuration

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
          "instance_size": "basic-xxs",
          "instance_count": 1,
          "region": "lon"
        }
      }
    }
  ]
}
```

### Instance sizes

| Size | vCPUs | Memory | Monthly cost |
|------|-------|--------|--------------|
| `basic-xxs` | Shared | 512MB | ~$5 |
| `basic-xs` | Shared | 1GB | ~$10 |
| `basic-s` | 1 | 2GB | ~$20 |
| `basic-m` | 2 | 4GB | ~$40 |
| `professional-xs` | 1 | 1GB | ~$12 |
| `professional-s` | 1 | 2GB | ~$25 |
| `professional-m` | 2 | 4GB | ~$50 |

Professional instances include dedicated CPU resources.

### Regions

| Code | Location |
|------|----------|
| `lon` | London, UK |
| `ams` | Amsterdam, Netherlands |
| `fra` | Frankfurt, Germany |
| `nyc` | New York, USA |
| `sfo` | San Francisco, USA |
| `sgp` | Singapore |
| `blr` | Bangalore, India |
| `syd` | Sydney, Australia |

### Domain configuration

Configure domains for your app:

```json
{
  "apps": [
    {
      "name": "web",
      "domains": {
        "primary": "example.com",
        "aliases": ["www.example.com"]
      }
    }
  ]
}
```

WebKit automatically:
- Configures the domain in App Platform
- Creates DNS records if the domain is managed by DigitalOcean
- Provisions SSL certificates via Let's Encrypt

### Alerts

App Platform alerts are automatically configured when you enable Slack notifications. Alerts trigger for:
- CPU utilisation above 80%
- Memory utilisation above 80%
- Application restarts (more than 3 in 5 minutes)

## Droplets

Droplets are virtual machines for workloads requiring more control than App Platform provides.

### Configuration

```json
{
  "apps": [
    {
      "name": "api",
      "type": "golang",
      "path": "./apps/api",
      "infrastructure": {
        "provider": "digital_ocean",
        "type": "vm",
        "config": {
          "size": "s-1vcpu-1gb",
          "region": "lon1",
          "image": "ubuntu-22-04-x64"
        }
      }
    }
  ]
}
```

### Droplet sizes

| Size | vCPUs | Memory | Storage | Monthly cost |
|------|-------|--------|---------|--------------|
| `s-1vcpu-512mb-10gb` | 1 | 512MB | 10GB | ~$4 |
| `s-1vcpu-1gb` | 1 | 1GB | 25GB | ~$6 |
| `s-1vcpu-2gb` | 1 | 2GB | 50GB | ~$12 |
| `s-2vcpu-2gb` | 2 | 2GB | 60GB | ~$18 |
| `s-2vcpu-4gb` | 2 | 4GB | 80GB | ~$24 |
| `s-4vcpu-8gb` | 4 | 8GB | 160GB | ~$48 |

### Regions

Droplet regions use slightly different codes than App Platform:

| Code | Location |
|------|----------|
| `lon1` | London, UK |
| `ams3` | Amsterdam, Netherlands |
| `fra1` | Frankfurt, Germany |
| `nyc1`, `nyc3` | New York, USA |
| `sfo3` | San Francisco, USA |
| `sgp1` | Singapore |
| `blr1` | Bangalore, India |

### SSH access

WebKit configures SSH keys for Droplet access. Add your public key to the project's SSH configuration or use DigitalOcean's team SSH keys.

## Managed Postgres

DigitalOcean Managed Databases provides fully managed PostgreSQL with automatic backups, failover, and maintenance.

### Configuration

```json
{
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
    }
  ]
}
```

### Database sizes

| Size | vCPUs | Memory | Storage | Monthly cost |
|------|-------|--------|---------|--------------|
| `db-s-1vcpu-1gb` | 1 | 1GB | 10GB | ~$15 |
| `db-s-1vcpu-2gb` | 1 | 2GB | 25GB | ~$30 |
| `db-s-2vcpu-4gb` | 2 | 4GB | 38GB | ~$60 |
| `db-s-4vcpu-8gb` | 4 | 8GB | 115GB | ~$120 |

### Outputs

After provisioning, these outputs are available for environment variables:

| Output | Description |
|--------|-------------|
| `postgres.connection_url` | Full connection string |
| `postgres.host` | Database host |
| `postgres.port` | Database port |
| `postgres.database` | Database name |
| `postgres.user` | Database user |
| `postgres.password` | Database password |

Use them in your manifest:

```json
{
  "environment": {
    "production": {
      "DATABASE_URL": {
        "source": "resource",
        "value": "postgres.connection_url"
      }
    }
  }
}
```

## Spaces (Object Storage)

DigitalOcean Spaces provides S3-compatible object storage for files, assets, and backups.

### Configuration

```json
{
  "resources": [
    {
      "name": "storage",
      "type": "s3",
      "provider": "digital_ocean",
      "config": {
        "region": "ams3",
        "acl": "private"
      }
    }
  ]
}
```

### Regions

| Code | Location |
|------|----------|
| `ams3` | Amsterdam |
| `fra1` | Frankfurt |
| `nyc3` | New York |
| `sfo3` | San Francisco |
| `sgp1` | Singapore |
| `syd1` | Sydney |

### ACL options

| ACL | Description |
|-----|-------------|
| `private` | No public access (default) |
| `public-read` | Public read access |

### Outputs

| Output | Description |
|--------|-------------|
| `storage.bucket_name` | Bucket name |
| `storage.endpoint` | S3 endpoint URL |
| `storage.region` | Bucket region |

### CDN

Enable CDN for public buckets:

```json
{
  "resources": [
    {
      "name": "storage",
      "type": "s3",
      "provider": "digital_ocean",
      "config": {
        "acl": "public-read",
        "cdn": {
          "enabled": true,
          "ttl": 3600
        }
      }
    }
  ]
}
```

## DNS management

WebKit can manage DNS records for domains hosted on DigitalOcean:

```json
{
  "apps": [
    {
      "name": "web",
      "domains": {
        "primary": "example.com",
        "managed": true
      }
    }
  ]
}
```

When `managed: true`, WebKit creates:
- A records pointing to App Platform
- CNAME records for aliases
- Required verification records

## Example: Full-stack application

A complete example deploying a SvelteKit frontend with Payload CMS and Postgres:

```json
{
  "project": {
    "name": "my-app",
    "title": "My Application",
    "repo": "github.com/myorg/my-app"
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
        "primary": "myapp.com",
        "aliases": ["www.myapp.com"]
      }
    },
    {
      "name": "cms",
      "type": "payload",
      "path": "./apps/cms",
      "infrastructure": {
        "provider": "digital_ocean",
        "type": "app",
        "config": {
          "instance_size": "basic-s",
          "region": "lon"
        }
      },
      "domains": {
        "primary": "cms.myapp.com"
      },
      "environment": {
        "production": {
          "DATABASE_URL": {
            "source": "resource",
            "value": "postgres.connection_url"
          }
        }
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
      }
    },
    {
      "name": "storage",
      "type": "s3",
      "provider": "digital_ocean",
      "config": {
        "region": "ams3"
      }
    }
  ]
}
```

## Further reading

- [DigitalOcean API documentation](https://docs.digitalocean.com/reference/api/)
- [App Platform documentation](https://docs.digitalocean.com/products/app-platform/)
- [Managed Databases documentation](https://docs.digitalocean.com/products/databases/)
- [Spaces documentation](https://docs.digitalocean.com/products/spaces/)
