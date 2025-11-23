#
# Turso Database
# Provisions a SQLite database on Turso with authentication token.
#
# Ref: https://registry.terraform.io/providers/jpedroh/turso/latest/docs/resources/database
#

#
# Database
#
resource "turso_database" "this" {
  organization_name = var.organisation
  name              = var.name
  group             = var.group
  size_limit        = var.size_limit

  # Workaround for Turso provider limitation: the 'group' attribute is not
  # captured during import, causing Terraform to want to replace the database.
  # We ignore changes to 'group' since databases cannot be moved between groups anyway.
  lifecycle {
    ignore_changes = [group]
  }
}

#
# Database Token
#
# Ref: https://registry.terraform.io/providers/jpedroh/turso/latest/docs/resources/database_token
#
resource "turso_database_token" "this" {
  organization_name = var.organisation
  database_name     = turso_database.this.name
  expiration        = "never"
  authorization     = "full-access"
}
