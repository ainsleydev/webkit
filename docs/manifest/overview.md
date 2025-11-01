# The Manifest (`app.json`)

The manifest is the single source of truth. All generated code, infrastructure and configuration derive from this
document. Every project should define an `app.json` file within its root directory, even if it doesn't have any apps.
This will in turn, help with managing the formatting of files and templates.

## Description

- **Format**: JSON, validated against a published JSON Schema.
- **Schema Reference**: `$schema` points to a remote or versioned schema definition.
- **Versioning**: A webkit_version field is automatically added to the schema during CLI generation. When running webkit
  update, the tool will update this version to match the installed CLI version, ensuring consistency between your
  manifest and the WebKit tooling. This allows the CLI to detect version drift and apply necessary migrations or warn
  about incompatible changes.

## Example

Below is an example of a fully fledged `app.json` manifest.

```json
{
    "$schema": "https://raw.githubusercontent.com/ainsley/webkit-schema/v1.0.0/schema.json",
    "webkit_version": "0.1.0",
    "project": {
        "name": "my-website",
        "title": "My Website",
        "description": "My website is a bespoke sales platform for developers and designers.",
        "repo": "git@github.com:ainsley/my-website.git"
    },
    "shared": {
        "env": {
            "dev": [
                {
                    "key": "FRONTEND_URL",
                    "type": "value",
                    "value": "http://localhost:3000"
                }
            ],
            "staging": [
                {
                    "key": "FRONTEND_URL",
                    "type": "value",
                    "value": "https://staging.my-website.com"
                }
            ],
            "production": [
                {
                    "key": "SENTRY_DSN",
                    "type": "secret",
                    "from": "sops:secrets/shared.yaml:/SENTRY_DSN"
                },
                {
                    "key": "FRONTEND_URL",
                    "type": "value",
                    "value": "https://my-website.com"
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
                "size": "db-s-1vcpu-1gb",
                "engine_version": "17",
                "region": "ams3"
            },
            "outputs": [
                "connection_url",
                "host",
                "port",
                "database"
            ]
        },
        {
            "name": "store",
            "type": "s3",
            "provider": "digitalocean",
            "config": {
                "region": "ams3",
                "acl": "public-read"
            },
            "outputs": [
                "bucket_name",
                "endpoint",
                "region"
            ]
        }
    ],
    "apps": [
        {
            "name": "cms",
            "type": "payload",
            "description": "Payload CMS for managing content.",
            "path": "services/cms",
            "build": {
                "dockerfile": "Dockerfile"
            },
            "infra": {
                "provider": "digitalocean",
                "type": "vm",
                "config": {
                    "size": "s-2vcpu-4gb",
                    "region": "ams3",
                    "domain": "cms.my-website.com",
                    "ssh_keys": [
                        "your-ssh-key-id"
                    ]
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
                        "type": "secret",
                        "from": "sops:secrets/cms.yaml:/PAYLOAD_SECRET"
                    }
                ],
                "staging": [
                    {
                        "key": "DATABASE_URL",
                        "type": "from_resource",
                        "from": "resource:db:connection_url"
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
                        "from": "github-secrets:PAYLOAD_SECRET"
                    }
                ]
            },
            "depends_on": [
                "db"
            ]
        },
        {
            "name": "web",
            "type": "sveltekit",
            "path": "apps/web",
            "build": {
                "dockerfile": "Dockerfile"
            },
            "infra": {
                "provider": "digitalocean",
                "type": "app",
                "config": {
                    "region": "fra1",
                    "domain": "www.my-website.com",
                    "instance_count": 2,
                    "env_from_shared": true
                }
            },
            "env": {
                "dev": [
                    {
                        "key": "PUBLIC_API_URL",
                        "type": "value",
                        "value": "http://localhost:3000"
                    }
                ],
                "production": [
                    {
                        "key": "PUBLIC_API_URL",
                        "type": "value",
                        "value": "https://api.my-website.com"
                    },
                    {
                        "key": "ASSETS_BUCKET",
                        "type": "from_resource",
                        "from": "resource:object-store:bucket"
                    }
                ]
            }
        }
    ]
}
```
