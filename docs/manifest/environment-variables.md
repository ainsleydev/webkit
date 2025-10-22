# Environment Variables

Variables can either be defined within the `shared` section of the app manifest, or defined per service.

Variables can be configured for these environments:

- `dev`
- `staging`
- `production`

## Canonical object form (recommended)

Each entry is a small object with a `key`, a `source` object that describes where the value comes from, and optional metadata such as `sensitive`.

Supported `source.type` values:
- `value` — literal value (local or public value)
- `resource` — an output from a defined resource in the manifest
- `sops` — a secret stored in a SOPS-encrypted file

## Example:

```json
{
    "env": {
        "production": [
            {
                "key": "DATABASE_URL",
                "type": "resource",
                "resource": "db-primary",
                "output": "connection_url",
                "sensitive": true
            },
            {
                "key": "REDIS_HOST",
                "type": "resource",
                "resource": "cache",
                "output": "host"
            },
            {
                "key": "S3_BUCKET",
                "source": { "type": "resource", "name": "object-store", "output": "bucket_name" }
            },
            {
                "key": "PAYLOAD_SECRET",
                "type": "sops",
                "file": "secrets/prod.yaml",
                "path": "PAYLOAD_SECRET",
                "sensitive": true
            },
            {
                "key": "PUBLIC_API_URL",
                "type": "value",
                "value": "https://api.my-website.com"
            }
        ]
    }
}
```

OR

```json
{
    "env": {
        "production": {
            "DATABASE_URL": { "type": "resource", "resource": "db-primary", "output": "connection_url", "sensitive": true },
            "REDIS_HOST": { "type": "resource", "resource": "cache", "output": "host" },
            "PAYLOAD_SECRET": { "type": "sops", "file": "secrets/prod.yaml", "path": "PAYLOAD_SECRET" },
            "PUBLIC_API_URL": { "type": "value", "value": "https://api.my-website.com" }
        }
    }
}
```