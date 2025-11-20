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
      version = "~> 0.1.1"
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
  email    = var.peekaping_email
  password = var.peekaping_password
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
  uses_hetzner_vms      = anytrue([for a in var.apps : a.platform_provider == "hetzner" && a.platform_type == "vm"])
}

# DigitalOcean SSH Keys (only lookup if DO VMs are in use)
data "digitalocean_ssh_key" "personal_keys" {
  for_each = local.uses_digitalocean_vms ? toset(var.digitalocean_ssh_keys) : toset([])
  name     = each.value
}

# Hetzner SSH Keys (only lookup if Hetzner VMs are in use)
data "hcloud_ssh_key" "personal_keys" {
  for_each = local.uses_hetzner_vms ? toset(var.hetzner_ssh_keys) : toset([])
  name     = each.value
}

locals {
  # Provider-specific SSH key ID lists
  do_ssh_key_ids      = [for k in data.digitalocean_ssh_key.personal_keys : k.id]
  hetzner_ssh_key_ids = [for k in data.hcloud_ssh_key.personal_keys : k.id]
}

#
# Default B2 Bucket (always provisioned for every project)
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
#
# Create a Slack channel for CI/CD alerts and notifications.
# The channel is archived (not deleted) on destroy to preserve message history.
# This is created before resources/apps so the channel ID is available for alert configuration.
#
resource "slack_conversation" "project_channel" {
  name              = "alerts-${var.project_name}"
  topic             = "CI/CD alerts and notifications for ${replace(var.project_title, "/[^a-zA-Z0-9 ]/", " ")}"
  is_private        = false
  action_on_destroy = "archive"

  # Permanent members (you + bot).
  permanent_members = [
    "U035SMG9XFG" # Ainsley Clark
  ]
}

#
# Slack Channel GitHub Secret
#
# Create the Slack channel ID secret immediately after channel creation.
# This ensures the channel ID is available in GitHub Actions even if resource/app
# provisioning fails later, allowing CI/CD pipelines to send notifications.
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
  do_ssh_key_ids     = local.do_ssh_key_ids
  hetzner_ssh_key_ids = local.hetzner_ssh_key_ids
  domains            = try(each.value.domains, [])
  env_vars           = try(each.value.env_vars, [])
  tags               = local.common_tags
  slack_webhook_url  = var.slack_webhook_url
  slack_channel_name = slack_conversation.project_channel.name

  # Apps may depend on resources being created first.
  resource_outputs = module.resources
  depends_on       = [module.resources]
}

#
# Monitoring
#
# Only create the monitoring module if there are monitors configured.
# This prevents provider initialization when monitoring is not in use.
#
module "monitoring" {
  count  = length(var.monitors) > 0 ? 1 : 0
  source = "../modules/monitoring"

  providers = {
    peekaping = peekaping
  }

  project_name         = var.project_name
  project_title        = var.project_title
  environment          = var.environment
  monitors             = var.monitors
  slack_webhook_url    = var.slack_webhook_url
  brand_primary_color  = var.brand_primary_color
  brand_logo_url       = var.brand_logo_url

  # Monitoring depends on apps and resources being created.
  depends_on = [module.apps, module.resources]
}

#
# Resource and App GitHub Secrets
#
# These secrets are created after resources/apps are provisioned.
# The Slack channel secret is created separately (earlier) to ensure it's available
# even if resource/app provisioning fails.
#
locals {
  # Define which outputs each resource type has (known at plan time)
  resource_output_map = {
    postgres = ["id", "urn", "connection_url"]
    s3       = ["id", "urn", "bucket_name", "bucket_url", "region", "endpoint"]
    sqlite   = ["id", "connection_url", "auth_token", "host", "database"]
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

  # Merge resource and app GitHub secrets (Slack channel secret created separately)
  github_secrets = merge(
    local.github_secrets_resources,
    local.github_secrets_apps
  )
}

resource "github_actions_secret" "resource_outputs" {
  for_each = local.github_secrets

  repository  = var.github_config.repo
  secret_name = each.key
  plaintext_value = try(
    each.value["source_type"] == "resource" ? tostring(module.resources[each.value["resource_name"]][each.value["output_name"]]) :
    each.value["source_type"] == "app" ? tostring(module.apps[each.value["app_name"]][each.value["output_name"]]) :
    "NOT_SET"
  )
  depends_on = [module.resources, module.apps]
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
#
# Additionally, we preserve manually-added domains by querying the existing project
# via the DigitalOcean API and merging domain URNs with Terraform-managed resources.

# Query existing project to get manually-added domain URNs
# This allows manual domain management while Terraform manages other resources
data "external" "project_domains" {
  program = ["bash", "${path.module}/scripts/get_project_domains.sh"]

  query = {
    project_id    = try(var.digitalocean_project_id, "")
    project_title = var.project_title
    do_token      = var.do_token
  }
}

# Count total projects in the account to determine if this should be default
data "external" "project_count" {
  program = ["bash", "${path.module}/scripts/count_projects.sh"]

  query = {
    do_token = var.do_token
  }
}

locals {
  # Parse comma-separated domain URNs from external script
  manual_domain_urns = data.external.project_domains.result.domain_urns != "" ? split(",", data.external.project_domains.result.domain_urns) : []

  # Terraform-managed resources (apps, databases, buckets, etc.)
  terraform_managed_urns = concat(
    [for r in module.resources : r.urn if r.platform_provider == "digitalocean"],
    [for a in module.apps : a.urn if a.platform_provider == "digitalocean"]
  )

  # Merge Terraform-managed resources with manually-added domains
  all_project_resources = concat(
    local.terraform_managed_urns,
    local.manual_domain_urns
  )

  # Set as default if this is the only project in the account
  is_only_project = tonumber(data.external.project_count.result.count) == 1
}

# Wait for DigitalOcean API propagation after app/resource creation
resource "time_sleep" "wait_for_propagation" {
  create_duration = "30s"
  depends_on      = [module.resources, module.apps]
}

resource "digitalocean_project" "this" {
  name        = var.project_title
  description = var.project_description
  purpose     = "Web Application"
  environment = title(var.environment)
  resources   = local.all_project_resources
  is_default  = local.is_only_project

  depends_on = [time_sleep.wait_for_propagation]
}
