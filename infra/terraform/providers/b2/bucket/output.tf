output "id" {
  description = "The ID of the S3 bucket"
  value       = b2_bucket.this.id
}

output "bucket_name" {
  description = "The name of the S3 bucket"
  value       = b2_bucket.this.bucket_name
}

output "bucket_info" {
  description = "The metadata associated with the S3 bucket"
  value       = b2_bucket.this.bucket_info
}
