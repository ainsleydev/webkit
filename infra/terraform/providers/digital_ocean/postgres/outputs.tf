output "id" {
  description = "The ID of the Postgres DB"
  value       = digitalocean_database_cluster.this.id
}

output "urn" {
  description = "The URN of the Postgres DB"
  value       = digitalocean_database_cluster.this.urn
}

output "host" {
  description = "The host of the Postgres cluster"
  value       = digitalocean_database_cluster.this.host
}

output "port" {
  description = "The port Postgres is running on"
  value       = digitalocean_database_cluster.this.port
  sensitive   = true
}

output "database" {
  description = "The database name"
  value       = digitalocean_database_db.this.name
}

output "user" {
  description = "The database user name"
  value       = digitalocean_database_user.this.name
}

output "connection_url" {
  description = "Alias for pool_uri - the connection pool URI"
  value       = digitalocean_database_connection_pool.this.uri
  sensitive   = true
}
