terraform {
  required_version = ">= 1.13.0"

  //backend "s3" {
  # Backend configuration will be provided via backend.hcl
  # at runtime: terraform init -backend-config=backend.hcl
  //}

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

# Provider configurations will be set via environment variables:
# - DIGITALOCEAN_TOKEN
# - B2_APPLICATION_KEY_ID
# - B2_APPLICATION_KEY
provider "digitalocean" {}
provider "b2" {}

#
# Resources (databases, storage, etc.)
#
module "resources" {
  for_each = { for r in var.resources : r.name => r }
  source   = "./modules/resources"

  project_name   = var.project_name
  name           = each.value.name
  type           = each.value.type
  cloud_provider = each.value.provider
  config         = each.value.config
  tags           = concat([var.project_name], var.tags)
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
  infra_type        = each.value.infra.type
  cloud_provider    = each.value.infra.provider
  infra_config      = each.value.infra.config
  image_tag         = try(each.value.image_tag, "latest")
  github_config     = var.github_config
  user_ssh_key_name = var.user_ssh_key_name
  env_vars          = try(each.value.env_vars, [])
  tags              = try(var.tags, [])

  # Apps may depend on resources being created first.
  depends_on = [module.resources]
}
