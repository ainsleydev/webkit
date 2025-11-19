#
# Monitoring Module
# Generates module calls based on monitoring[] in resources and apps.
#
# TODO: Re-enable when resource monitoring is implemented
# Postgres and Push monitors have been removed to simplify initial implementation.
#
#
locals {
  http_monitors = [for m in var.monitors : m if m.type == "http"]
  push_monitors = [for m in var.monitors : m if m.type == "push"]
}

#
# HTTP Monitors
#
# HTTP monitors check the availability of web applications via HTTP/HTTPS requests.
# They are created for each domain (primary + aliases) of apps with monitoring enabled.
#
resource "uptimekuma_monitor" "http" {
  for_each = { for m in local.http_monitors : m.name => m }

  name           = "${var.project_name}-${each.value.name}"
  type           = "http"
  url            = each.value.url
  method         = coalesce(each.value.method, "GET")
  interval       = 60 # 1 minute
  retry_interval = 60 # 1 minute
  max_retries    = 3
}

#
# Push Monitors (Heartbeats)
#
# Push monitors expect periodic heartbeat signals from external systems.
# They are used for backup job monitoring - the job pings the monitor on success.
#
resource "uptimekuma_monitor" "push" {
  for_each = { for m in local.push_monitors : m.name => m }

  name        = "${var.project_name}-${each.value.name}"
  type        = "push"
  max_retries = 2
}

