terraform {
  required_version = ">= 2.0.0"

  backend "s3" {
    # Backend configuration will be provided via backend.hcl
    # at runtime: terraform init -backend-config=backend.hcl
  }

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
    b2 = {
      source  = "Backblaze/b2"
      version = "~> 1.0"
    }
  }
}

# Provider configurations will be set via environment variables:
# - DIGITALOCEAN_TOKEN
# - B2_APPLICATION_KEY_ID
# - B2_APPLICATION_KEY
provider "digitalocean" {}
provider "b2" {}

# Instantiate each resource from the manifest
module "resources" {
  for_each = { for r in var.resources : r.name => r }
  source   = "./modules/resources"

  project_name = var.project_name
  name         = each.value.name
  type         = each.value.type
  provider     = each.value.provider
  config       = each.value.config
  tags         = var.tags
}

# TODO: "apps"
