output "id" {
  description = "Status page ID"
  value       = peekaping_status_page.this.id
}

output "slug" {
  description = "Status page URL slug"
  value       = peekaping_status_page.this.slug
}

output "url" {
  description = "Full status page URL"
  value       = "https://${peekaping_status_page.this.slug}"
}
