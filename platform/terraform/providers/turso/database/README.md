# Turso Database Provider

This module provisions a Turso SQLite database with automatic token generation.

## Provider

This module uses the [jpedroh/turso](https://registry.terraform.io/providers/jpedroh/turso/latest) Terraform provider.

## Usage in app.json

```json
{
  "resources": [
    {
      "name": "db",
      "type": "sqlite",
      "provider": "turso",
      "config": {
        "organisation": "my-org",
        "group": "default"
      }
    }
  ]
}
```

## Configuration Options

| Option | Type | Required | Default | Description |
|--------|------|----------|---------|-------------|
| `organisation` | string | Yes | - | Your Turso organisation name |
| `group` | string | No | `"default"` | The Turso group to create the database in |
| `size_limit` | string | No | `null` | Optional size limit for the database |

## Authentication

The Turso provider requires a `TURSO_API_TOKEN` environment variable:

### Local Development
Add to your `tf_env` file:
```bash
export TURSO_API_TOKEN="your-turso-api-token"
```

### GitHub Actions
Add `TURSO_API_TOKEN` as an organization or repository secret in GitHub.

To get your API token:
```bash
turso auth token
```

## Outputs

The following outputs are automatically exported as environment variables:

| Output | Environment Variable | Description |
|--------|---------------------|-------------|
| `connection_url` | `TF_{ENV}_{NAME}_CONNECTION_URL` | Full libsql connection URL with auth token |
| `hostname` | `TF_{ENV}_{NAME}_HOST` | Database hostname |
| `database` | `TF_{ENV}_{NAME}_DATABASE` | Database name |
| `auth_token` | `TF_{ENV}_{NAME}_AUTH_TOKEN` | Authentication token |
| `id` | `TF_{ENV}_{NAME}_ID` | Unique database ID |

Example: For a resource named `db` in production, the connection URL will be available as `TF_PROD_DB_CONNECTION_URL`.

## Example with Application

```json
{
  "resources": [
    {
      "name": "db",
      "type": "sqlite",
      "provider": "turso",
      "config": {
        "organisation": "my-org",
        "group": "default"
      }
    }
  ],
  "apps": [
    {
      "name": "api",
      "type": "golang",
      "infra": {
        "provider": "digitalocean",
        "type": "container"
      },
      "env": {
        "production": {
          "DATABASE_URL": {
            "source": "resource",
            "value": "db.connection_url"
          }
        }
      }
    }
  ]
}
```

## Resources Created

- `turso_database` - The SQLite database instance
- `turso_database_token` - Authentication token with full access and no expiration

## Notes

- Tokens are created with `expiration: "never"` and `authorization: "full-access"`
- The connection URL includes the authentication token for convenience
- Turso databases are globally distributed edge databases using libSQL
