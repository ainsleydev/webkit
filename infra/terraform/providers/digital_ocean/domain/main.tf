resource "digitalocean_record" "this" {
  domain = var.domain_name
  type   = "A"
  name   = "cms"
  value  = var.domain_ip
  ttl    = 1800
}
