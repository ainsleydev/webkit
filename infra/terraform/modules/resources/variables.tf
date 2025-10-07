variable "project_name" {
  description = "Project name for tagging"
  type        = string
}

variable "name" {
  description = "Resource name from manifest"
  type        = string
}

variable "type" {
  description = "Resource type (postgres, s3, etc.)"
  type        = string

  validation {
    condition = contains(["postgres", "s3"], var.type)
    error_message = "Type must be one of: postgres, s3"
  }
}

variable "provider" {
  description = "Cloud provider (digitalocean, aws, b2, etc.)"
  type        = string

  validation {
    condition = contains(["digitalocean", "b2"], var.type)
    error_message = "Type must be one of: digitalocean, b2"
  }
}

variable "config" {
  description = "Provider-specific configuration from manifest"
  type        = any
}

variable "tags" {
  description = "Additional tags"
  type = list(string)
  default = []
}
