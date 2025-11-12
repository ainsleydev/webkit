#
# Resource Module
# Generates module calls based on resources[] in app.json
#

# DigitalOcean Postgres
module "do_postgres" {
  count  = var.platform_provider == "digitalocean" && var.platform_type == "postgres" ? 1 : 0
  source = "../../providers/digital_ocean/postgres"

  name       = "${var.project_name}-${var.name}"
  pg_version = try(var.platform_config.pg_version, "17")
  size       = try(var.platform_config.size, "db-s-1vcpu-1gb")
  region     = try(var.platform_config.region, "lon1")
  node_count = try(var.platform_config.node_count, 1)
  tags       = try(var.tags, [])

  allowed_ips_addr = try(
    jsondecode(var.platform_config.allowed_ips_addr),
    []
  )
  allowed_droplet_ips = try(
    jsondecode(var.platform_config.allowed_droplet_ips),
    []
  )
}

# DigitalOcean S3 Bucket (Spaces)
module "do_bucket" {
  count  = var.platform_provider == "digitalocean" && var.platform_type == "s3" ? 1 : 0
  source = "../../providers/digital_ocean/bucket"

  name   = "${var.project_name}-${var.name}"
  region = try(var.platform_config.region, "ams3")
  acl    = try(var.platform_config.acl, "public-read")
}

# B2 S3 Bucket
module "b2_bucket" {
  count  = var.platform_provider == "b2" && var.platform_type == "s3" ? 1 : 0
  source = "../../providers/b2/bucket"

  bucket_name                     = var.name
  acl                             = try(var.platform_config.acl, null)
  days_from_hiding_to_deleting    = try(var.platform_config.days_from_hiding_to_deleting, null)
  days_from_uploading_to_hiding   = try(var.platform_config.days_from_uploading_to_hiding, null)
  lifecycle_rule_file_name_prefix = try(var.platform_config.lifecycle_rule_file_name_prefix, null)
}

# Turso Database
module "turso_database" {
  count  = var.platform_provider == "turso" && var.platform_type == "turso" ? 1 : 0
  source = "../../providers/turso/database"

  name         = var.name
  organization = var.platform_config.organization
  group        = try(var.platform_config.group, "default")
  size_limit   = try(var.platform_config.size_limit, null)
  tags         = try(var.tags, [])
}
