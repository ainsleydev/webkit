output "postgres_id" {
  description = "The ID of the Postgres DB"
  value       = digitalocean_database_cluster.this.id
}

output "postgres_urn" {
  description = "The URN of the Postgres DB"
  value = digitalocean_database_cluster.this.urn
}

output "postgres_port" {
  description = "The port Postgres is running on"
  value     = digitalocean_database_cluster.this.port
  sensitive = true
}

output "postgres_pool_uri" {
  description = "The connection pool URI"
  value     = digitalocean_database_connection_pool.this.uri
  sensitive = true
}
