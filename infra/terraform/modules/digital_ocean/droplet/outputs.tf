output "droplet_ip_address" {
  description = "The IPv4 address of the Droplet"
  value       = digitalocean_droplet.this.ipv4_address
}

output "droplet_id" {
  description = "The ID of the Droplet"
  value       = digitalocean_droplet.this.id
}

output "droplet_urn" {
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