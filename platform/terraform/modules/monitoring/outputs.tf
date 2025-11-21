output "http_monitors" {
  description = "HTTP monitor details"
  value = {
    for name, monitor in peekaping_monitor.http : name => {
      id   = monitor.id
      name = monitor.name
    }
  }
}

output "dns_monitors" {
  description = "DNS monitor details"
  value = {
    for name, monitor in peekaping_monitor.dns : name => {
      id   = monitor.id
      name = monitor.name
    }
  }
}

output "push_monitors" {
  description = "Push monitor details including ping URLs"
  value = {
    for name, monitor in peekaping_monitor.push : name => {
      id         = monitor.id
      name       = monitor.name
      push_token = monitor.push_token
      ping_url   = monitor.push_token != null ? "${var.peekaping_endpoint}/api/push/${monitor.push_token}" : null
    }
  }
}

output "all_monitor_ids" {
  description = "All monitor IDs for reference"
  value = concat(
    [for m in peekaping_monitor.http : m.id],
    [for m in peekaping_monitor.dns : m.id],
    [for m in peekaping_monitor.push : m.id]
  )
}

output "tag_ids" {
  description = "Tag IDs for reference"
  value = {
    project     = peekaping_tag.project.id
    # environment = data.peekaping_tag.environment.id
    # webkit      = data.peekaping_tag.webkit.id
  }
}

output "status_page_url" {
  description = "Status page URL if created"
  value       = length(peekaping_status_page.this) > 0 ? "https://${peekaping_status_page.this[0].slug}" : null
}

# output "slack_notification_id" {
#   description = "Slack notification ID if created"
#   value       = length(peekaping_notification.slack) > 0 ? peekaping_notification.slack[0].id : null
# }
