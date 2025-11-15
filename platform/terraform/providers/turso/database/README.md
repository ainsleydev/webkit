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

The Turso provider requires a `TURSO_TOKEN` environment variable:

### Local Development
Add to your `tf_env` file:
```bash
export TURSO_TOKEN="your-turso-token"
```

### GitHub Actions
Add `TURSO_TOKEN` as `ORG_TURSO_TOKEN` organisation secret in GitHub.

To get your token:
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

## Database Naming

Turso databases created by webkit follow the naming pattern `{project-name}-{resource-name}`.

For example, if your project is named "my-app" and you define a resource named "db", the actual database created in Turso will be named `my-app-db`. This matches the naming convention used for DigitalOcean resources.

**Importing Existing Databases:**
When importing an existing Turso database, use the full database name in the import ID:
```bash
webkit infra import --resource db --id my-org/my-app-db --env production
```

## Notes

- Tokens are created with `expiration: "never"` and `authorization: "full-access"`
- The connection URL includes the authentication token for convenience
- Turso databases are globally distributed edge databases using libSQL
