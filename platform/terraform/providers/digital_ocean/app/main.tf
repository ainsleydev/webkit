# App Platform
# Ref: https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/app

# Track env var changes to trigger app replacement when they change
# This allows us to ignore DO's encryption drift while still detecting real changes
locals {
  # Create a stable, sorted map of env vars for consistent hashing
  # Only include key, value, and type (ignore scope as it doesn't affect the env var value)
  env_vars_for_hash = {
    for env in var.envs : env.key => {
      value = env.value
      type  = lookup(env, "type", "GENERAL")
    }
  }
  # Use image_name if provided, otherwise fall back to name for backward compatibility
  repository_name = var.image_name != "" ? var.image_name : var.name
}

resource "terraform_data" "env_vars_hash" {
  input = sha256(jsonencode(local.env_vars_for_hash))
}

# Track registry credentials changes to trigger app replacement when they change
# This allows us to ignore DO's encryption drift while still detecting real changes
resource "terraform_data" "registry_credentials_hash" {
  input = sha256(jsonencode({
    owner = var.github_config.owner
    token = var.github_config.token
  }))
}

resource "digitalocean_app" "this" {

  spec {
    name   = var.name
    region = var.region

    dynamic "domain" {
      for_each = var.domains
      content {
        name     = domain.value.name
        type     = domain.value.type
        zone     = lookup(domain.value, "zone", null)
        wildcard = lookup(domain.value, "wildcard", false)
      }
    }

    service {
      name               = var.service_name
      instance_size_slug = var.instance_size_slug
      instance_count     = var.instance_count
      http_port          = var.http_port

      image {
        registry_type = "GHCR"
        registry      = "ghcr.io"
        # Use the image repository name (based on repo name) to match what GitHub Actions publishes
        # GitHub Actions publishes as: ghcr.io/{owner}/{repo-name}-{app-name}
        # For example: ghcr.io/ainsleydev/player2clubs-cms
        repository = "${var.github_config.owner}/${local.repository_name}"
        tag        = var.image_tag
        # We have to use a classic token here as packages don't support fine-grained
        # PATs right now, so this should use ghp_ token formats.
        # See: https://github.com/github/roadmap/issues/558
        registry_credentials = "${var.github_config.owner}:${var.github_config.token}"
      }

      health_check {
        http_path             = var.health_check_path
        failure_threshold     = 10
        initial_delay_seconds = 90
        period_seconds        = 5
      }

      alert {
        value    = 80
        operator = "GREATER_THAN"
        window   = "FIVE_MINUTES"
        rule     = "CPU_UTILIZATION"
      }

      alert {
        value    = 80
        operator = "GREATER_THAN"
        window   = "FIVE_MINUTES"
        rule     = "MEM_UTILIZATION"
      }

      alert {
        value    = 3
        operator = "GREATER_THAN"
        window   = "FIVE_MINUTES"
        rule     = "RESTART_COUNT"
      }

      dynamic "env" {
        for_each = var.envs
        content {
          key   = env.value.key
          value = env.value.value
          type  = lookup(env.value, "type", "GENERAL")
          # Potential to make this more flexible in the future if needed.
          scope = "RUN_AND_BUILD_TIME"
        }
      }
    }
  }

  lifecycle {
    # Ignore changes to env vars and registry credentials caused by DigitalOcean's server-side encryption
    # This prevents perpetual drift between Terraform state (plain text) and DO API (encrypted)
    # See: https://github.com/digitalocean/terraform-provider-digitalocean/issues/869
    ignore_changes = [
      spec[0].service[0].env,
      spec[0].service[0].image[0].registry_credentials
    ]

    # Force app replacement when env vars or registry credentials actually change in our code
    # The terraform_data resources track hashes of var.envs and var.github_config
    # When these values change in code, the hashes change, triggering replacement
    replace_triggered_by = [
      terraform_data.env_vars_hash,
      terraform_data.registry_credentials_hash
    ]
  }
}
