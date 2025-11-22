#
# Monitoring Module
#
# Orchestrates Peekaping resources for monitoring apps and infrastructure.
#

#
# Locals
#
locals {
  # Tag colour: use brand primary colour if set, otherwise fallback to blue.
  tag_colour = var.brand_primary_color != null ? var.brand_primary_color : "#3B82F6"

  # Currently - Slack - #alerts
  notification_ids = ["7e4f8d2e-5720-4b07-9ce4-3f639b5e4647"]

  # Static tag IDs (environment and webkit tags).
  static_tag_ids = [
    "ac5b2626-3425-4496-a318-ede51ce7baa8",
    "dd1151bc-1d7a-42d1-8166-f87b7b180798",
  ]
}

#
# Project Tag
#
module "project_tag" {
  source = "../../providers/peekaping/tag"

  name        = var.project_title
  colour      = local.tag_colour
  description = "Monitor for ${var.project_title}"
}

#
# Monitors
#
module "monitors" {
  source = "../../providers/peekaping/monitors"

  monitors           = var.monitors
  peekaping_endpoint = var.peekaping_endpoint
  notification_ids   = local.notification_ids
  tag_ids            = concat([module.project_tag.id], local.static_tag_ids)

  depends_on = [module.project_tag]
}

#
# Status Page
#
module "status_page" {
  source = "../../providers/peekaping/status_page"
  count  = length(var.monitors) > 0 ? 1 : 0

  title       = "${var.project_title} Status"
  description = "Public status page for ${var.project_title} services"
  slug        = lower(replace(var.project_name, "_", "-"))
  icon_url    = var.brand_icon_url
  domains     = var.status_page_domain != null ? [var.status_page_domain] : []
  monitor_ids = module.monitors.all_ids

  depends_on = [module.monitors]
}
