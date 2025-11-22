output "id" {
  description = "The ID of the status page."
  value       = peekaping_status_page.this.id
}

output "slug" {
  description = "The slug of the status page."
  value       = peekaping_status_page.this.slug
}

output "url" {
  description = "The public URL of the status page."
  value       = "https://${peekaping_status_page.this.slug}"
}
