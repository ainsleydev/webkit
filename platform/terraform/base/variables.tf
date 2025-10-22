variable "project_name" {
  type        = string
  description = "Name of the client that will be prefixed on all resources"
}

# variable "project_title" {
#   type        = string
#   description = "Nice name of the client that will appear in project settings"
# }

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

variable "b2_application_key" {
  type      = string
  sensitive = true
}

variable "b2_application_key_id" {
  type      = string
  sensitive = true
}

variable "github_token" {
  type = string
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

variable "ssh_keys" {
  description = "List of SSH key names to apply to droplets"
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
#
# variable "slack_config" {
#   type = object({
#     base_user  = string
#     bot_token  = string
#     user_token = string
#   })
#   description = "Configuration for Slack"
# }
