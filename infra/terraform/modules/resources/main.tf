#
# Resource Module
# Generates module calls based on resources[] in app.json
#

# DigitalOcean Postgres
module "do_postgres" {
  count  = var.provider == "digitalocean" && var.type == "postgres" ? 1 : 0
  source = "../../providers/digital_ocean/postgres"

  name                = var.name
  pg_version          = lookup(var.config, "engine_version", "17")
  size                = lookup(var.config, "size", "db-s-1vcpu-1gb")
  region              = lookup(var.config, "region", "lon1")
  node_count          = lookup(var.config, "node_count", 1)
  allowed_droplet_ips = lookup(var.config, "allowed_droplet_ips", [])
  allowed_ips_addr    = lookup(var.config, "allowed_ips_addr", [])
  tags                = concat([var.project_name], var.tags)
}

# DigitalOcean S3 Bucket (Spaces)
module "do_bucket" {
  count  = var.provider == "digitalocean" && var.type == "s3" ? 1 : 0
  source = "../../providers/digital_ocean/bucket"

  name   = var.name
  region = lookup(var.config, "region", "ams3")
  acl    = lookup(var.config, "acl", "public-read")
}

# B2 S3 Bucket
module "b2_bucket" {
  count  = var.provider == "b2" && var.type == "s3" ? 1 : 0
  source = "../../providers/b2/bucket"

  bucket_name = var.name
  acl         = lookup(var.config, "acl", "allPrivate")
}
