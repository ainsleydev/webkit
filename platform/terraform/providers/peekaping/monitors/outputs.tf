output "monitors" {
  description = "All monitors as a flat array with type field."
  value = concat(
    [for name, monitor in peekaping_monitor.http : {
      id   = monitor.id
      name = monitor.name
      type = "http"
    }],
    [for name, monitor in peekaping_monitor.dns : {
      id   = monitor.id
      name = monitor.name
      type = "dns"
    }],
    [for name, monitor in peekaping_monitor.push : {
      id   = monitor.id
      name = monitor.name
      type = "push"
    }]
  )
}

output "push_monitors" {
  description = "Push monitor details including ping URLs for GitHub variables."
  value = {
    for name, monitor in peekaping_monitor.push : name => {
      id            = monitor.id
      name          = monitor.name
      variable_name = local.push_monitors_map[name].variable_name
      push_token    = random_id.push_token[name].b64_url
      ping_url      = "${var.peekaping_endpoint}/api/v1/push/${random_id.push_token[name].b64_url}?status=up&msg=OK&ping="
    }
  }
}

output "all_ids" {
  description = "All monitor IDs."
  value = concat(
    [for m in peekaping_monitor.http : m.id],
    [for m in peekaping_monitor.dns : m.id],
    [for m in peekaping_monitor.push : m.id]
  )
}
