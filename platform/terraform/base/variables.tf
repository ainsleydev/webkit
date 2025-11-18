variable "project_name" {
  type        = string
  description = "Name of the client that will be prefixed on all resources"

  validation {
    condition     = length(var.project_name) > 0
    error_message = "The project_name variable is required and cannot be empty."
  }
}

variable "project_title" {
  type        = string
  description = "Human-readable title of the project"

  validation {
    condition     = length(var.project_title) > 0
    error_message = "The project_title variable is required and cannot be empty."
  }
}

variable "project_description" {
  type        = string
  description = "Description of the project"
  default     = "Managed by WebKit"
}

variable "environment" {
  type        = string
  description = "The environment the platform is currently running on"

  validation {
    condition     = contains(["development", "staging", "production"], var.environment)
    error_message = "Type must be one of: development, stgaing, production"
  }
}

variable "do_token" {
  type      = string
  sensitive = true
}

variable "do_spaces_access_id" {
  type      = string
  sensitive = true
}

variable "do_spaces_secret_key" {
  type      = string
  sensitive = true
}

variable "digitalocean_project_id" {
  type        = string
  description = "DigitalOcean project ID for preserving manually-added domains. Leave empty on first apply, then populate with project ID from output."
  default     = ""
}

variable "hetzner_token" {
  type        = string
  description = "Hetzner Cloud API token for authentication"
  sensitive   = true
  default     = ""
}

variable "b2_application_key" {
  type      = string
  sensitive = true
}

variable "b2_application_key_id" {
  type      = string
  sensitive = true
}

variable "turso_api_token" {
  type        = string
  description = "Turso token for authentication"
  sensitive   = true
  default     = ""
}

variable "github_token" {
  type      = string
  sensitive = true
}

variable "github_token_classic" {
  type      = string
  sensitive = true
}

variable "resources" {
  type = list(object({
    name              = string
    platform_type     = string
    platform_provider = string
    config            = any
    outputs           = optional(list(string), [])
  }))
  description = "List of resources from the app.json manifest"
  default     = []
}

variable "monitors" {
  type = list(object({
    name              = string
    type              = string # "http", "postgres", "push"
    enabled           = bool
    url               = optional(string)
    method            = optional(string)
    expected_status   = optional(list(number))
    health_check_path = optional(string)
    database_url      = optional(string)
    connection_type   = optional(string)
    expected_interval = optional(number)
    interval          = number
    retry_interval    = number
    max_retries       = number
    upside_down       = optional(bool, false)
    ignore_tls        = optional(bool, false)
  }))
  description = "List of monitors to create in Uptime Kuma"
  default     = []
}

variable "uptime_kuma_notification_ids" {
  type        = list(number)
  description = "List of Uptime Kuma notification IDs to attach to monitors"
  default     = []
}

variable "apps" {
  type = list(object({
    name              = string
    platform_type     = string
    platform_provider = string
    app_type          = string
    path              = optional(string)
    image_tag         = optional(string, "latest")
    config            = any
    domains = optional(list(object({
      name     = string
      type     = string
      zone     = optional(string)
      wildcard = optional(bool, false)
    })), [])
    env_vars = optional(list(object({
      key    = string
      value  = string
      source = string
      type   = optional(string, "GENERAL")
    })), [])
  }))
  description = "List of apps from the app.json manifest"
  default     = []

  validation {
    condition = alltrue([
      for app in var.apps : alltrue([
        for ev in app.env_vars : contains(["GENERAL", "SECRET"], ev.type)
      ])
    ])
    error_message = "Each env_var 'type' must be either 'GENERAL' or 'SECRET'."
  }
}

variable "tags" {
  type        = list(string)
  description = "Additional tags to apply to all resources"
  default     = []
}

variable "digitalocean_ssh_keys" {
  description = "List of SSH key names for DigitalOcean VMs"
  type        = list(string)
  default     = []
}

variable "hetzner_ssh_keys" {
  description = "List of SSH key names for Hetzner VMs"
  type        = list(string)
  default     = []
}

variable "github_config" {
  type = object({
    owner = string
    repo  = string
  })
  description = "Configuration for the Github repo"
  sensitive   = true
}

# --------------------------------------- TODO --------------------------------------- #

# variable "better_stack_token" {
#   type      = string
#   sensitive = true
# }

variable "slack_bot_token" {
  type        = string
  description = "Slack bot token for CI/CD notifications"
  sensitive   = true
}

variable "uptime_kuma_url" {
  type        = string
  description = "Uptime Kuma Web API adapter URL"
  default     = "https://uptime.ainsley.dev"
}

variable "uptime_kuma_username" {
  type        = string
  description = "Uptime Kuma username for authentication"
  sensitive   = true
}

variable "uptime_kuma_password" {
  type        = string
  description = "Uptime Kuma password for authentication"
  sensitive   = true
}

variable "slack_user_token" {
  type        = string
  description = "Slack user token for CI/CD notifications"
  sensitive   = true
}

variable "notifications_webhook_url" {
  type        = string
  description = "Webhook URL for notifications (Slack, Discord, etc.) - sourced from app.json"
  sensitive   = false
  default     = ""
}
