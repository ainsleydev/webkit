# Lookup all existing SSH keys from DigitalOcean based on the names provided
data "digitalocean_ssh_key" "this" {
  for_each = toset(var.ssh_keys)
  name     = each.value
}

resource "tls_private_key" "this" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "digitalocean_ssh_key" "this" {
  name       = "${var.name}-key"
  public_key = tls_private_key.this.public_key_openssh
}

# Droplet
# Ref: https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/droplet
resource "digitalocean_droplet" "this" {
  name   = var.name
  image  = "ubuntu-24-04-x64"
  size   = var.droplet_size
  region = var.droplet_region
  tags   = var.tags

  # Sort to ensure deterministic ordering across Terraform runs
  # Without sort(), for_each iteration order is not guaranteed
  ssh_keys = sort(concat(
    [for k in data.digitalocean_ssh_key.this : k.id],
    [digitalocean_ssh_key.this.id]
  ))

  lifecycle {
    create_before_destroy = true
    ignore_changes        = [user_data, ssh_keys]
  }

  user_data = templatefile("${path.module}/templates/server.yaml", {
    name = var.name
  })
}

# Firewall
# Ref: https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/firewall
resource "digitalocean_firewall" "this" {
  name        = "${var.name}-firewall"
  droplet_ids = [digitalocean_droplet.this.id]

  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "80"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "443"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "icmp"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "tcp"
    port_range            = "all"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
}
