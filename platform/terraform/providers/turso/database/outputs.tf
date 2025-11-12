output "id" {
  description = "The unique ID of the Turso database"
  value       = turso_database.this.db_id
}

output "hostname" {
  description = "The hostname/URL of the Turso database"
  value       = turso_database.this.hostname
}

output "database" {
  description = "The database name"
  value       = turso_database.this.name
}

output "auth_token" {
  description = "The authentication token for the database"
  value       = turso_database_token.this.token
  sensitive   = true
}

output "connection_url" {
  description = "The libsql connection URL (libsql://hostname)"
  value       = "libsql://${turso_database.this.hostname}"
  sensitive   = false
}

output "connection_url_with_token" {
  description = "The complete libsql connection URL with authentication token"
  value       = "libsql://${turso_database.this.hostname}?authToken=${turso_database_token.this.token}"
  sensitive   = true
}
