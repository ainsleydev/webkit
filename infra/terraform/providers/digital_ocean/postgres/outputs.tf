output "id" {
  description = "The ID of the Postgres DB"
  value       = digitalocean_database_cluster.this.id
}

output "urn" {
  description = "The URN of the Postgres DB"
  value       = digitalocean_database_cluster.this.urn
}

output "port" {
  description = "The port Postgres is running on"
  value       = digitalocean_database_cluster.this.port
  sensitive   = true
}

output "pool_uri" {
  description = "The connection pool URI"
  value       = digitalocean_database_connection_pool.this.uri
  sensitive   = true
}
