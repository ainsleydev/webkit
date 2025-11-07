#
# Apps Module
# Generates module calls based on apps[] in app.json
# Maps generic infra types to provider-specific resources
#

# DigitalOcean Droplet (VM)
module "do_droplet" {
  count  = var.platform_provider == "digitalocean" && var.platform_type == "vm" ? 1 : 0
  source = "../../providers/digital_ocean/droplet"

  name           = "${var.project_name}-${var.name}"
  droplet_size   = try(var.platform_config.size, "s-1vcpu-1gb")
  droplet_region = try(var.platform_config.region, "lon1")
  ssh_key_ids    = var.do_ssh_key_ids
  tags           = try(var.tags, [])
  server_user    = var.server_user
}

# DigitalOcean App Platform (Container)
module "do_app" {
  count  = var.platform_provider == "digitalocean" && var.platform_type == "container" ? 1 : 0
  source = "../../providers/digital_ocean/app"

  project_name       = var.project_name
  name               = var.name
  service_name       = var.app_type
  region             = try(var.platform_config.region, "lon")
  instance_size_slug = try(var.platform_config.size, "apps-s-1vcpu-1gb")
  instance_count     = try(var.platform_config.instance_count, 1)
  http_port          = try(var.platform_config.port, 3000)
  image_tag          = var.image_tag
  github_config      = var.github_config
  health_check_path  = try(var.platform_config.health_check_path, "/")
  notifications_webhook_url = var.notifications_webhook_url

  envs = [
    for env in var.env_vars : {
      key = env.key
      value = (
        startswith(env.value, "resource:")
        ? var.resource_outputs[split(".", trimprefix(env.value, "resource:"))[0]][split(".", trimprefix(env.value, "resource:"))[1]]
        : env.value
      )
      scope = try(env.scope, "RUN_TIME")
      type  = try(env.type, "GENERAL")
    }
  ]

  domains = [
    for d in var.domains : {
      name     = d.name
      type     = upper(d.type) # PRIMARY, ALIAS, UNMANAGED
      zone     = d.zone
      wildcard = d.wildcard
    }
  ]
}

# DNS Record for VM-based apps
# App Platform handles its own DNS, but VMs need A records
# module "dns_record" {
#   count  = local.is_do_vm && local.domain != null ? 1 : 0
#   source = "../../providers/digital_ocean/domain_record"
#
#   # Extract root domain and subdomain from FQDN
#   # e.g., "cms.my-website.com" -> domain: "my-website.com", name: "cms"
#   domain = join(".", slice(split(".", local.domain), length(split(".", local.domain)) - 2, length(split(".", local.domain))))
#   name   = join(".", slice(split(".", local.domain), 0, length(split(".", local.domain)) - 2))
#   type   = "A"
#   value  = module.do_droplet[0].ip_address
# }
