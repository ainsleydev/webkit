resource "b2_bucket" "this" {
  bucket_name = var.bucket_name
  bucket_type = var.acl
}
