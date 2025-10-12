#
# Postgres Outputs
# These should match the outputs[] array in app.json resources
#
output "connection_url" {
  description = "Database connection URL (pool)"
  value = (
    var.platform_type == "postgres" && var.platform_provider == "digitalocean" ? module.do_postgres[0].connection_url :
    null
  )
  sensitive = true
}

output "host" {
  description = "Database host"
  value = (
    var.platform_type == "postgres" && var.platform_provider == "digitalocean" ? module.do_postgres[0].host :
    null
  )
}

output "port" {
  description = "Database port"
  value = (
    var.platform_type == "postgres" && var.platform_provider == "digitalocean" ? module.do_postgres[0].port :
    null
  )
  sensitive = true
}

output "database" {
  description = "Database name"
  value = (
    var.platform_type == "postgres" && var.platform_provider == "digitalocean" ? module.do_postgres[0].database :
    null
  )
}

#
# Storage Outputs (S3-compatible)
# These should match the outputs[] array in app.json resources
#
output "bucket_name" {
  description = "S3 bucket name"
  value = (
    var.platform_provider == "digitalocean" && var.platform_type == "s3" ? module.do_bucket[0].name :
    var.platform_provider == "b2" && var.platform_type == "s3" ? module.b2_bucket[0].name :
    null
  )
}

output "endpoint" {
  description = "S3 endpoint URL"
  value = (
    var.platform_provider == "digitalocean" && var.platform_type == "s3" ? module.do_bucket[0].domain_name :
    null
  )
}

output "region" {
  description = "S3 region"
  value = (
    var.platform_provider == "digitalocean" && var.platform_type == "s3" ? module.do_bucket[0].region :
    null
  )
}
