#
# WebKit Base Infrastructure
# Core Terraform configuration for provisioning project infrastructure.
#

terraform {
  # Keep in sync with internal/infra/tf.go (TerraformVersion constant).
  required_version = ">= 1.13.0"

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = "~> 1.0"
    }
    b2 = {
      source  = "Backblaze/b2"
      version = "~> 0.10.0"
    }
    turso = {
      source  = "jpedroh/turso"
      version = "~> 0.3.0"
    }
    github = {
      source  = "integrations/github"
      version = "~> 5.0"
    }
    slack = {
      source  = "pablovarela/slack"
      version = "~> 1.0"
    }
    time = {
      source  = "hashicorp/time"
      version = "~> 0.9"
    }
    peekaping = {
      source  = "tafaust/peekaping"
      version = "~> 0.2.1"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.5"
    }
    json-formatter = {
      source  = "TheNicholi/json-formatter"
      version = "~> 0.1.0"
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

provider "hcloud" {
  token = var.hetzner_token
}

provider "b2" {
  application_key    = var.b2_application_key
  application_key_id = var.b2_application_key_id
}

provider "turso" {
  api_token = var.turso_api_token
}

provider "github" {
  owner = var.github_config.owner
  token = var.github_token
}

provider "slack" {
  token = var.slack_bot_token
}

provider "peekaping" {
  endpoint = var.peekaping_endpoint
  api_key  = var.peekaping_api_key
}

#
# Locals
#
locals {
  # Default tags applied to all resources (normalised to lowercase).
  default_tags = [
    lower(var.project_name),
    lower(var.environment),
    "terraform",
  ]

  # Combined tags: default + custom (all normalised to lowercase).
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
# Default B2 Bucket
#
module "default_b2_bucket" {
  source = "../providers/b2/bucket"

  bucket_name                     = var.project_name
  acl                             = "allPrivate"
  days_from_hiding_to_deleting    = 1
  days_from_uploading_to_hiding   = 0
  lifecycle_rule_file_name_prefix = ""
}

#
# Slack Channel
# Created before resources/apps so the channel ID is available for alerts.
#
resource "slack_conversation" "project_channel" {
  name              = "alerts-${var.project_name}"
  topic             = "CI/CD alerts and notifications for ${replace(var.project_title, "/[^a-zA-Z0-9 ]/", " ")}"
  is_private        = false
  action_on_destroy = "archive"

  permanent_members = [
    "U035SMG9XFG" # Ainsley Clark
  ]
}

#
# Slack Channel GitHub Secret
# Created early to ensure availability even if later provisioning fails.
#
resource "github_actions_secret" "slack_channel_id" {
  repository      = var.github_config.repo
  secret_name     = "TF_SLACK_CHANNEL_ID"
  plaintext_value = slack_conversation.project_channel.id

  depends_on = [slack_conversation.project_channel]
}

#
# Resources (databases, storage, etc.)
#
module "resources" {
  for_each = { for r in var.resources : r.name => r }
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
  for_each = { for a in var.apps : a.name => a }
  source   = "../modules/apps"

  project_name      = var.project_name
  name              = each.value.name
  app_type          = each.value.app_type
  platform_type     = each.value.platform_type
  platform_provider = each.value.platform_provider
  platform_config   = each.value.config
  image_tag         = try(each.value.image_tag, "latest")
  github_config = {
    owner = var.github_config.owner
    repo  = var.github_config.repo
    token = var.github_token_classic
  }
  do_ssh_key_ids      = local.do_ssh_key_ids
  hetzner_ssh_key_ids = local.hetzner_ssh_key_ids
  domains             = try(each.value.domains, [])
  env_vars            = try(each.value.env_vars, [])
  tags                = local.common_tags
  slack_webhook_url   = var.slack_webhook_url
  slack_channel_name  = slack_conversation.project_channel.name

  resource_outputs = module.resources
  depends_on       = [module.resources]
}

#
# Monitoring
# Only created if monitors are configured.
#
module "monitoring" {
  count  = length(var.monitors) > 0 ? 1 : 0
  source = "../modules/monitoring"

  project_name        = var.project_name
  project_title       = var.project_title
  environment         = var.environment
  monitors            = var.monitors
  slack_webhook_url   = var.slack_webhook_url
  brand_primary_color = var.brand_primary_color
  brand_logo_url      = var.brand_logo_url
  brand_icon_url      = var.brand_icon_url
  status_page_domain  = var.status_page_domain
  status_page_slug    = var.status_page_slug
  status_page_theme   = var.status_page_theme
  peekaping_endpoint  = var.peekaping_endpoint

  depends_on = [module.apps, module.resources]
}

#
# Local Outputs File
# Writes monitoring and Slack data to .webkit/outputs.json for README badges.
# Only created when monitoring is enabled and project_root is provided.
#
data "json-formatter_format_json" "webkit_outputs" {
  count = length(var.monitors) > 0 && var.project_root != "" ? 1 : 0

  indent = "\t"
  json = jsonencode({
    peekaping_endpoint = var.peekaping_endpoint
    monitors = [for m in module.monitoring[0].monitors : {
      id   = m.id
      name = m.name
      type = m.type
    }]
    slack = {
      channel_name = slack_conversation.project_channel.name
      channel_id   = slack_conversation.project_channel.id
    }
  })
}

resource "local_file" "webkit_outputs" {
  count    = length(var.monitors) > 0 && var.project_root != "" ? 1 : 0
  filename = "${var.project_root}/.webkit/outputs.json"
  content  = data.json-formatter_format_json.webkit_outputs[0].result

  depends_on = [module.monitoring]
}
