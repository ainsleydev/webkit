#
# Peekaping Tag
# Creates a tag for organising monitors in Peekaping.
#
# Ref: https://registry.terraform.io/providers/tafaust/peekaping/latest/docs/resources/tag
#
resource "peekaping_tag" "this" {
  name        = var.name
  color       = var.colour
  description = var.description
}
