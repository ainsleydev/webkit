# Spaces Bucket
# Ref: https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/spaces_bucket
resource "digitalocean_spaces_bucket" "this" {
  name   = var.name
  region = var.region
  acl    = var.acl
}

# CORS Configuration
# Ref: https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/spaces_bucket_cors_configuration
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

# CDN
# Ref: https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/cdn
resource "digitalocean_cdn" "this" {
  origin = digitalocean_spaces_bucket.this.bucket_domain_name
}
