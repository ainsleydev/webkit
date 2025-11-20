output "http_monitors" {
  description = "HTTP monitor details"
  value = {
    for name, monitor in module.monitor_http : name => {
      id   = monitor.id
      name = monitor.name
    }
  }
}

output "dns_monitors" {
  description = "DNS monitor details"
  value = {
    for name, monitor in module.monitor_dns : name => {
      id   = monitor.id
      name = monitor.name
    }
  }
}

output "push_monitors" {
  description = "Push monitor details"
  value = {
    for name, monitor in module.monitor_push : name => {
      id   = monitor.id
      name = monitor.name
    }
  }
}

output "all_monitor_ids" {
  description = "All monitor IDs for reference"
  value = concat(
    [for m in module.monitor_http : m.id],
    [for m in module.monitor_dns : m.id],
    [for m in module.monitor_push : m.id]
  )
}

output "tag_ids" {
  description = "Tag IDs for reference"
  value = {
    project     = module.tag_project.id
    environment = module.tag_environment.id
    webkit      = module.tag_webkit.id
  }
}

output "status_page_url" {
  description = "Status page URL if created"
  value       = length(module.status_page) > 0 ? module.status_page[0].url : null
}

output "slack_notification_id" {
  description = "Slack notification ID if created"
  value       = length(module.notification_slack) > 0 ? module.notification_slack[0].id : null
}
