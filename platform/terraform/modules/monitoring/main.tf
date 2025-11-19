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
  method               = coalesce(each.value.method, "GET")
  expected_status_code = [200]
  interval             = 60  # 1 minute
  retry_interval       = 60  # 1 minute
  max_retries          = 3
  upside_down          = false
  ignore_tls           = false
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
  database_url         = each.value.url
  interval             = 300 # 5 minutes
  retry_interval       = 60  # 1 minute
  max_retries          = 3
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
  expected_interval    = 95040 # 26.4 hours (daily backup with 10% buffer)
  max_retries          = 2
  notification_id_list = var.notification_ids
}
