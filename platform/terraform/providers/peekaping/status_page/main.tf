#
# Peekaping Status Page
# Creates a public status page for displaying monitor status.
#
# Ref: https://registry.terraform.io/providers/tafaust/peekaping/latest/docs/resources/status_page
#
resource "peekaping_status_page" "this" {
  title       = var.title
  description = var.description
  slug        = var.slug
  published   = var.published
  theme       = var.theme
  icon        = var.icon_url
  domains     = var.domains
  monitor_ids = var.monitor_ids

  # Workaround: Provider bug - doesn't return domains after apply.
  # See: https://github.com/tafaust/terraform-provider-peekaping/issues/16
  lifecycle {
    ignore_changes = [domains]
  }
}
