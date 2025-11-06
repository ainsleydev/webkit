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
    slack = {
      source  = "pablovarela/slack"
      version = "~> 1.0"
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

provider "slack" {
  token = var.slack_bot_token
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
# SSH Keys
# Lookup personal SSH keys at plan time to avoid deferred reads
# Data sources are evaluated here (outside modules) to prevent deferred reads
#

# Determine which providers are actually being used for VMs
locals {
  uses_digitalocean_vms = anytrue([for a in var.apps : a.platform_provider == "digitalocean" && a.platform_type == "vm"])
  # Future providers can be added here:
  # uses_hetzner_vms = anytrue([for a in var.apps : a.platform_provider == "hetzner" && a.platform_type == "vm"])
}

# DigitalOcean SSH Keys (only lookup if DO VMs are in use)
data "digitalocean_ssh_key" "personal_keys" {
  for_each = local.uses_digitalocean_vms ? toset(var.ssh_keys) : toset([])
  name     = each.value
}

locals {
  # Provider-specific SSH key ID lists
  do_ssh_key_ids = [for k in data.digitalocean_ssh_key.personal_keys : k.id]

  # Future providers:
  # hetzner_ssh_key_ids = [for k in data.hcloud_ssh_key.personal_keys : k.id]
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
  do_ssh_key_ids = local.do_ssh_key_ids
  domains = try(each.value.domains, [])
  env_vars = try(each.value.env_vars, [])
  tags = local.common_tags

  # Apps may depend on resources being created first.
  resource_outputs = module.resources
  depends_on = [module.resources]
}

#
# Slack Channel
#
# Create a Slack channel for CI/CD alerts and notifications.
# The channel is archived (not deleted) on destroy to preserve message history.
#
resource "slack_conversation" "project_channel" {
  name              = "${var.project_name}-alerts"
  topic             = "CI/CD alerts and notifications for ${var.project_title}"
  is_private        = false
  action_on_destroy = "archive"

  # Permanent members (you + bot).
  permanent_members = [
    "U035SMG9XFG" # Ainsley Clark
  ]
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

  # Slack channel secret
  github_secrets_slack = {
    "TF_SLACK_CHANNEL_ID" = tomap({
      source_type = "slack"
      output_name = "channel_id"
    })
  }

  # Merge all GitHub secrets
  github_secrets = merge(
    local.github_secrets_resources,
    local.github_secrets_apps,
    local.github_secrets_slack
  )
}

resource "github_actions_secret" "resource_outputs" {
  for_each = local.github_secrets

  repository      = var.github_config.repo
  secret_name     = each.key
  plaintext_value = try(
    each.value["source_type"] == "resource" ? tostring(module.resources[each.value["resource_name"]][each.value["output_name"]]) :
    each.value["source_type"] == "app" ? tostring(module.apps[each.value["app_name"]][each.value["output_name"]]) :
    each.value["source_type"] == "slack" ? tostring(slack_conversation.project_channel.id) :
    "NOT_SET"
  )
  depends_on = [module.resources, module.apps, slack_conversation.project_channel]
}

#
# DigitalOcean Project
#
# Create the project and assign all DigitalOcean resources to it.
# Only includes resources where platform_provider is "digitalocean".
#
# Important: We use direct URN references (not try/compact) to maintain a static
# list structure. This allows Terraform to properly track changes even when URN
# values are unknown at plan time, preventing the two-apply cycle.
resource "digitalocean_project" "this" {
  name        = var.project_title
  description = var.project_description
  purpose     = "Web Application"
  environment = title(var.environment)

  # Direct references maintain static list structure even with unknown URN values
  resources = concat(
    [for r in module.resources : r.urn if r.platform_provider == "digitalocean"],
    [for a in module.apps : a.urn if a.platform_provider == "digitalocean"]
  )

  depends_on = [module.resources, module.apps]
}
