#
# App Outputs
# Returns relevant outputs based on the infrastructure type
#

# VM (Droplet) Outputs
output "ip_address" {
  description = "IP address of the VM"
  value = (
    var.infra_type == "vm" && var.cloud_provider == "digitalocean" ? module.do_droplet[0].ip_address :
    null
  )
}

output "droplet_id" {
  description = "ID of the droplet"
  value = (
    var.infra_type == "vm" && var.cloud_provider == "digitalocean" ? module.do_droplet[0].id :
    null
  )
}

output "ssh_private_key" {
  description = "SSH private key for the VM"
  value = (
    var.infra_type == "vm" && var.cloud_provider == "digitalocean" ? module.do_droplet[0].ssh_private_key :
    null
  )
  sensitive = true
}

# Container (App Platform) Outputs
output "app_id" {
  description = "ID of the App Platform app"
  value = (
    var.infra_type == "container" && var.cloud_provider == "digitalocean" ? module.do_app[0].id :
    null
  )
}

output "app_url" {
  description = "Live URL of the app"
  value = (
    var.infra_type == "container" && var.cloud_provider == "digitalocean" ? module.do_app[0].live_url :
    null
  )
}

output "app_domain" {
  description = "Live domain of the app"
  value = (
    var.infra_type == "container" && var.cloud_provider == "digitalocean" ? module.do_app[0].live_domain :
    null
  )
}

# DNS Outputs
output "dns_fqdn" {
  description = "Fully qualified domain name"
  value = (
    var.infra_type == "vm" && length(module.dns_record) > 0 ? module.dns_record[0].fqdn :
    null
  )
}

# Common Outputs
output "infra_type" {
  description = "Type of infrastructure provisioned"
  value       = var.infra_type
}

output "cloud_provider" {
  description = "Cloud provider used"
  value       = var.cloud_provider
}
