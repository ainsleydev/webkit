locals {
  # Filter monitors by type for easier resource creation.
  http_monitors     = [for m in var.monitors : m if m.type == "http"]
  postgres_monitors = [for m in var.monitors : m if m.type == "postgres"]
  push_monitors     = [for m in var.monitors : m if m.type == "push"]
}
