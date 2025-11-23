#
# Data Sources
# Lookups and external data required at plan time.
#

#
# SSH Keys
# Evaluated outside modules to prevent deferred reads.
#
locals {
  uses_digitalocean_vms = anytrue([for a in var.apps : a.platform_provider == "digitalocean" && a.platform_type == "vm"])
  uses_hetzner_vms      = anytrue([for a in var.apps : a.platform_provider == "hetzner" && a.platform_type == "vm"])
}

data "digitalocean_ssh_key" "personal_keys" {
  for_each = local.uses_digitalocean_vms ? toset(var.digitalocean_ssh_keys) : toset([])
  name     = each.value
}

data "hcloud_ssh_key" "personal_keys" {
  for_each = local.uses_hetzner_vms ? toset(var.hetzner_ssh_keys) : toset([])
  name     = each.value
}

locals {
  do_ssh_key_ids      = [for k in data.digitalocean_ssh_key.personal_keys : k.id]
  hetzner_ssh_key_ids = [for k in data.hcloud_ssh_key.personal_keys : k.id]
}

#
# DigitalOcean Project Data
# Queries for project management and resource assignment.
#
data "external" "project_domains" {
  program = ["bash", "${path.module}/scripts/get_project_domains.sh"]

  query = {
    project_id    = try(var.digitalocean_project_id, "")
    project_title = var.project_title
    do_token      = var.do_token
  }
}

data "external" "project_count" {
  program = ["bash", "${path.module}/scripts/count_projects.sh"]

  query = {
    do_token = var.do_token
  }
}
