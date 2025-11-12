variable "project_name" {
  description = "Project name for tagging"
  type        = string
}

variable "name" {
  description = "Resource name from manifest"
  type        = string
}

variable "platform_type" {
  description = "Resource type (postgres, s3, sqlite, etc.)"
  type        = string

  validation {
    condition     = contains(["postgres", "s3", "sqlite"], var.platform_type)
    error_message = "Type must be one of: postgres, s3, sqlite"
  }
}

variable "platform_provider" {
  description = "Platform provider (digitalocean, aws, b2, turso, etc.)"
  type        = string

  validation {
    condition     = contains(["digitalocean", "b2", "turso"], var.platform_provider)
    error_message = "Provider must be one of: digitalocean, b2, turso"
  }
}

variable "platform_config" {
  description = "Provider-specific configuration from manifest"
  type        = any
}

variable "tags" {
  description = "Additional tags"
  type        = list(string)
  default     = []
}
