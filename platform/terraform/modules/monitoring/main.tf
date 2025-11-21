#
# Monitoring Module
# Orchestrates Peekaping resources for monitoring apps and infrastructure.
#
# This module creates:
# - Project-specific tag (unique per project)
# - Slack notification channel for alerts
# - HTTP monitors for app domains
# - DNS monitors for domain resolution
# - Public status page
#
# This module references:
# - Shared "WebKit" tag (looked up via data source)
# - Shared environment tags like "Production" (looked up via data source)
#

#
# Locals
#
locals {
  # Filter monitors by type.
  http_monitors = [for m in var.monitors : m if m.type == "http"]
  dns_monitors  = [for m in var.monitors : m if m.type == "dns"]
  push_monitors = [for m in var.monitors : m if m.type == "push"]

  # Slack notification enabled if webhook URL is provided.
  slack_enabled = var.slack_webhook_url != ""

  # Tag colour: use brand primary colour if set, otherwise fallback to blue.
  tag_color = var.brand_primary_color != null ? var.brand_primary_color : "#3B82F6"

  # All tag IDs for monitors.
  tag_ids = [
    peekaping_tag.project.id,

    # See: https://github.com/tafaust/terraform-provider-peekaping/pull/12
    # data.peekaping_tag.environment.id,
    # data.peekaping_tag.webkit.id
  ]

  # Notification IDs.
  notification_ids = local.slack_enabled ? [peekaping_notification.slack[0].id] : []
}

#
# Shared Tags (Data Sources)
#
# Look up existing shared tags that are used across multiple webkit repos.
# These tags should be created once manually or by a central terraform config.
#
# Reference: https://registry.terraform.io/providers/tafaust/peekaping/latest/docs/data-sources/tag
#

# data "peekaping_tag" "webkit" {
#   name = "WebKit"
# }
#
# data "peekaping_tag" "environment" {
#   name = title(var.environment)
# }

#
# Project Tag (Resource)
#
# Create a unique tag for this specific project.
# This tag is project-specific and won't conflict with other repos.
#
# Reference: https://registry.terraform.io/providers/tafaust/peekaping/latest/docs/resources/tag
#

resource "peekaping_tag" "project" {
  name        = var.project_title
  color       = local.tag_color
  description = "Monitor for ${var.project_title}"
}

#
# Slack Notification
#
# Creates a Slack notification channel for monitor alerts.
# Only created if slack_webhook_url is provided.
#
# Reference: https://registry.terraform.io/providers/tafaust/peekaping/latest/docs/resources/notification
#

resource "peekaping_notification" "slack" {
  count = local.slack_enabled ? 1 : 0

  name = "${var.project_title} Slack Alerts"
  type = "slack"
  config = jsonencode({
    webhook_url = var.slack_webhook_url
  })
}

#
# HTTP Monitors
#
# HTTP monitors check the availability of web applications via HTTP/HTTPS requests.
# They are created for each domain (primary + aliases) of apps with monitoring enabled.
#
# Reference: https://registry.terraform.io/providers/tafaust/peekaping/latest/docs/resources/monitor
#

resource "peekaping_monitor" "http" {
  for_each = { for m in local.http_monitors : m.name => m }

  name = each.value.name
  type = "http"
  config = jsonencode({
    url                  = each.value.url
    method               = coalesce(each.value.method, "GET")
    encoding             = "json"
    accepted_statuscodes = ["2XX"]
    authMethod           = "none"
  })

  interval         = 60 # 1 minute
  timeout          = 30 # 30 seconds
  max_retries      = 3
  retry_interval   = 60 # 1 minute
  resend_interval  = 10 # 10 minutes
  active           = true
  notification_ids = local.notification_ids
  tag_ids          = local.tag_ids
}

#
# DNS Monitors
#
# DNS monitors check domain name resolution.
# They are created for each domain to ensure DNS is correctly configured.
#
# Reference: https://registry.terraform.io/providers/tafaust/peekaping/latest/docs/resources/monitor
#

resource "peekaping_monitor" "dns" {
  for_each = { for m in local.dns_monitors : m.name => m }

  name = each.value.name
  type = "dns"
  config = jsonencode({
    host            = each.value.domain
    resolver_server = "1.1.1.1" # Cloudflare DNS
    port            = 53
    resolve_type    = "A" # A record lookup
  })

  interval         = 300 # 5 minutes (less frequent than HTTP)
  timeout          = 30  # 30 seconds
  max_retries      = 3
  retry_interval   = 60 # 1 minute
  resend_interval  = 10 # 10 minutes
  active           = true
  notification_ids = local.notification_ids
  tag_ids          = local.tag_ids
}

#
# Push Monitors (Heartbeats)
#
# Push monitors expect periodic heartbeat signals from external systems.
# They are used for backup job monitoring - the job pings the monitor on success.
#
# Reference: https://registry.terraform.io/providers/tafaust/peekaping/latest/docs/resources/monitor
#

resource "peekaping_monitor" "push" {
  for_each = { for m in local.push_monitors : m.name => m }

  name = each.value.name
  type = "push"
  config = jsonencode({
    # Push monitors don't require URL configuration
  })

  max_retries      = 2
  retry_interval   = 60
  resend_interval  = 10
  active           = true
  notification_ids = local.notification_ids
  tag_ids          = local.tag_ids
}

#
# Status Page
#
# Creates a public status page showing the health of all monitors.
# All HTTP, DNS, and push monitors are explicitly attached to the status page.
# A custom domain (e.g., status.example.com) can be configured for public access.
#
# Reference: https://registry.terraform.io/providers/tafaust/peekaping/latest/docs/resources/status_page
#

resource "peekaping_status_page" "this" {
  count = length(var.monitors) > 0 ? 1 : 0

  title       = "${var.project_title} Status"
  description = "Public status page for ${var.project_title} services"
  slug        = lower(replace(var.project_name, "_", "-"))
  published   = true
  theme       = "auto"
  icon        = var.brand_logo_url
  domains     = var.status_page_domain != null ? [var.status_page_domain] : []
  monitor_ids = concat(
    [for m in peekaping_monitor.http : m.id],
    [for m in peekaping_monitor.dns : m.id],
    [for m in peekaping_monitor.push : m.id]
  )
}
