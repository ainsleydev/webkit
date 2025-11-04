# App Platform
# Ref: https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/app
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

    alert { rule = "DEPLOYMENT_FAILED" }
    alert { rule = "DEPLOYMENT_LIVE" }

    service {
      name               = var.service_name
      instance_size_slug = var.instance_size_slug
      instance_count     = var.instance_count
      http_port          = var.http_port

      image {
        registry_type        = "GHCR"
        registry             = "ghcr.io"
        # The var.name variable should match the name of the image on GHCR.
        # For example: ainsleydev/search-spares-web
        repository           = "${var.github_config.owner}/${var.name}"
        tag                  = "sha-2eef5657ea0ae1e440512e352b82c2a0d0dfa456"
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
}
