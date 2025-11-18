# Hetzner Volume Module
# Creates a Hetzner Cloud volume and attaches it to a server
# Ref: https://registry.terraform.io/providers/hetznercloud/hcloud/latest/docs/resources/volume

resource "hcloud_volume" "this" {
  name     = var.name
  size     = var.size
  location = var.location
  format   = var.format
  labels   = { for tag in var.tags : tag => "true" }
}

resource "hcloud_volume_attachment" "this" {
  volume_id = hcloud_volume.this.id
  server_id = var.server_id
  automount = var.automount
}
