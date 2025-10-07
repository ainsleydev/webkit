output "app_id" {
  description = "The ID of the App"
  value       = digitalocean_app.this.id
}

output "app_urn" {
  description = "The URN of the App"
  value       = digitalocean_app.this.urn
}

output "app_ingress" {
  description = "Default URL to access the App"
  value       = digitalocean_app.this.default_ingress
}

output "app_live_url" {
  description = "Live URL of the App"
  value       = digitalocean_app.this.live_url
}

output "app_live_domain" {
  description = "Live domain of the App"
  value       = digitalocean_app.this.live_domain
}