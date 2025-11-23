#
# Hetzner Server
# Provisions a VM with SSH access and firewall configuration.
#

#
# SSH Key
#
resource "tls_private_key" "this" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "hcloud_ssh_key" "this" {
  name       = "${var.name}-key"
  public_key = tls_private_key.this.public_key_openssh
}

#
# Server
#
# Ref: https://registry.terraform.io/providers/hetznercloud/hcloud/latest/docs/resources/server
#
resource "hcloud_server" "this" {
  name        = var.name
  image       = "ubuntu-24.04"
  server_type = var.server_type
  location    = var.location
  labels      = { for tag in var.tags : tag => "true" }

  # Sort to ensure deterministic ordering across Terraform runs
  # Combine personal SSH keys (passed as IDs) with the Terraform-generated key
  ssh_keys = sort(concat(
    var.ssh_key_ids,
    [hcloud_ssh_key.this.id]
  ))

  lifecycle {
    create_before_destroy = true
    ignore_changes        = [user_data]
  }

  user_data = templatefile("${path.module}/templates/server.yaml", {
    name = var.name
  })
}

#
# Firewall
#
# Ref: https://registry.terraform.io/providers/hetznercloud/hcloud/latest/docs/resources/firewall
#
resource "hcloud_firewall" "this" {
  name = "${var.name}-firewall"

  rule {
    direction = "in"
    protocol  = "tcp"
    port      = "22"
    source_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }

  rule {
    direction = "in"
    protocol  = "tcp"
    port      = "80"
    source_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }

  rule {
    direction = "in"
    protocol  = "tcp"
    port      = "443"
    source_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }

  rule {
    direction = "in"
    protocol  = "icmp"
    source_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }

  rule {
    direction = "out"
    protocol  = "tcp"
    port      = "any"
    destination_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }

  rule {
    direction = "out"
    protocol  = "udp"
    port      = "any"
    destination_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }

  rule {
    direction = "out"
    protocol  = "icmp"
    destination_ips = [
      "0.0.0.0/0",
      "::/0"
    ]
  }
}

resource "hcloud_firewall_attachment" "this" {
  firewall_id = hcloud_firewall.this.id
  server_ids  = [hcloud_server.this.id]
}
