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

# TODO: Re-enable when resource monitoring is implemented
# output "postgres_monitors" - Removed, postgres monitoring disabled
# output "push_monitors" - Removed, push monitoring disabled

output "all_monitor_ids" {
  description = "All monitor IDs for reference"
  value       = [for m in uptimekuma_monitor.http : m.id]
}
