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
  interval       = 60  # 1 minute
  retry_interval = 60  # 1 minute
  max_retries    = 3
  upside_down    = false
  ignore_tls     = false
  # Note: accepted_status_codes and notification_id_list not supported by ehealth-co-id/uptimekuma provider
  # These need to be configured manually in Uptime Kuma UI or use ainsleydev/uptimekuma fork
}

#
# Postgres Monitors
#
# Postgres monitors check the connection health of PostgreSQL databases.
# They verify that the database is reachable and accepting connections.
#
resource "uptimekuma_monitor" "postgres" {
  for_each = { for m in local.postgres_monitors : m.name => m }

  name     = "${var.project_name}-${each.value.name}"
  type     = "postgres"
  hostname = each.value.url # Provider doesn't support database_connection_string, using hostname as workaround
  interval = 300            # 5 minutes
  retry_interval = 60       # 1 minute
  max_retries    = 3
  # Note: notification_id_list not supported by ehealth-co-id/uptimekuma provider
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
  # Note: expected_interval and notification_id_list not supported by ehealth-co-id/uptimekuma provider
  # Push monitors track heartbeat pings automatically. Notifications must be configured in Uptime Kuma UI.
}
