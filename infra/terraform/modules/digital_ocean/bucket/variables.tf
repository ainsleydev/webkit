variable "name" {
  description = "The name of the S3 bucket"
  type        = string
}

variable "region" {
  description = "In what region the bucket will reside"
  type        = string
  default     = "ams3"
}

variable "acl" {
  description = "The lifecycle policy, defaults to public-read"
  type        = string
  default     = "public-read"
}

variable "tags" {
  description = "List of tags to apply to the resource"
  type        = list(string)
  default     = []
}