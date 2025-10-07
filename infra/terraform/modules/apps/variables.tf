variable "project_name" {
  description = "Project name for tagging and prefixing"
  type        = string
}

variable "name" {
  description = "App name from manifest"
  type        = string
}

variable "app_type" {
  description = "App type (sveltekit, golang, payload-cms)"
  type        = string

  validation {
    condition     = contains(["sveltekit", "golang", "payload-cms"], var.app_type)
    error_message = "App type must be one of: sveltekit, golang, payload-cms"
  }
}

variable "cloud_provider" {
  description = "Cloud provider (digitalocean, aws, etc.)"
  type        = string

  validation {
    condition     = contains(["digitalocean"], var.cloud_provider)
    error_message = "Provider must be: digitalocean"
  }
}

variable "infra_type" {
  description = "Infrastructure type (vm, container, serverless)"
  type        = string

  validation {
    condition     = contains(["vm", "container", "serverless"], var.infra_type)
    error_message = "Infra type must be one of: vm, container, serverless"
  }
}

variable "infra_config" {
  description = "Provider-specific infrastructure configuration from manifest"
  type        = any
}

variable "image_tag" {
  description = "Docker image tag to deploy"
  type        = string
  default     = "latest"
}

variable "github_config" {
  description = "GitHub Container Registry configuration"
  type = object({
    user  = string
    repo  = string
    token = string
  })
}

variable "env_vars" {
  description = "Environment variables for the app"
  type = list(object({
    key   = string
    value = string
    scope = optional(string, "RUN_TIME")
    type  = optional(string, "GENERAL")
  }))
  default = []
}

variable "user_ssh_key_name" {
  description = "Default SSH key name for droplets"
  type        = string
  default     = ""
}

variable "tags" {
  description = "Additional tags"
  type        = list(string)
  default     = []
}
