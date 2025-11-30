#
# DigitalOcean Project
# Assigns all DigitalOcean resources to a project for organisation.
#
# Uses direct URN references to maintain a static list structure, allowing
# Terraform to track changes even when URN values are unknown at plan time.
#
# Manually-added domains are preserved by querying the existing project
# and merging domain URNs with Terraform-managed resources.
#

locals {
  manual_domain_urns = data.external.project_domains.result.domain_urns != "" ? split(",", data.external.project_domains.result.domain_urns) : []

  terraform_managed_urns = concat(
    [for r in module.resources : r.urn if r.platform_provider == "digitalocean"],
    [for a in module.apps : a.urn if a.platform_provider == "digitalocean"]
  )

  all_project_resources = concat(
    local.terraform_managed_urns,
    local.manual_domain_urns
  )

  is_only_project = tonumber(data.external.project_count.result.count) == 1

  # Only create project when there are Terraform-managed DigitalOcean resources.
  should_create_project = length(local.terraform_managed_urns) > 0
}

resource "time_sleep" "wait_for_propagation" {
  count           = local.should_create_project ? 1 : 0
  create_duration = "30s"
  depends_on      = [module.resources, module.apps]
}

resource "digitalocean_project" "this" {
  count       = local.should_create_project ? 1 : 0
  name        = var.project_title
  description = var.project_description
  purpose     = "Web Application"
  environment = title(var.environment)
  resources   = local.all_project_resources
  is_default  = local.is_only_project

  depends_on = [time_sleep.wait_for_propagation]
}
