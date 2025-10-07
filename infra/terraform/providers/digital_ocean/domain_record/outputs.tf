output "id" {
  description = "The ID of the DNS record"
  value       = digitalocean_record.this.id
}

output "fqdn" {
  description = "The FQDN of the record"
  value       = digitalocean_record.this.fqdn
}
