# Peekaping Notification
# Creates a notification channel for monitor alerts.

resource "peekaping_notification" "this" {
  name = var.name
  type = var.type
  config = jsonencode({
    webhook_url = var.webhook_url
  })
}
