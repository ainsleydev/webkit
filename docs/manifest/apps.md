# Apps

An app, refers to a service that will exist both locally and in the cloud, this is the backbone of the application.
Every app defines a Dockerfile so it can be ran in many different environments.

## Attributes

| Key           | Description                                         | Required | Notes                                       |
|---------------|-----------------------------------------------------|----------|---------------------------------------------|
| `name`        | App machine-readable name                           | Yes      | kebab-case                                  |
| `type`        | The type of app                                     | Yes      | Supported: `payload`, `sveltekit`, `go`     |
| `description` | Description of the app                              | No       |                                             |
| `path`        | The relative path of where the application resides  | Yes      |                                             |
| `build`       | Instructions for compilation                        | No       |                                             |
| `infra`       | Infrastructure and provisioning details for the app | Yes      |                                             |
| `env`         | Per-environment variables                           | No       | [Ref](/manifest/environment-variables.html) |

## Example

```json
{
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

## Infrastructure

The `infrastructure` block defines how the application will be provisioned in cloud environments. It is similar to the
resources block but focuses on the runtime environment, compute instances, and cloud provider configurations required to
deploy the app.

### Attributes

| Key           | Description                                                  | Required | Notes                           |
|---------------|--------------------------------------------------------------|----------|---------------------------------|
| `type`        | Type of infrastructure or deployment unit                    | Yes      | Supported values below          |
| `provider`    | Cloud provider where the infrastructure is provisioned       | Yes      | Supported: `digitalocean`, `b2` |
| `config`      | Terraform input configuration based on the type and provider | Yes      |                                 |
| `description` | Description of the resource                                  | No       |                                 |

### Types

WebKit uses generic type names. The CLI maps these to provider-specific resources automatically.

| Type         | Description                | DigitalOcean | 
|--------------|----------------------------|--------------|
| `vm`         | Virtual machine            | Droplet      | 
| `container`  | Managed container platform | App Platform |
| `serverless` | Function-as-a-service      | Functions    |

### Depends On

Controls startup order in local development (Docker Compose). In production, Terraform automatically handles
provisioning order through environment variable references. You don't need to explicitly declare
dependencies—referencing a resource's output creates the dependency.

**When to use:**
Only needed if your app requires a dependency for local development that isn't referenced in environment variables.

### Example

```json
{
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
    }
}
```

## Build

The `build` block defines how each app app is compiled and packaged. All apps must include a `Dockerfile` in their root
path. It's assumed that every app will define its own Dockerfile so it can be executed and ran on cloud environments.
Dockerfile paths can be overridden with the `dockerfile` key.

Arguments can be passed in to each dockerfile using the `args` parameter as a key value pair.

### Example

```json
{
    "build": {
        "dockerfile": "Dockerfile.custom",
        "args": {
            "NODE_VERSION": "20",
            "BUILD_ENV": "production"
        }
    }
}
```

### Attributes

| Key          | Description                     | Required | Default      | Notes |
|--------------|---------------------------------|----------|--------------|-------|
| `dockerfile` | Custom Dockerfile name          | No       | `Dockerfile` |       |
| `args`       | Build-time arguments for Docker | No       | `{}`         |       |

### Notes

- Every app must have a `Dockerfile` at `{app.path}/Dockerfile`
- Build args are passed to Docker with `--build-arg`
- For advanced Docker features (multi-stage builds, BuildKit), modify your Dockerfile directly.

## Tools and dependencies

The `tools` field allows you to specify build tools and their versions that are required for CI/CD pipelines. WebKit automatically installs these tools in GitHub Actions workflows before running your commands.

### Default behaviour

WebKit provides sensible defaults for common tools based on your app type. For Go applications, the following tools are installed automatically:

- `golangci-lint` - For linting
- `templ` - For template generation
- `sqlc` - For SQL code generation

JavaScript applications (Payload, SvelteKit) don't have default tools, as they typically install dependencies via pnpm.

### Overriding tool versions

You can override the version of any default tool by specifying it in the `tools` field:

```json
{
    "apps": [
        {
            "name": "api",
            "type": "golang",
            "path": "services/api",
            "tools": {
                "templ": "v0.2.543",
                "golangci-lint": "v1.55.2"
            }
        }
    ]
}
```

### Adding custom tools

You can add tools that aren't in the default set by including them in the `tools` field:

```json
{
    "tools": {
        "buf": "v1.28.1",
        "wire": "latest"
    }
}
```

WebKit recognises common Go tool names and knows how to install them. Supported tool names include:

- `golangci-lint`
- `templ`
- `sqlc`
- `buf`
- `wire`
- `mockgen`

For tools not in the registry, specify the full Go module path:

```json
{
    "tools": {
        "github.com/custom/tool/cmd/mytool": "v1.0.0"
    }
}
```

### Disabling default tools

If you don't need a default tool, you can disable it by setting its version to an empty string or `"disabled"`:

```json
{
    "tools": {
        "sqlc": "disabled",
        "templ": ""
    }
}
```

### Attributes

| Key   | Description                                                             | Required | Default                           | Notes                                |
|-------|-------------------------------------------------------------------------|----------|-----------------------------------|--------------------------------------|
| tools | Map of tool names to versions (e.g., `"golangci-lint": "latest"`)      | No       | Auto-populated for Go apps        | Only applies to Go apps in CI/CD     |

### Example

```json
{
    "apps": [
        {
            "name": "api",
            "type": "golang",
            "path": "services/api",
            "tools": {
                "golangci-lint": "v1.55.2",
                "templ": "v0.2.543",
                "sqlc": "disabled",
                "buf": "v1.28.1"
            },
            "commands": {
                "lint": "golangci-lint run",
                "generate": "templ generate && buf generate"
            }
        }
    ]
}
```

### Notes

- Tools are only installed in CI/CD workflows, not in local development.
- For local development, install tools manually or use a tool like `asdf`.
- Version `"latest"` installs the most recent release of the tool.
- Tool installation happens after Go is set up but before commands are run.
