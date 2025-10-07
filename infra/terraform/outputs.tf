# Output all resources in a structured format
# This allows CI/CD and applications to reference resource outputs
# Format: resources[resource_name][output_name]
output "resources" {
  description = "All resource outputs by resource name"
  value = {
    for name, resource in module.resources : name => {
      # Postgres outputs
      connection_url = try(resource.connection_url, null)
      host           = try(resource.host, null)
      port           = try(resource.port, null)
      database       = try(resource.database, null)

      # S3/Storage outputs
      bucket_name = try(resource.bucket_name, null)
      endpoint    = try(resource.endpoint, null)
      region      = try(resource.region, null)
    }
  }
  sensitive = true
}

# Individual resource outputs for easier CLI access
# These can be referenced as: terraform output -json | jq '.resource_names.value'
output "resource_names" {
  description = "List of all provisioned resource names"
  value       = [for r in var.resources : r.name]
}

output "project_name" {
  description = "Project name used for tagging"
  value       = var.project_name
}
