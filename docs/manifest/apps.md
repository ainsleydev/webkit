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

### Configuring tools

Tools are defined as objects with a `type` field that determines how they're installed. WebKit supports three tool types:

#### Go tools

Go tools are installed via `go install`:

```json
{
    "apps": [
        {
            "name": "api",
            "type": "golang",
            "path": "services/api",
            "tools": {
                "custom-tool": {
                    "type": "go",
                    "name": "github.com/custom/tool/cmd/mytool",
                    "version": "v1.0.0"
                }
            }
        }
    ]
}
```

This generates: `go install github.com/custom/tool/cmd/mytool@v1.0.0`

#### pnpm tools

Node.js tools are installed globally via pnpm:

```json
{
    "tools": {
        "eslint": {
            "type": "pnpm",
            "name": "eslint",
            "version": "8.0.0"
        }
    }
}
```

This generates: `pnpm add -g eslint@8.0.0`

#### Script tools

For custom installation methods (downloading binaries, curl scripts, etc.), use the `script` type:

```json
{
    "tools": {
        "goreleaser": {
            "type": "script",
            "install": "curl -sSL https://github.com/goreleaser/goreleaser/releases/download/v1.18.2/goreleaser_Linux_x86_64.tar.gz | tar xz"
        }
    }
}
```

The `install` command is executed exactly as written.

### Overriding default tools

Default tools are automatically populated by `applyDefaults()`. To customise a default tool (like changing the version), simply include it in your `tools` configuration:

```json
{
    "tools": {
        "templ": {
            "type": "go",
            "name": "github.com/a-h/templ/cmd/templ",
            "version": "v0.2.543"
        }
    }
}
```

### Install command override

You can override the auto-generated install command for any tool type by providing an `install` field:

```json
{
    "tools": {
        "custom": {
            "type": "go",
            "name": "github.com/foo/bar",
            "version": "v1.0.0",
            "install": "custom install command"
        }
    }
}
```

### Attributes

| Key     | Description                                     | Required | Default                    | Notes                                                     |
|---------|-------------------------------------------------|----------|----------------------------|-----------------------------------------------------------|
| tools   | Map of tool names to tool configurations        | No       | Auto-populated for Go apps | Each tool is an object with type, name, version, install |
| type    | Installation method                             | Yes      | -                          | One of: "go", "pnpm", "script"                            |
| name    | Package path (go) or package name (pnpm)        | No       | -                          | Required for "go" and "pnpm" types                        |
| version | Version to install                              | No       | -                          | Required for "go" and "pnpm" types                        |
| install | Custom installation command                     | No       | -                          | Required for "script" type, optional override for others  |

### Example

```json
{
    "apps": [
        {
            "name": "api",
            "type": "golang",
            "path": "services/api",
            "tools": {
                "golangci-lint": {
                    "type": "go",
                    "name": "github.com/golangci/golangci-lint/cmd/golangci-lint",
                    "version": "v1.55.2"
                },
                "templ": {
                    "type": "go",
                    "name": "github.com/a-h/templ/cmd/templ",
                    "version": "v0.2.543"
                },
                "buf": {
                    "type": "go",
                    "name": "github.com/bufbuild/buf/cmd/buf",
                    "version": "v1.28.1"
                },
                "custom-binary": {
                    "type": "script",
                    "install": "curl -sSL https://example.com/install.sh | sh"
                }
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
- Go tools are installed via `go install`, pnpm tools via `pnpm add -g`.
- For other package managers or custom installation methods, use the `script` type.
