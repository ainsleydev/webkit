variable "name" {
  description = "Notification channel name"
  type        = string
}

variable "type" {
  description = "Notification type (slack, email, etc.)"
  type        = string
}

variable "webhook_url" {
  description = "Webhook URL for Slack notifications"
  type        = string
}
