resource "digitalocean_record" "this" {
  domain = var.domain
  type   = var.type
  name   = var.name
  value  = var.value
  ttl    = 1800
}
