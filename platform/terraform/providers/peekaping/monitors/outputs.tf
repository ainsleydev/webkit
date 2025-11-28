output "monitors" {
  description = "All monitors as a flat array with type field. Push monitors include extra fields for CI/CD."
  # IMPORTANT: When adding new monitor types in main.tf, you MUST update this output
  # to include them. Each monitor type (http, http_keyword, dns, push) must be listed.
  value = concat(
    [for name, monitor in peekaping_monitor.http : {
      id   = monitor.id
      name = monitor.name
      type = "http"
    }],
    [for name, monitor in peekaping_monitor.http_keyword : {
      id   = monitor.id
      name = monitor.name
      type = "http-keyword"
    }],
    [for name, monitor in peekaping_monitor.dns : {
      id   = monitor.id
      name = monitor.name
      type = "dns"
    }],
    [for name, monitor in peekaping_monitor.push : {
      id            = monitor.id
      name          = monitor.name
      type          = "push"
      variable_name = local.push_monitors_map[name].variable_name
      push_token    = random_id.push_token[name].b64_url
      ping_url      = "${var.peekaping_endpoint}/api/v1/push/${random_id.push_token[name].b64_url}?status=up&msg=OK&ping="
    }]
  )
}

output "all_ids" {
  description = "All monitor IDs."
  # IMPORTANT: When adding new monitor types in main.tf, you MUST update this output
  # to include them. Each monitor type must be listed to appear on the status page.
  value = concat(
    [for m in peekaping_monitor.http : m.id],
    [for m in peekaping_monitor.http_keyword : m.id],
    [for m in peekaping_monitor.dns : m.id],
    [for m in peekaping_monitor.push : m.id]
  )
}
