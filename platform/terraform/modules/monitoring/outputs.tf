output "monitors" {
  description = "All monitors as a flat array with type field."
  value       = module.monitors.monitors
}

output "all_monitor_ids" {
  description = "All monitor IDs for reference."
  value       = module.monitors.all_ids
}

output "peekaping" {
  description = "Peekaping configuration and outputs."
  value = {
    project_tag = module.project_tag.id
  }
}
