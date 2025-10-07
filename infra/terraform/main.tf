terraform {
  required_version = ">= 1.13.0"

  # Backend configuration created on CI (backend.hcl)

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
    b2 = {
      source  = "Backblaze/b2"
      version = "~> 0.10.0"
    }
  }
}

#
# Providers
#
# Config will be set via environment variables:
# - DIGITALOCEAN_TOKEN
# - B2_APPLICATION_KEY_ID
# - B2_APPLICATION_KEY
#
provider "digitalocean" {}
provider "b2" {}

#
# Locals
#
locals {
  # Default tags applied to all resources (normalized to lowercase)
  default_tags = [
    lower(var.project_name),
    lower(var.environment),
    "terraform",
  ]

  # Combined tags: default + custom (all normalized to lowercase)
  common_tags = concat(
    local.default_tags,
    [for tag in var.tags : lower(tag)]
  )
}

#
# Resources (databases, storage, etc.)
#
module "resources" {
  for_each = { for r in var.resources : r.name => r }
  source   = "./modules/resources"

  project_name      = var.project_name
  name              = each.value.name
  platform_type     = each.value.type
  platform_provider = each.value.provider
  platform_config   = each.value.config
  tags              = local.common_tags
}

#
# Apps (services, applications)
#
module "apps" {
  for_each = { for a in var.apps : a.name => a }
  source   = "./modules/apps"

  project_name      = var.project_name
  name              = each.value.name
  app_type          = each.value.type
  platform_type     = each.value.infra.type
  platform_provider = each.value.infra.provider
  platform_config   = each.value.infra.config
  image_tag         = try(each.value.image_tag, "latest")
  github_config     = var.github_config
  ssh_keys          = try(var.ssh_keys, [])
  domains           = try(each.value.domains, [])
  env_vars          = try(each.value.env_vars, [])
  tags              = local.common_tags

  # Apps may depend on resources being created first.
  resource_outputs = module.resources
  depends_on = [module.resources]
}
