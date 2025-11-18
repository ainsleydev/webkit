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

    service {
      name               = var.service_name
      instance_size_slug = var.instance_size_slug
      instance_count     = var.instance_count
      http_port          = var.http_port

      image {
        registry_type        = "GHCR"
        registry             = "ghcr.io"
        repository           = var.repository
        tag                  = var.image_tag
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
        disabled = var.notifications_webhook_url == ""

        destinations {
          slack_webhooks {
            channel = var.slack_channel_name
            url     = var.notifications_webhook_url
          }
        }
      }

      alert {
        value    = 80
        operator = "GREATER_THAN"
        window   = "FIVE_MINUTES"
        rule     = "MEM_UTILIZATION"
        disabled = var.notifications_webhook_url == ""

        destinations {
          slack_webhooks {
            channel = var.slack_channel_name
            url     = var.notifications_webhook_url
          }
        }
      }

      alert {
        value    = 3
        operator = "GREATER_THAN"
        window   = "FIVE_MINUTES"
        rule     = "RESTART_COUNT"
        disabled = var.notifications_webhook_url == ""

        destinations {
          slack_webhooks {
            channel = var.slack_channel_name
            url     = var.notifications_webhook_url
          }
        }
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
