output "id" {
  description = "The ID of the App"
  value       = digitalocean_app.this.id
}

output "urn" {
  description = "The URN of the App"
  value       = digitalocean_app.this.urn
}

output "ingress" {
  description = "Default URL to access the App"
  value       = digitalocean_app.this.default_ingress
}

output "live_url" {
  description = "Live URL of the App"
  value       = digitalocean_app.this.live_url
}

output "live_domain" {
  description = "Live domain of the App"
  value       = digitalocean_app.this.live_domain
}
