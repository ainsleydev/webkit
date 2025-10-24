variable "bucket_name" {
  description = "Name of the S3 bucket"
  type        = string
}

variable "acl" {
  description = "ACL of the bucket, defaults to allPrivate"
  type        = string
  default     = "allPrivate"
}

variable "days_from_hiding_to_deleting" {
  description = "How long to keep file versions that are not the current version (in days)"
  type        = number
  default     = null
}

variable "days_from_uploading_to_hiding" {
  description = "Causes files to be hidden automatically after the given number of days"
  type        = number
  default     = null
}

variable "lifecycle_rule_file_name_prefix" {
  description = "Specifies which files in the bucket the lifecycle rule applies to"
  type        = string
  default     = null
}
