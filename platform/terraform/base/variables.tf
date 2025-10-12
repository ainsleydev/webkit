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

variable "resources" {
  type = list(object({
    name     = string
    type     = string
    provider = string
    config   = map(any)
    outputs  = optional(list(string), [])
  }))
  description = "List of resources from the app.json manifest"
  default     = []
}

variable "apps" {
  type = list(object({
    name      = string
    type      = string
    path      = optional(string)
    image_tag = optional(string, "latest")
    infra = object({
      provider = string
      type     = string
      config   = map(any)
    })
    domains = optional(list(object({
      name     = string
      type     = string
      zone     = optional(string)
      wildcard = optional(bool, false)
    })), [])
    env_vars = optional(list(object({
      key   = string
      value = string
      scope = optional(string, "RUN_TIME")
      type  = optional(string, "GENERAL")
    })), [])
  }))
  description = "List of apps from the app.json manifest"
  default     = []
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

# variable "digital_ocean_config" {
#   type = object({
#     api_key           = string
#     spaces_access_key = string
#     spaces_secret_key = string
#   })
#   description = "Configuration for the Digital Ocean provider"
# }

variable "github_config" {
  type = object({
    user  = string
    repo  = string
    token = string
  })
  description = "Configuration for the Github repo"
}


# variable "back_blaze_config" {
#   type = object({
#     application_key_id = string
#     application_key    = string
#   })
#   description = "Configuration for BackBlaze B2"
# }
#
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
