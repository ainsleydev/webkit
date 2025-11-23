output "http_monitors" {
  description = "HTTP monitor details."
  value = {
    for name, monitor in peekaping_monitor.http : name => {
      id   = monitor.id
      name = monitor.name
    }
  }
}

output "dns_monitors" {
  description = "DNS monitor details."
  value = {
    for name, monitor in peekaping_monitor.dns : name => {
      id   = monitor.id
      name = monitor.name
    }
  }
}

output "push_monitors" {
  description = "Push monitor details including ping URLs."
  value = {
    for name, monitor in peekaping_monitor.push : name => {
      id         = monitor.id
      name       = monitor.name
      identifier = local.push_monitors_map[name].identifier
      push_token = random_id.push_token[name].b64_url
      ping_url   = "${var.peekaping_endpoint}/api/v1/push/${random_id.push_token[name].b64_url}?status=up&msg=OK&ping="
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
