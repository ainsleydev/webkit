# Resources

Resources define infrastructure components that applications depend on. They act as declarative specifications of
services such as databases, caches and storage buckets that apps can use.

## Attributes

| Key           | Description                                                  | Required | Notes                           |
|---------------|--------------------------------------------------------------|----------|---------------------------------|
| `name`        | Project machine-readable name                                | Yes      | kebab-case                      |
| `type`        | Type of infrastructure or deployment unit                    | Yes      | Supported: `s3`, `postgres`     |
| `provider`    | Cloud provider where the infrastructure is provisioned       | Yes      | Supported: `digitalocean`, `b2` |
| `description` | Description of the resource                                  | No       |                                 |
| `config`      | Terraform input configuration based on the type and provider | Yes      |                                 |
| `outputs`     | Terraform outputs based on the type and provider             | No       |                                 |

::: warning
Each provider variable and output needs to be documented according to each module. To be confirmed how this should be
done.
:::

## Example

 ```json
{
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
    ]
}
```

## Config

The `config` key directly relates to the Terraform configuration for a provider. Each provider has a subset of variables
that are exposed in each module. For example, a `digitalocean`, `postgres` resource exposes `name`, `pg_version` etc.

::: info
More resources can be added at a later date such as Redis and other components.
:::

## Outputs

**How Outputs Work:**

Resources expose outputs that can be referenced by apps and other resources. Outputs must be explicitly declared in the
`outputs` block of each resource definition.

**Output Declaration:**

Each resource declares which outputs are available. The values come from Terraform outputs after provisioning.

**How This Maps to Terraform:**

WebKit generates Terraform modules for each resource. The `outputs` array in `app.json` corresponds directly to
Terraform outputs:

```terraform
# Generated: terraform/modules/db/outputs.tf
output "connection_url" {
  value     = digitalocean_database_cluster.primary.uri
  sensitive = true
}

output "host" {
  value = digitalocean_database_cluster.primary.host
}

output "port" {
  value = digitalocean_database_cluster.primary.port
}

output "database" {
  value = digitalocean_database_cluster.primary.database
}
```

