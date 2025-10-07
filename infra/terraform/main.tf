terraform {
  backend "s3" {
    bucket                      = "ainsley-dev-terraform"
    key                         = "search-spares/terraform.tfstate"
    region                      = "eu-central-003"
    skip_credentials_validation = true
    skip_region_validation      = true
    skip_requesting_account_id  = true
    use_path_style              = true

    endpoints = {
      s3 = "https://s3.eu-central-003.backblazeb2.com"
    }
  }

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
    logtail = {
      source  = "BetterStackHQ/logtail"
      version = ">= 0.1.0"
    }
    github = {
      source  = "integrations/github"
      version = "~> 5.0"
    }
    b2 = {
      source = "Backblaze/b2"
    }
    slack = {
      source  = "pablovarela/slack"
      version = "~> 1.0"
    }
  }
}

provider "digitalocean" {
  token             = var.digital_ocean_config.api_key
  spaces_access_id  = var.digital_ocean_config.spaces_access_key
  spaces_secret_key = var.digital_ocean_config.spaces_secret_key
}

provider "github" {
  owner = "ainsleydev"
  token = var.github_config.token
}

provider "b2" {
  application_key    = var.back_blaze_config.application_key
  application_key_id = var.back_blaze_config.application_key_id
}

provider "slack" {
  token = var.slack_config.bot_token
}

# ------------------------------------------------------------------------------
# Bucket
# ------------------------------------------------------------------------------

module "bucket" {
  source = "./modules/digital_ocean/bucket"
  name   = "${var.project_name}-store"
}

# ------------------------------------------------------------------------------
# Database
# ------------------------------------------------------------------------------

module "postgres" {
  source = "./modules/digital_ocean/postgres"

  name                = "${var.project_name}-db"
  size                = "db-s-1vcpu-1gb"
  region              = "lon1"
  node_count          = 1
  allowed_droplet_ips = [module.cms.droplet_id]
  allowed_ips_addr    = ["185.16.161.205"] // School Lane
}

# ------------------------------------------------------------------------------
# CMS
# ------------------------------------------------------------------------------

module "cms" {
  source = "./modules/digital_ocean/droplet"

  name              = "${var.project_name}-cms"
  droplet_size      = "s-1vcpu-1gb"
  droplet_region    = "lon1"
  user_ssh_key_name = "Ainsley - Mac Studio" # TODO: This should be an array
}


# ------------------------------------------------------------------------------
# Domain(s)
# ------------------------------------------------------------------------------

resource "digitalocean_domain" "domain" {
  name = "searchspares.com"
}

resource "digitalocean_domain" "domain-co-uk" {
  name = "searchspares.co.uk"
}

resource "digitalocean_record" "cms_a" {
  domain = digitalocean_domain.domain.name
  type   = "A"
  name   = "cms"
  value  = module.cms.droplet_ip_address
  ttl    = 1800
}

# ------------------------------------------------------------------------------
# Web App
# ------------------------------------------------------------------------------

module "web" {
  source = "./modules/digital_ocean/app"

  name               = "${var.project_name}-web"
  github_config      = var.github_config
  service_name       = "sveltekit"
  instance_size_slug = "apps-s-1vcpu-1gb"
  instance_count     = 1
  http_port          = 3001
  image_tag          = "latest"
  health_check_path  = "/api/ping"

  domains = [
    { name = var.url, type = "PRIMARY", zone = var.url },
    { name = "www.${var.url}", type = "ALIAS", zone = var.url }
  ]

  depends_on = [
    module.bucket,
    module.postgres,
    module.cms,
  ]

