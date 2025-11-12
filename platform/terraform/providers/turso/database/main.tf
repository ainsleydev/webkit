# Turso Database
# Ref: https://registry.terraform.io/providers/jpedroh/turso/latest/docs/resources/database
resource "turso_database" "this" {
  organization_name = var.organisation
  name              = var.name
  group             = var.group
  size_limit        = var.size_limit
}

# Turso Database Token
# Ref: https://registry.terraform.io/providers/jpedroh/turso/latest/docs/resources/database_token
resource "turso_database_token" "this" {
  organization_name = var.organisation
  database_name     = turso_database.this.name
  expiration        = "never"
  authorization     = "full-access"
}
