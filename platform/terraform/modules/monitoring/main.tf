#
# HTTP Monitors
#
# HTTP monitors check the availability of web applications via HTTP/HTTPS requests.
# They are created for each domain (primary + aliases) of apps with monitoring enabled.
#
resource "uptimekuma_monitor" "http" {
  for_each = { for m in local.http_monitors : m.name => m }

  name                 = "${var.project_name}-${each.value.name}"
  type                 = "http"
  url                  = each.value.url
  method               = each.value.method
  expected_status_code = each.value.expected_status
  interval             = each.value.interval
  retry_interval       = each.value.retry_interval
  max_retries          = each.value.max_retries
  upside_down          = each.value.upside_down
  ignore_tls           = each.value.ignore_tls
  notification_id_list = var.notification_ids
}

#
# Postgres Monitors
#
# Postgres monitors check the connection health of PostgreSQL databases.
# They verify that the database is reachable and accepting connections.
#
resource "uptimekuma_monitor" "postgres" {
  for_each = { for m in local.postgres_monitors : m.name => m }

  name                 = "${var.project_name}-${each.value.name}"
  type                 = "postgres"
  database_url         = each.value.database_url
  interval             = each.value.interval
  retry_interval       = each.value.retry_interval
  max_retries          = each.value.max_retries
  notification_id_list = var.notification_ids
}

#
# Push Monitors (Heartbeats)
#
# Push monitors expect periodic heartbeat signals from external systems.
# They are used for backup job monitoring - the job pings the monitor on success.
#
resource "uptimekuma_monitor" "push" {
  for_each = { for m in local.push_monitors : m.name => m }

  name                 = "${var.project_name}-${each.value.name}"
  type                 = "push"
  expected_interval    = each.value.expected_interval
  max_retries          = each.value.max_retries
  notification_id_list = var.notification_ids
}
