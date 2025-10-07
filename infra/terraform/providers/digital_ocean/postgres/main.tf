locals {
  db_prefix = lower(replace(var.name, "-", "_"))
  has_firewall_rules = length(var.allowed_droplet_ips) > 0 || length(var.allowed_ips_addr) > 0
}

# Database Cluster
# Ref: https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/database_cluster
resource "digitalocean_database_cluster" "this" {
  name       = var.name
  engine     = "pg"
  version    = var.pg_version
  size       = var.size
  region     = var.region
  node_count = var.node_count
  tags       = var.tags
}

# Database User
# Ref: https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/database_user
resource "digitalocean_database_user" "this" {
  cluster_id = digitalocean_database_cluster.this.id
  name       = "${local.db_prefix}_admin"
}

# Database
# Ref: https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/database_db
resource "digitalocean_database_db" "this" {
  cluster_id = digitalocean_database_cluster.this.id
  name       = local.db_prefix
}

# Connection Pool
# Ref: https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/data-sources/database_connection_pool
resource "digitalocean_database_connection_pool" "this" {
  cluster_id = digitalocean_database_cluster.this.id
  name       = "${local.db_prefix}_pool"
  mode       = "transaction"
  size       = 20
  db_name    = digitalocean_database_db.this.name
  user       = digitalocean_database_user.this.name
}

# Firewall
# Ref: https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/database_firewall
resource "digitalocean_database_firewall" "this" {
  count      = local.has_firewall_rules ? 1 : 0
  cluster_id = digitalocean_database_cluster.this.id

  dynamic "rule" {
    for_each = var.allowed_droplet_ips
    content {
      type  = "droplet"
      value = rule.value
    }
  }

  dynamic "rule" {
    for_each = var.allowed_ips_addr
    content {
      type  = "ip_addr"
      value = rule.value
    }
  }
}

# Grant privileges to the user, Payload seed scrips will fail if
# you try to drop the database (as a seed script).
resource "null_resource" "grant_permissions" {
  depends_on = [
    digitalocean_database_user.this,
    digitalocean_database_db.this
  ]

  # This assumes you have psql installed locally
  provisioner "local-exec" {
    command = <<-EOT
PGPASSWORD='${digitalocean_database_cluster.this.password}' psql \
-h ${digitalocean_database_cluster.this.host} \
-p ${digitalocean_database_cluster.this.port} \
-U ${digitalocean_database_cluster.this.user} \
-d ${digitalocean_database_db.this.name} \
-c "GRANT ALL PRIVILEGES ON SCHEMA public TO ${digitalocean_database_user.this.name};" \
-c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ${digitalocean_database_user.this.name};" \
-c "ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL PRIVILEGES ON TABLES TO ${digitalocean_database_user.this.name};"
EOT
  }
}
