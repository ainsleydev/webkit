#
# Monitoring Module
#
# Orchestrates Peekaping resources for monitoring apps and infrastructure.
#

#
# Locals
#
locals {
  # Filter monitors by type.
  http_monitors = [for m in var.monitors : m if m.type == "http"]
  dns_monitors  = [for m in var.monitors : m if m.type == "dns"]
  push_monitors = [for m in var.monitors : m if m.type == "push"]

  # Tag colour: use brand primary colour if set, otherwise fallback to blue.
  tag_color        = var.brand_primary_color != null ? var.brand_primary_color : "#3B82F6"

  # Currently - Slack - #alerts
  notification_ids = ["7e4f8d2e-5720-4b07-9ce4-3f639b5e4647"]

  # Default values
  defaults = {
    timeout         = 30
    http_max_retries = 3
    dns_max_retries  = 3
    push_max_retries = 2
    retry_interval   = 60
    resend_interval  = 10
  }

  # Tag IDs
  tag_ids = [
    peekaping_tag.project.id,
    "ac5b2626-3425-4496-a318-ede51ce7baa8",
    "dd1151bc-1d7a-42d1-8166-f87b7b180798",
  ]
}

#
# Tags (Data Sources)
#
# Reference: https://registry.terraform.io/providers/tafaust/peekaping/latest/docs/data-sources/tag
#
resource "peekaping_tag" "project" {
  name        = var.project_title
  color       = local.tag_color
  description = "Monitor for ${var.project_title}"
}

#
# Slack Notification
#
# Reference: https://registry.terraform.io/providers/tafaust/peekaping/latest/docs/resources/notification
#

#
# HTTP Monitors
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

  interval         = each.value.interval
  timeout          = local.defaults.timeout
  max_retries      = local.defaults.http_max_retries
  retry_interval   = local.defaults.retry_interval
  resend_interval  = local.defaults.resend_interval
  active           = true
  notification_ids = local.notification_ids
  tag_ids          = local.tag_ids

  depends_on = [peekaping_tag.project]
}

#
# DNS Monitors
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

  interval         = each.value.interval
  timeout          = local.defaults.timeout
  max_retries      = local.defaults.dns_max_retries
  retry_interval   = local.defaults.retry_interval
  resend_interval  = local.defaults.resend_interval
  active           = true
  notification_ids = local.notification_ids
  tag_ids          = local.tag_ids

  depends_on = [peekaping_tag.project]
}

#
# Push Token Generation
#
# Generates deterministic push tokens for monitors.
# Tokens only change when the monitor name changes.
#
resource "random_id" "push_token" {
  for_each = { for m in local.push_monitors : m.name => m }

  byte_length = 24 # 24 bytes = 32 characters in base64

  keepers = {
    monitor_name = each.value.name
  }
}

#
# Push Monitors (Heartbeats)
#
# Reference: https://registry.terraform.io/providers/tafaust/peekaping/latest/docs/resources/monitor
#
resource "peekaping_monitor" "push" {
  for_each = { for m in local.push_monitors : m.name => m }

  name = each.value.name
  type = "push"
  config = jsonencode({
    pushToken = random_id.push_token[each.key].b64_url
  })

  interval         = each.value.interval
  timeout          = local.defaults.timeout
  max_retries      = local.defaults.push_max_retries
  retry_interval   = local.defaults.retry_interval
  resend_interval  = local.defaults.resend_interval
  active           = true
  notification_ids = local.notification_ids
  tag_ids          = local.tag_ids

  depends_on = [peekaping_tag.project]
}

#
# Status Page
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
  icon        = var.brand_icon_url
  domains     = var.status_page_domain != null ? [var.status_page_domain] : []
  monitor_ids = concat(
    [for m in peekaping_monitor.http : m.id],
    [for m in peekaping_monitor.dns : m.id],
    [for m in peekaping_monitor.push : m.id]
  )
}
