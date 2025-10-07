variable "project_name" {
  type        = string
  description = "Name of the client that will be prefixed on all resources"
}

variable "project_title" {
  type        = string
  description = "Nice name of the client that will appear in project settings"
}

variable digital_ocean_api_key {

}

variable "digital_ocean_config" {
  type = object({
    api_key           = string
    spaces_access_key = string
    spaces_secret_key = string
  })
  description = "Configuration for the Digital Ocean provider"
}

variable "github_config" {
  type = object({
    user  = string
    repo  = string
    token = string
  })
  description = "Configuration for the Github repo"
}

variable "back_blaze_config" {
  type = object({
    application_key_id = string
    application_key    = string
  })
  description = "Configuration for BackBlaze B2"
}

variable "better_stack_token" {
  type      = string
  sensitive = true
}

variable "payload_config" {
  type = object({
    secret  = string
    api_key = string
  })
  description = "Configuration for Payload CMS"
}

variable "slack_config" {
  type = object({
    base_user  = string
    bot_token  = string
    user_token = string
  })
  description = "Configuration for Slack"
}
