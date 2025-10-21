# Infrastructure overview

WebKit generates and manages Terraform configurations from your `app.json` manifest, allowing you to deploy infrastructure without writing Terraform code directly.

## How it works

WebKit acts as a Terraform wrapper that:

1. **Generates Terraform modules** from your `app.json` resources and apps
2. **Creates tfvars files** with configuration from your manifest
3. **Manages Terraform state** in Backblaze B2 for team collaboration
4. **Resolves secrets** from SOPS files and injects them as Terraform variables
5. **Executes Terraform commands** (plan, apply, destroy) on your behalf

```mermaid
graph LR
    A[app.json] --> B[webkit infra]
    B --> C[Generate Terraform]
    C --> D[Resolve secrets]
    D --> E[Execute Terraform]
    E --> F[Cloud infrastructure]
```

## Key concepts

### Resources become Terraform modules

Each resource in `app.json` maps to a Terraform module:

```json
{
  "resources": [
    {
      "name": "db",
      "type": "postgres",
      "provider": "digitalocean",
      "config": {
        "size": "db-s-2vcpu-4gb",
        "region": "lon1"
      }
    }
  ]
}
```

WebKit generates:
```hcl
module "db" {
  source = "../modules/digitalocean/postgres"
  
  name   = "my-project-db"
  size   = "db-s-2vcpu-4gb"
  region = "lon1"
}
```

### Apps become deployment resources

Apps with `infra` blocks also generate Terraform:

```json
{
  "apps": [
    {
      "name": "web",
      "infra": {
        "provider": "digitalocean",
        "type": "container",
        "config": {
          "region": "lon1",
          "domain": "example.com"
        }
      }
    }
  ]
}
```

Generates App Platform or VM configuration depending on the `type`.

### Environment variables create dependencies

When an app references a resource output:

```json
{
  "env": {
    "production": [
      {
        "key": "DATABASE_URL",
        "type": "from_resource",
        "from": "resource:db:connection_url"
      }
    ]
  }
}
```

WebKit creates a Terraform dependency:

```hcl
resource "digitalocean_app" "web" {
  env {
    key   = "DATABASE_URL"
    value = module.db.connection_url
  }
  
  depends_on = [module.db]
}
```

This ensures resources are provisioned in the correct order.

## Deployment workflow

### Prerequisites

1. **Install Terraform** (1.0 or higher)
2. **Export provider credentials**:
   ```bash
   export TF_VAR_digitalocean_token="your-token"
   export TF_VAR_b2_application_key_id="your-b2-key-id"
   export TF_VAR_b2_application_key="your-b2-key"
   ```

### 1. Plan changes

Preview infrastructure changes without applying them:

```bash
webkit infra plan
```

This shows you what Terraform will create, update, or destroy.

### 2. Apply changes

Deploy the infrastructure:

```bash
webkit infra apply
```

WebKit:
- Generates Terraform configurations in a temporary directory
- Resolves secrets from SOPS files
- Runs `terraform apply`
- Outputs resource information

### 3. View outputs

After applying, view resource outputs:

```bash
webkit infra output
```

Example output:
```
db_connection_url = "postgresql://user:pass@db-host:5432/dbname"
storage_endpoint = "https://ams3.digitaloceanspaces.com"
storage_bucket_name = "my-project-storage"
```

### 4. Destroy infrastructure

When you're done (e.g., tearing down a staging environment):

```bash
webkit infra destroy
```

::: danger
This permanently deletes all infrastructure. Use with caution.
:::

## State management

WebKit uses Backblaze B2 for remote Terraform state storage. This allows teams to collaborate without conflicts.

**State configuration** is automatically generated from your provider credentials. WebKit creates a B2 bucket named `webkit-terraform-state` and stores state files there.

**Benefits:**
- Team members share the same state
- State is versioned and backed up
- Prevents concurrent modifications with state locking

## Provider support

WebKit currently supports:

| Provider       | Resources                             | Status     |
|----------------|---------------------------------------|------------|
| DigitalOcean   | Postgres, Redis, S3, VMs, Containers  | ✅ Stable   |
| Backblaze B2   | S3-compatible storage, state backend  | ✅ Stable   |
| AWS            | Planned (Postgres, S3, VMs, ECS)      | 🚧 Planned |

## Supported resource types

### Databases

- **postgres** - Managed PostgreSQL clusters
  - Providers: DigitalOcean
  - Outputs: `connection_url`, `host`, `port`, `database`, `user`, `password`

- **redis** - Managed Redis clusters
  - Providers: DigitalOcean
  - Outputs: `connection_url`, `host`, `port`

### Storage

- **s3** - Object storage (S3-compatible)
  - Providers: DigitalOcean (Spaces), Backblaze (B2)
  - Outputs: `bucket_name`, `endpoint`, `region`, `access_key_id`, `secret_access_key`

### Compute

- **vm** - Virtual machines (droplets)
  - Providers: DigitalOcean
  - Configuration: size, region, domain, ssh_keys

- **container** - Managed container platforms
  - Providers: DigitalOcean (App Platform)
  - Configuration: region, domain, instance_count

- **serverless** - Function-as-a-service (planned)

## Environment-specific deployments

You can deploy different infrastructure for different environments by creating separate `app.json` files or using environment variables in your Terraform configuration.

**Approach 1: Separate manifests** (recommended for significant differences)

```bash
# Staging
cp app.json app.staging.json
# Edit app.staging.json with smaller instance sizes

# Production
webkit infra apply  # Uses app.json
```

**Approach 2: Parameterised config** (for minor differences)

Use Terraform variables to adjust sizes, regions, or counts based on environment.

## Customising generated Terraform

WebKit generates Terraform in a temporary directory by default. If you need to customise the generated Terraform:

1. Generate it to a permanent location:
   ```bash
   webkit infra plan --output ./terraform
   ```

2. Customise the Terraform files in `./terraform`

3. Apply manually:
   ```bash
   cd terraform
   terraform apply
   ```

::: warning
Custom Terraform changes won't persist through `webkit infra` commands. For extensive customisation, consider managing Terraform separately.
:::

## Troubleshooting

### "Provider credentials not found"

Ensure you've exported the required environment variables:

```bash
export TF_VAR_digitalocean_token="your-token"
export TF_VAR_b2_application_key_id="your-key-id"
export TF_VAR_b2_application_key="your-key"
```

### "State lock" errors

If Terraform state is locked (usually from a previous failed run):

```bash
terraform force-unlock <lock-id>
```

### Resource already exists

If Terraform reports a resource already exists, either:
- Import it: `terraform import <resource> <id>`
- Remove it manually from your cloud provider
- Adjust the resource name in `app.json`

## Best practices

### Use remote state

Always use remote state (B2) for team projects. This prevents conflicts and ensures everyone has the latest infrastructure state.

### Plan before applying

Always run `webkit infra plan` before `webkit infra apply` to preview changes, especially in production.

### Version control your manifest

Commit `app.json` to git. Infrastructure changes should go through pull requests just like code changes.

### Test in staging first

Deploy infrastructure changes to a staging environment before production to catch issues early.

### Monitor costs

Cloud resources cost money. Regularly review your infrastructure and destroy unused resources.

## Next steps

- **[Terraform integration](/infrastructure/terraform-integration)** - Deep dive into how WebKit uses Terraform
- **[State management](/infrastructure/state-management)** - Learn about remote state
- **[Deployment workflows](/infrastructure/deployment-workflows)** - CI/CD for infrastructure
- **[CLI reference - webkit infra](/cli/webkit-infra)** - Command documentation
