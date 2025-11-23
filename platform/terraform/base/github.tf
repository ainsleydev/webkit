#
# GitHub Secrets and Variables
# Exposes Terraform outputs to GitHub Actions for CI/CD workflows.
#

#
# Resource and App Secrets
# Created after resources/apps are provisioned.
#
locals {
  # Output mappings for each resource type.
  resource_output_map = {
    postgres = ["id", "urn", "connection_url"]
    s3       = ["id", "urn", "bucket_name", "bucket_url", "region", "endpoint"]
    sqlite   = ["id", "connection_url", "auth_token", "host", "database"]
  }

  # Output mappings for each app type.
  app_output_map = {
    vm        = ["ip_address", "ssh_private_key", "server_user"]
    container = []
  }

  # Build secret keys from var.resources.
  github_secrets_resources = merge([
    for resource in var.resources : {
      for output_name in lookup(local.resource_output_map, resource.platform_type, []) :
      upper("TF_${local.environment_short}_${replace(resource.name, "-", "_")}_${output_name}") => tomap({
        source_type   = "resource"
        resource_name = resource.name
        output_name   = output_name
      })
    }
  ]...)

  # Build secret keys from var.apps.
  github_secrets_apps = merge([
    for app in var.apps : {
      for output_name in lookup(local.app_output_map, app.platform_type, []) :
      upper("TF_${local.environment_short}_${replace(app.name, "-", "_")}_${output_name}") => tomap({
        source_type = "app"
        app_name    = app.name
        output_name = output_name
      })
    }
  ]...)

  github_secrets = merge(
    local.github_secrets_resources,
    local.github_secrets_apps
  )
}

resource "github_actions_secret" "resource_outputs" {
  for_each = local.github_secrets

  repository  = var.github_config.repo
  secret_name = each.key
  plaintext_value = try(
    each.value["source_type"] == "resource" ? tostring(module.resources[each.value["resource_name"]][each.value["output_name"]]) :
    each.value["source_type"] == "app" ? tostring(module.apps[each.value["app_name"]][each.value["output_name"]]) :
    "NOT_SET"
  )
  depends_on = [module.resources, module.apps, module.monitoring]
}

#
# Monitor Ping URL Variables
# Stored as variables (not secrets) for easier debugging.
#
# Naming: {ENV}_{IDENTIFIER}_{TYPE}_PING_URL
# Example: PROD_DB_BACKUP_PING_URL
#
resource "github_actions_variable" "monitor_ping_urls" {
  # Only create variables for push monitors that have a variable_name computed by Go.
  for_each = length(var.monitors) > 0 ? {
    for k, v in module.monitoring[0].push_monitors : k => v if v.variable_name != null && v.variable_name != ""
  } : {}

  repository    = var.github_config.repo
  variable_name = each.value.variable_name
  value         = each.value.ping_url
}
