# Peekaping Monitor
# Creates a monitor for checking service availability.

resource "peekaping_monitor" "this" {
  name             = var.name
  type             = var.type
  config           = var.config
  interval         = var.interval
  timeout          = var.timeout
  max_retries      = var.max_retries
  retry_interval   = var.retry_interval
  resend_interval  = var.resend_interval
  active           = var.active
  notification_ids = var.notification_ids
  tag_ids          = var.tag_ids
}
