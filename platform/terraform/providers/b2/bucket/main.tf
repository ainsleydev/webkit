# B2 Bucket
# Ref: https://registry.terraform.io/providers/Backblaze/b2/latest/docs/resources/bucket
resource "b2_bucket" "this" {
  bucket_name = var.bucket_name
  bucket_type = var.acl

  lifecycle_rule {
    days_from_hiding_to_deleting = 1
    days_from_uploading_to_hiding = 0
    file_name_prefix              = ""
  }
}
