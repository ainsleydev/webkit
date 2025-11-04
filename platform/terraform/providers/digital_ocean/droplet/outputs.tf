output "id" {
  description = "The ID of the Droplet"
  value       = digitalocean_droplet.this.id
}

output "ip_address" {
  description = "The IPv4 address of the Droplet"
  value       = digitalocean_droplet.this.ipv4_address
}

output "urn" {
  description = "The URN of the Droplet"
  value       = digitalocean_droplet.this.urn
}

output "ssh_private_key" {
  description = "Private key for the server (Terraform generated)"
  value       = tls_private_key.this.private_key_pem
  sensitive   = true
}

output "ssh_public_key" {
  description = "Public key for the server (Terraform generated)"
  value       = tls_private_key.this.public_key_openssh
}

output "server_user" {
  description = "The SSH user for the server"
  value       = var.server_user
}
