# B2 Bucket
# Ref: https://registry.terraform.io/providers/Backblaze/b2/latest/docs/resources/bucket
resource "b2_bucket" "this" {
  bucket_name = var.bucket_name
  bucket_type = var.acl

  dynamic "lifecycle_rules" {
    for_each = var.lifecycle_rule_file_name_prefix != null ? [1] : []
    content {
      days_from_hiding_to_deleting  = var.days_from_hiding_to_deleting
      days_from_uploading_to_hiding = var.days_from_uploading_to_hiding
      file_name_prefix              = var.lifecycle_rule_file_name_prefix
    }
  }
}
