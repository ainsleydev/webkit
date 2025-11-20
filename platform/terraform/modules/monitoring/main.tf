#
# Monitoring Module
# Orchestrates Peekaping resources for monitoring apps and infrastructure.
#
# This module creates:
# - Tags for organising monitors (project, environment, webkit)
# - Slack notification channel for alerts
# - HTTP monitors for app domains
# - DNS monitors for domain resolution
# - Public status page
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
    module.tag_project.id,
    module.tag_environment.id,
    module.tag_webkit.id
  ]

  # Notification IDs.
  notification_ids = local.slack_enabled ? [module.notification_slack[0].id] : []
}

#
# Tags
#
# Tags are used to organise and categorise monitors.
# Each project gets: project name tag, environment tag, and webkit tag.
#

module "tag_project" {
  source = "../../providers/peekaping/tag"

  providers = {
    peekaping = peekaping
  }

  name        = var.project_title
  color       = local.tag_color
  description = "Monitor for ${var.project_title}"
}

module "tag_environment" {
  source = "../../providers/peekaping/tag"

  providers = {
    peekaping = peekaping
  }

  name        = title(var.environment)
  color       = local.tag_color
  description = "${title(var.environment)} environment"
}

module "tag_webkit" {
  source = "../../providers/peekaping/tag"

  providers = {
    peekaping = peekaping
  }

  name        = "WebKit"
  color       = "#10B981" # Green for webkit
  description = "Managed by WebKit infrastructure"
}

#
# Slack Notification
#
# Creates a Slack notification channel for monitor alerts.
# Only created if slack_webhook_url is provided.
#

module "notification_slack" {
  count  = local.slack_enabled ? 1 : 0
  source = "../../providers/peekaping/notification"

  providers = {
    peekaping = peekaping
  }

  name        = "${var.project_title} Slack Alerts"
  type        = "slack"
  webhook_url = var.slack_webhook_url
}

#
# HTTP Monitors
#
# HTTP monitors check the availability of web applications via HTTP/HTTPS requests.
# They are created for each domain (primary + aliases) of apps with monitoring enabled.
#

module "monitor_http" {
  for_each = { for m in local.http_monitors : m.name => m }
  source   = "../../providers/peekaping/monitor"

  providers = {
    peekaping = peekaping
  }

  name = "${var.project_name}-${each.value.name}"
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

module "monitor_dns" {
  for_each = { for m in local.dns_monitors : m.name => m }
  source   = "../../providers/peekaping/monitor"

  providers = {
    peekaping = peekaping
  }

  name = "${var.project_name}-dns-${each.value.name}"
  type = "dns"
  config = jsonencode({
    hostname = each.value.domain
    resolver = "1.1.1.1" # Cloudflare DNS
    dns_type = "A"       # A record lookup
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

module "monitor_push" {
  for_each = { for m in local.push_monitors : m.name => m }
  source   = "../../providers/peekaping/monitor"

  providers = {
    peekaping = peekaping
  }

  name = "${var.project_name}-${each.value.name}"
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
# The status page is automatically populated with all monitors via tag filtering.
#

module "status_page" {
  count  = length(var.monitors) > 0 ? 1 : 0
  source = "../../providers/peekaping/status_page"

  providers = {
    peekaping = peekaping
  }

  title       = "${var.project_title} Status"
  description = "Public status page for ${var.project_title} services"
  slug        = lower(replace(var.project_name, "_", "-"))
  published   = true
  theme       = "auto"
  icon        = var.brand_logo_url
}
