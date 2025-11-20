# Peekaping Status Page
# Creates a public status page for monitoring display.

resource "peekaping_status_page" "this" {
  title       = var.title
  description = var.description
  slug        = var.slug
  published   = var.published
  theme       = var.theme
  icon        = var.icon
}
