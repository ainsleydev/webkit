output "slack_channel_id" {
  value       = slack_conversation.project_channel.id
  description = "ID of the Slack channel"
}

output "slack_channel" {
  value       = slack_conversation.project_channel.name
  description = "Name of the Slack channel"
}

# CMS

output "cms_server_ip" {
  value       = module.cms.droplet_ip_address
  description = "The IP of the CMS droplet"
}

output "cms_server_private_key" {
  value       = module.cms.ssh_private_key
  description = "The private key of the CMS droplet"
  sensitive   = true
}

# Database

output "database_uri" {
  value       = module.postgres.postgres_pool_uri
  description = "URI of the Postgres pool"
  sensitive   = true
}

# Bucket

output "bucket_region" {
  value = module.bucket.region
}

output "bucket_name" {
  value = module.bucket.name
}

output "bucket_endpoint" {
  value = module.bucket.domain_name
}