#
# Apps Module
# Generates module calls based on apps[] in app.json
# Maps generic infra types to provider-specific resources
#

locals {
  # Determine which provider module to use
  is_do_vm        = var.cloud_provider == "digitalocean" && var.infra_type == "vm"
  is_do_container = var.cloud_provider == "digitalocean" && var.infra_type == "container"

  # Extract domain from config for DNS record
  domain = try(var.infra_config.domain, null)
}

# DigitalOcean Droplet (VM)
module "do_droplet" {
  count  = local.is_do_vm ? 1 : 0
  source = "../../providers/digital_ocean/droplet"

  name              = "${var.project_name}-${var.name}"
  user_ssh_key_name = try(var.infra_config.ssh_keys[0], var.user_ssh_key_name)
  droplet_size      = try(var.infra_config.size, "s-1vcpu-1gb")
  droplet_region    = try(var.infra_config.region, "lon1")
  tags              = try(var.tags, [])
}

# DigitalOcean App Platform (Container)
module "do_app" {
  count  = local.is_do_container ? 1 : 0
  source = "../../providers/digital_ocean/app"

  name               = "${var.project_name}-${var.name}"
  service_name       = var.name
  region             = try(var.infra_config.region, "lon")
  instance_size_slug = try(var.infra_config.size, "apps-s-1vcpu-1gb")
  instance_count     = try(var.infra_config.instance_count, 1)
  http_port          = try(var.infra_config.port, 3000)
  image_tag          = var.image_tag
  github_config      = var.github_config
  health_check_path  = try(var.infra_config.health_check_path, "/")

  # Configure domain if provided
  domains = local.domain != null ? [
    {
      name = local.domain
      type = "PRIMARY"
    }
  ] : []

  # Environment variables will be passed through
  envs = var.env_vars
}

# DNS Record for VM-based apps
# App Platform handles its own DNS, but VMs need A records
module "dns_record" {
  count  = local.is_do_vm && local.domain != null ? 1 : 0
  source = "../../providers/digital_ocean/domain_record"

  # Extract root domain and subdomain from FQDN
  # e.g., "cms.my-website.com" -> domain: "my-website.com", name: "cms"
  domain = join(".", slice(split(".", local.domain), length(split(".", local.domain)) - 2, length(split(".", local.domain))))
  name   = join(".", slice(split(".", local.domain), 0, length(split(".", local.domain)) - 2))
  type   = "A"
  value  = module.do_droplet[0].ip_address
}
