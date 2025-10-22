# B2 Bucket
# Ref: https://registry.terraform.io/providers/Backblaze/b2/latest/docs/resources/bucket
resource "b2_bucket" "this" {
  bucket_name = var.bucket_name
  bucket_type = var.acl
}
