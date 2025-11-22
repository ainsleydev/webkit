output "http_monitors" {
  description = "HTTP monitor details."
  value       = module.monitors.http_monitors
}

output "dns_monitors" {
  description = "DNS monitor details."
  value       = module.monitors.dns_monitors
}

output "push_monitors" {
  description = "Push monitor details including ping URLs."
  value       = module.monitors.push_monitors
}

output "all_monitor_ids" {
  description = "All monitor IDs for reference."
  value       = module.monitors.all_ids
}

output "tag_ids" {
  description = "Tag IDs for reference."
  value = {
    project = module.project_tag.id
  }
}

output "status_page_url" {
  description = "Status page URL if created."
  value       = length(module.status_page) > 0 ? module.status_page[0].url : null
}
