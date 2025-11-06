variable "project_name" {
  description = "Project name for tagging and prefixing"
  type        = string
}

variable "name" {
  description = "App name from manifest"
  type        = string
}

variable "platform_type" {
  description = "Infrastructure type (vm, container, serverless)"
  type        = string

  validation {
    condition     = contains(["vm", "container", "serverless"], var.platform_type)
    error_message = "Platform type must be one of: vm, container, serverless"
  }
}

variable "platform_provider" {
  description = "Platform provider (digitalocean, aws, etc.)"
  type        = string

  validation {
    condition     = contains(["digitalocean"], var.platform_provider)
    error_message = "Provider must be: digitalocean"
  }
}

variable "platform_config" {
  description = "Provider-specific configuration from manifest"
  type        = any
}

variable "app_type" {
  description = "App type (svelte-kit, golang, payload)"
  type        = string

  validation {
    condition     = contains(["svelte-kit", "golang", "payload"], var.app_type)
    error_message = "App type must be one of: svelte-kit, golang, payload"
  }
}

variable "resource_outputs" {
  description = "Outputs from resources module for env var resolution"
  type        = any
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

variable "domains" {
  description = "List of domains from the app manifest"
  type = list(object({
    name     = string
    type     = string
    zone     = optional(string)
    wildcard = optional(bool, false)
  }))
  default = []
}

variable "image_tag" {
  description = "Docker image tag to deploy"
  type        = string
  default     = "latest"
}

variable "github_config" {
  description = "GitHub Container Registry configuration"
  type = object({
    owner  = string
    repo  = string
    token = string
  })
}

variable "do_ssh_key_ids" {
  description = "List of DigitalOcean SSH key IDs to apply to VMs"
  type        = list(string)
  default     = []
}

# Future providers can add their own SSH key variables here:
# variable "hetzner_ssh_key_ids" {
#   description = "List of Hetzner SSH key IDs to apply to VMs"
#   type        = list(string)
#   default     = []
# }

variable "tags" {
  description = "Additional tags"
  type        = list(string)
  default     = []
}

variable "server_user" {
  description = "SSH user for VM deployments"
  type        = string
  default     = "root"
}
