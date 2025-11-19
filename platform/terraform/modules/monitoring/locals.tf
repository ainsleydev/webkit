locals {
  # Filter monitors by type for easier resource creation.
  # Currently only HTTP monitors are supported.
  http_monitors = [for m in var.monitors : m if m.type == "http"]

  # TODO: Re-enable when resource monitoring is implemented
  postgres_monitors = [for m in var.monitors : m if m.type == "postgres"]
  push_monitors     = [for m in var.monitors : m if m.type == "push"]
}
