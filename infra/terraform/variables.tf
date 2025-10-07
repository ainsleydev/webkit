variable "project_name" {
  type        = string
  description = "Name of the client that will be prefixed on all resources"
}

variable "project_title" {
  type        = string
  description = "Nice name of the client that will appear in project settings"
}

variable "resources" {
  type = list(object({
    name     = string
    type     = string
    provider = string
    config   = any
    outputs  = optional(list(string), [])
  }))
  description = "List of resources from the app.json manifest"
  default     = []
}

variable "apps" {
  type = list(object({
    name = string
    title = string
    type = string
    provider = string
    config   = any
    outputs  = optional(list(string), [])
  }))
  description = "List of apps from the app.json manifest"
  default     = []
}

variable "tags" {
  type        = list(string)
  description = "Additional tags to apply to all resources"
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
#
# variable "github_config" {
#   type = object({
#     user  = string
#     repo  = string
#     token = string
#   })
#   description = "Configuration for the Github repo"
# }
#
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
