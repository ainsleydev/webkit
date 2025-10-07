variable "bucket_name" {
  description = "Name of the S3 bucket"
  type        = string
}

variable "acl" {
  description = "ACL of the bucket, defaults to allPrivate"
  type        = string
  default     = "allPrivate"
}
