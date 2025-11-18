output "http_monitors" {
  description = "HTTP monitor details"
  value = {
    for name, monitor in uptimekuma_monitor.http : name => {
      id   = monitor.id
      name = monitor.name
      url  = monitor.url
    }
  }
}

output "postgres_monitors" {
  description = "Postgres monitor details"
  value = {
    for name, monitor in uptimekuma_monitor.postgres : name => {
      id   = monitor.id
      name = monitor.name
    }
  }
  sensitive = true # Contains database connection info.
}

output "push_monitors" {
  description = "Push monitor details (for heartbeats)"
  value = {
    for name, monitor in uptimekuma_monitor.push : name => {
      id       = monitor.id
      name     = monitor.name
      push_url = try(monitor.push_url, "")
    }
  }
}

output "all_monitor_ids" {
  description = "All monitor IDs for reference"
  value = concat(
    [for m in uptimekuma_monitor.http : m.id],
    [for m in uptimekuma_monitor.postgres : m.id],
    [for m in uptimekuma_monitor.push : m.id]
  )
}
