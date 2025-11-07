#
# Default B2 Bucket
#
output "default_b2_bucket" {
  description = "Default B2 bucket details"
  value = {
    id   = module.default_b2_bucket.id
    name = module.default_b2_bucket.name
    info = module.default_b2_bucket.info
  }
  sensitive = true
}

#
# Resources (databases, storage, etc.)
#
output "resources" {
  description = "All resource outputs by resource name"
  value       = {
    for name, resource in module.resources : name => merge(
      # Common outputs (always present)
      {
        platform_type     = resource.platform_provider
        platform_provider = resource.platform_provider
      },

      # Common outputs (if they exist)
        resource.id != null ? { id = resource.id } : {},
        resource.urn != null ? { urn = resource.urn } : {},

      # Postgres-specific outputs
        resource.connection_url != null ? { connection_url = resource.connection_url } : {},
        resource.host != null ? { host = resource.host } : {},
        resource.port != null ? { port = resource.port } : {},
        resource.database != null ? { database = resource.database } : {},

      # S3-specific outputs
        resource.bucket_name != null ? { bucket_name = resource.bucket_name } : {},
        resource.bucket_url != null ? { bucket_url = resource.bucket_url } : {},
        resource.endpoint != null ? { endpoint = resource.endpoint } : {},
        resource.region != null ? { region = resource.region } : {}
    )
  }
  sensitive = true
}

output "resource_names" {
  description = "List of all provisioned resource names"
  value       = [for r in var.resources : r.name]
}

#
# Apps (services, applications)
#
output "apps" {
  description = "All app outputs by app name"
  value = {
    for name, app in module.apps : name => merge(
      # Common outputs (always present)
      {
        platform_type     = app.platform_type
        platform_provider = app.platform_provider
      },

      # VM-specific outputs
        app.ip_address != null ? { ip_address = app.ip_address } : {},
        app.droplet_id != null ? { droplet_id = app.droplet_id } : {},
        app.ssh_private_key != null ? { ssh_private_key = app.ssh_private_key } : {},
        app.server_user != null ? { server_user = app.server_user } : {},

      # Container (App Platform) outputs
        app.app_id != null ? { app_id = app.app_id } : {},
        app.app_url != null ? { app_url = app.app_url } : {},
        app.app_domain != null ? { app_domain = app.app_domain } : {}
    )
  }
  sensitive = true
}

output "app_names" {
  description = "List of all provisioned app names"
  value       = [for a in var.apps : a.name]
}

#
# Meta
#
output "project_name" {
  description = "Project name used for tagging"
  value       = var.project_name
}

output "github_secrets_created" {
  description = "List of GitHub secrets that were created"
  value = keys(local.github_secrets)
}

output "github_secrets_count" {
  description = "Number of GitHub secrets created"
  value = length(local.github_secrets)
}

#
# Slack
#
output "slack_channel_id" {
  description = "Slack channel ID for CI/CD notifications"
  value       = slack_conversation.project_channel.id
  sensitive   = false
}

output "slack_channel_name" {
  description = "Slack channel name"
  value       = slack_conversation.project_channel.name
  sensitive   = false
}
