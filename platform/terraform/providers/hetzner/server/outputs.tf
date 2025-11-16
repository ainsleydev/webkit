output "id" {
  description = "Server ID"
  value       = hcloud_server.this.id
}

output "ip_address" {
  description = "Server IPv4 address"
  value       = hcloud_server.this.ipv4_address
}

output "ssh_private_key" {
  description = "Generated SSH private key"
  value       = tls_private_key.this.private_key_pem
  sensitive   = true
}

output "ssh_public_key" {
  description = "Generated SSH public key"
  value       = tls_private_key.this.public_key_openssh
}

output "server_user" {
  description = "SSH user for server access"
  value       = var.server_user
}