  envs = [
    {
      key   = "PAYLOAD_URL"
      value = "https://cms.${var.url}"
    },
    {
      key   = "PAYLOAD_API_KEY"
      value = var.payload_config.api_key
      type  = "SECRET"
    },
    {
      key   = "PUBLIC_PAYLOAD_DASHBOARD_URL"
      value = "https://cms.${var.url}/admin"
    },
    {
      key   = "PUBLIC_GOOGLE_MAPS_API_KEY"
      value = var.google_places_api_key
    },
    {
      key   = "PUBLIC_GOOGLE_RECAPTCHA_SITE_KEY"
      value = var.google_recaptcha_site_key
    },
    {
      key   = "GOOGLE_RECAPTCHA_SECRET_KEY"
      value = var.google_recaptcha_secret_key
      type  = "SECRET"
    }
  ]
}

# ------------------------------------------------------------------------------
# Storage Backup
# ------------------------------------------------------------------------------

module "b2_backup" {
  source      = "./modules/b2/bucket"
  bucket_name = "${var.project_name}-backup"
  acl         = "allPrivate"
}

# ------------------------------------------------------------------------------
# Project
# Import using terraform import -var-file=.tfvars digitalocean_project.project 7b4015f0-b722-4802-b334-1dde7692ac51
# ------------------------------------------------------------------------------

resource "digitalocean_project" "project" {
  name        = var.project_nice_name
  purpose     = "Web Application"
  environment = "Production"
  resources = [
    digitalocean_domain.domain.urn,
    digitalocean_domain.domain-co-uk.urn,
    module.bucket.urn,
    module.postgres.postgres_urn,
    module.cms.droplet_urn,
    module.web.app_urn,
  ]
}

# ------------------------------------------------------------------------------
# Github
# ------------------------------------------------------------------------------

locals {
  github_secrets = {
    # Digital Ocean
    "DO_ACCESS_TOKEN"      = var.digital_ocean_config.api_key
    "DO_SPACES_ACCESS_KEY" = var.digital_ocean_config.spaces_access_key
    "DO_SPACES_SECRET_KEY" = var.digital_ocean_config.spaces_secret_key
    "DO_SPACES_BUCKET"     = module.bucket.name
    "DO_SPACES_REGION"     = module.bucket.region

    # Server (CMS)
    "CMS_SERVER_IP"       = module.cms.droplet_ip_address
    "CMS_SERVER_USER"     = "root"
    "CMS_SSH_PRIVATE_KEY" = module.cms.ssh_private_key
    "CMS_PAYLOAD_SECRET"  = var.payload_config.secret

    # Database
    "DATABASE_ID"  = module.postgres.postgres_id
    "DATABASE_URL" = module.postgres.postgres_pool_uri

    # BackBlaze
    "B2_APPLICATION_KEY"    = var.back_blaze_config.application_key
    "B2_APPLICATION_KEY_ID" = var.back_blaze_config.application_key_id
    "B2_BUCKET_NAME"        = module.b2_backup.bucket_name

    # Slack
    "SLACK_BOT_TOKEN"  = var.slack_config.bot_token
    "SLACK_CHANNEL_ID" = slack_conversation.project_channel.id

    # Misc - TODO, these will become dynamic at some point.
    "API_KEY"               = var.api_key
    "GOOGLE_PLACES_API_KEY" = var.google_places_api_key
    "RESEND_TOKEN"          = var.resend_token
  }
}

resource "github_actions_secret" "secrets" {
  for_each = local.github_secrets

  repository      = var.github_config.repo
  secret_name     = each.key
  plaintext_value = each.value
}

# ------------------------------------------------------------------------------
# Slack
# ------------------------------------------------------------------------------

resource "slack_conversation" "project_channel" {
  name              = "${var.project_name}-alerts"
  topic             = "Alerts for ${var.project_name} project"
  is_private        = false
  action_on_destroy = "archive"

  # Add yourself as a permanent member using your Slack user ID.
  # You can find your user ID by running:
  # curl -H "Authorization: Bearer xoxb-your-bot-token" https://slack.com/api/users.list | jq '.members[] | {id, name, real_name}'
  # Look for your username or real name in the output, and copy the "id" value.
  permanent_members = ["U035SMG9XFG"]
}
