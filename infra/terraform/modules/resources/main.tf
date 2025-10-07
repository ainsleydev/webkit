#
# Resource Module
# Generates module calls based on resources[] in app.json
#

# DigitalOcean Postgres
module "do_postgres" {
  count  = var.cloud_provider == "digitalocean" && var.type == "postgres" ? 1 : 0
  source = "../../providers/digital_ocean/postgres"

  name                = var.name
  pg_version          = try(var.config.engine_version, null)
  size                = try(var.config.size, null)
  region              = try(var.config.region, null)
  node_count          = try(var.config.node_count, null)
  allowed_droplet_ips = try(var.config.allowed_droplet_ips, [])
  allowed_ips_addr    = try(var.config.allowed_ips_addr, [])
  tags                = concat([var.project_name], var.tags)
}

# DigitalOcean S3 Bucket (Spaces)
module "do_bucket" {
  count  = var.cloud_provider == "digitalocean" && var.type == "s3" ? 1 : 0
  source = "../../providers/digital_ocean/bucket"

  name   = var.name
  region = try(var.config.region, null)
  acl    = try(var.config.acl, null)
}

# B2 S3 Bucket
module "b2_bucket" {
  count  = var.cloud_provider == "b2" && var.type == "s3" ? 1 : 0
  source = "../../providers/b2/bucket"

  bucket_name = var.name
  acl         = try(var.config.acl, null)
}
