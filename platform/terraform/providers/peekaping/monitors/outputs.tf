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

output "all_ids" {
  description = "All monitor IDs."
  value = concat(
    [for m in peekaping_monitor.http : m.id],
    [for m in peekaping_monitor.dns : m.id],
    [for m in peekaping_monitor.push : m.id]
  )
}
