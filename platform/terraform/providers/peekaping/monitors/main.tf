#
# Peekaping Monitors
# Creates HTTP, DNS, and Push monitors in Peekaping.
#
# Ref: https://registry.terraform.io/providers/tafaust/peekaping/latest/docs/resources/monitor
#

#
# Locals
#
locals {
  # Filter monitors by type.
  http_monitors = [for m in var.monitors : m if m.type == "http"]
  dns_monitors  = [for m in var.monitors : m if m.type == "dns"]
  push_monitors = [for m in var.monitors : m if m.type == "push"]

  # Default values with optional overrides.
  defaults = {
    timeout          = coalesce(var.defaults.timeout, 30)
    http_max_retries = coalesce(var.defaults.http_max_retries, 3)
    dns_max_retries  = coalesce(var.defaults.dns_max_retries, 3)
    push_max_retries = coalesce(var.defaults.push_max_retries, 2)
    retry_interval   = coalesce(var.defaults.retry_interval, 60)
    resend_interval  = coalesce(var.defaults.resend_interval, 10)
  }
}

#
# HTTP Monitors
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
    max_redirects        = coalesce(each.value.max_redirects, 0)
  })

  interval         = each.value.interval
  timeout          = local.defaults.timeout
  max_retries      = local.defaults.http_max_retries
  retry_interval   = local.defaults.retry_interval
  resend_interval  = local.defaults.resend_interval
  active           = true
  notification_ids = var.notification_ids
  tag_ids          = var.tag_ids
}

#
# DNS Monitors
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
  notification_ids = var.notification_ids
  tag_ids          = var.tag_ids
}

#
# Push Token Generation
# Generates deterministic push tokens for monitors.
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
resource "peekaping_monitor" "push" {
  for_each = { for m in local.push_monitors : m.name => m }

  name       = each.value.name
  type       = "push"
  push_token = random_id.push_token[each.key].b64_url
  config = jsonencode({
    pushToken = random_id.push_token[each.key].b64_url
  })

  interval         = each.value.interval
  timeout          = local.defaults.timeout
  max_retries      = local.defaults.push_max_retries
  retry_interval   = local.defaults.retry_interval
  resend_interval  = local.defaults.resend_interval
  active           = true
  notification_ids = var.notification_ids
  tag_ids          = var.tag_ids
}
