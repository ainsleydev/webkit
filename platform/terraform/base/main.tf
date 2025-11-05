terraform {
  # NOTE: Keep this version in sync with internal/infra/tf.go (TerraformVersion constant)
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

  bucket_name                        = var.project_name
  acl                                = "allPrivate"
  days_from_hiding_to_deleting       = 1
  days_from_uploading_to_hiding      = 0
  lifecycle_rule_file_name_prefix    = ""
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
    token = var.github_token_classic
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

  # Define which outputs each app type has (known at plan time)
  app_output_map = {
    vm        = ["ip_address", "ssh_private_key", "server_user"]
    container = []
  }

  # Build secret keys from var.resources (fully known at plan time)
  github_secrets_resources = merge([
    for resource in var.resources : {
      for output_name in lookup(local.resource_output_map, resource.platform_type, []) :
      upper("TF_${local.environment_short}_${replace(resource.name, "-", "_")}_${output_name}") => tomap({
        source_type   = "resource"
        resource_name = resource.name
        output_name   = output_name
      })
    }
  ]...)

  # Build secret keys from var.apps (fully known at plan time)
  github_secrets_apps = merge([
    for app in var.apps : {
      for output_name in lookup(local.app_output_map, app.platform_type, []) :
      upper("TF_${local.environment_short}_${replace(app.name, "-", "_")}_${output_name}") => tomap({
        source_type = "app"
        app_name    = app.name
        output_name = output_name
      })
    }
  ]...)

  # Merge all GitHub secrets
  github_secrets = merge(local.github_secrets_resources, local.github_secrets_apps)
}

resource "github_actions_secret" "resource_outputs" {
  for_each = local.github_secrets

  repository      = var.github_config.repo
  secret_name     = each.key
  plaintext_value = try(
    each.value["source_type"] == "resource" ? tostring(module.resources[each.value["resource_name"]][each.value["output_name"]]) : tostring(module.apps[each.value["app_name"]][each.value["output_name"]]),
    "NOT_SET"
  )
  depends_on = [module.resources, module.apps]
}
