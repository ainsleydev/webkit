#
# App Outputs
# Returns relevant outputs based on the infrastructure type
#

# VM (Droplet) Outputs
output "ip_address" {
  description = "IP address of the VM"
  value = (
    var.platform_type == "vm" && var.platform_provider == "digitalocean" ? module.do_droplet[0].ip_address :
    var.platform_type == "vm" && var.platform_provider == "hetzner" ? module.hetzner_server[0].ip_address :
    null
  )
}

output "droplet_id" {
  description = "ID of the droplet"
  value = (
    var.platform_type == "vm" && var.platform_provider == "digitalocean" ? module.do_droplet[0].id :
    null
  )
}

output "ssh_private_key" {
  description = "SSH private key for the VM"
  value = (
    var.platform_type == "vm" && var.platform_provider == "digitalocean" ? module.do_droplet[0].ssh_private_key :
    var.platform_type == "vm" && var.platform_provider == "hetzner" ? module.hetzner_server[0].ssh_private_key :
    null
  )
  sensitive = true
}

output "server_user" {
  description = "SSH user for the VM"
  value = (
    var.platform_type == "vm" && var.platform_provider == "digitalocean" ? module.do_droplet[0].server_user :
    var.platform_type == "vm" && var.platform_provider == "hetzner" ? module.hetzner_server[0].server_user :
    null
  )
}

# Container (App Platform) Outputs
output "app_id" {
  description = "ID of the App Platform app"
  value = (
    var.platform_type == "container" && var.platform_provider == "digitalocean" ? module.do_app[0].id :
    null
  )
}

output "app_url" {
  description = "Live URL of the app"
  value = (
    var.platform_type == "container" && var.platform_provider == "digitalocean" ? module.do_app[0].live_url :
    null
  )
}

output "app_domain" {
  description = "Live domain of the app"
  value = (
    var.platform_type == "container" && var.platform_provider == "digitalocean" ? module.do_app[0].live_domain :
    null
  )
}

# DNS Outputs
# output "dns_fqdn" {
#   description = "Fully qualified domain name"
#   value = (
#     var.platform_type == "vm" && length(module.dns_record) > 0 ? module.dns_record[0].fqdn :
#     null
#   )
# }

# Common Outputs
output "platform_type" {
  description = "Type of infrastructure provisioned"
  value       = var.platform_type
}

output "platform_provider" {
  description = "Cloud provider used"
  value       = var.platform_provider
}

output "urn" {
  description = "Resource URN (DigitalOcean specific)"
  value = (
    var.platform_type == "vm" && var.platform_provider == "digitalocean" ? module.do_droplet[0].urn :
    var.platform_type == "container" && var.platform_provider == "digitalocean" ? module.do_app[0].urn :
    null
  )
}
