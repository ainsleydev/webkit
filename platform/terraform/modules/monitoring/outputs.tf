output "monitors" {
  description = "All monitors as a flat array with type field."
  value       = module.monitors.monitors
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
