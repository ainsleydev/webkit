# Peekaping Tag
# Tags are used to organise and categorise monitors.

resource "peekaping_tag" "this" {
  name        = var.name
  color       = var.color
  description = var.description
}
