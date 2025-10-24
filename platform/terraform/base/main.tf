terraform {
  required_version = ">= 1.13.0"

  # Backend configuration - details provided via backend.hcl
  # Remote backend configuration will be added dynamically via backend.tf

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
    b2 = {
      source  = "Backblaze/b2"
      version = "~> 0.10.0"
    }
    github = {
      source  = "integrations/github"
      version = "~> 5.0"
    }
  }
}

#
# Providers
#
provider "digitalocean" {
  token             = var.do_token
  spaces_access_id  = var.do_spaces_access_id
  spaces_secret_key = var.do_spaces_secret_key
}

provider "b2" {
  application_key    = var.b2_application_key
  application_key_id = var.b2_application_key_id
}

provider "github" {
  owner = var.github_config.owner
  token = var.github_token
}

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


  # Shortened environment names.
  environment_short_map = {
    PRODUCTION  = "prod"
    STAGING     = "stag"
    DEVELOPMENT = "dev"
    TEST        = "test"
  }
  environment_short = lookup(local.environment_short_map, upper(var.environment), lower(var.environment))
}

#
# Default B2 Bucket (always provisioned for every project)
#
module "default_b2_bucket" {
  source = "../providers/b2/bucket"

  bucket_name = "${var.project_name}-${var.environment}"
  acl         = "allPrivate"
}

#
# Resources (databases, storage, etc.)
#
module "resources" {
  for_each = {for r in var.resources : r.name => r}
  source   = "../modules/resources"

  project_name      = var.project_name
  name              = each.value.name
  platform_type     = each.value.platform_type
  platform_provider = each.value.platform_provider
  platform_config   = each.value.config
  tags              = local.common_tags
}

#
# Apps (services, applications)
#
module "apps" {
  for_each = {for a in var.apps : a.name => a}
  source   = "../modules/apps"

  project_name      = var.project_name
  name              = each.value.name
  app_type          = each.value.app_type
  platform_type     = each.value.platform_type
  platform_provider = each.value.platform_provider
  platform_config   = each.value.config
  image_tag = try(each.value.image_tag, "latest")
  github_config = {
    owner = var.github_config.owner
    repo  = var.github_config.repo
    token = var.github_token
  }
  ssh_keys = try(var.ssh_keys, [])
  domains = try(each.value.domains, [])
  env_vars = try(each.value.env_vars, [])
  tags = local.common_tags

  # Apps may depend on resources being created first.
  resource_outputs = module.resources
  depends_on = [module.resources]
}

#
# Secrets (GitHub)
#
locals {
  # Define which outputs each resource type has (known at plan time)
  resource_output_map = {
    postgres = ["id", "urn", "connection_url"]
    s3 = ["id", "urn", "bucket_name", "bucket_url", "region", "endpoint"]
  }

  # Build ALL secret keys from var.resources (fully known at plan time)
  github_secrets = merge([
    for resource in var.resources : {
      for output_name in lookup(local.resource_output_map, resource.platform_type, []) :
      upper("TF_${local.environment_short}_${replace(resource.name, "-", "_")}_${output_name}") => tomap({
        resource_name = resource.name
        output_name   = output_name
      })
    }
  ]...)
}

resource "github_actions_secret" "resource_outputs" {
  for_each = local.github_secrets

  repository      = var.github_config.repo
  secret_name     = each.key
  plaintext_value = try(
    tostring(module.resources[each.value["resource_name"]][each.value["output_name"]]),
    "NOT_SET"
  )
  depends_on = [module.resources, module.apps]
}
