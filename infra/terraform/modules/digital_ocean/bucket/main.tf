resource "digitalocean_spaces_bucket" "this" {
  name   = var.name
  region = var.region
  acl    = var.acl
}

resource "digitalocean_spaces_bucket_cors_configuration" "this" {
  bucket = digitalocean_spaces_bucket.this.id
  region = digitalocean_spaces_bucket.this.region

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET"]
    allowed_origins = ["*"]
    max_age_seconds = 31536000 // TODO, how long is this in hours?
  }
}

resource "digitalocean_cdn" "this" {
  origin = digitalocean_spaces_bucket.this.bucket_domain_name
}
