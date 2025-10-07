variable "name" {
  description = "The name of the DigitalOcean App."
  type        = string
}

variable "domains" {
  description = "List of domains to associate with the app."
  type = list(object({
    name     = string
    type     = string
    zone     = optional(string)
    wildcard = optional(bool, false)
  }))
  default = []
}

variable "region" {
  description = "The region slug where the app should be deployed."
  type        = string
  default     = "lon"
}

variable "service_name" {
  description = "The name of the app service inside App Platform."
  type        = string
}

variable "instance_size_slug" {
  description = "The size slug for the app service instance."
  type        = string
  default     = "apps-s-1vcpu-1gb"
}

variable "instance_count" {
  description = "The number of instances to run for the service."
  type        = number
  default     = 1
}

variable "http_port" {
  description = "The internal HTTP port the service listens on."
  type        = number
  default     = 3000
}

variable "image_tag" {
  description = "The image tag to deploy from the GitHub Container Registry."
  type        = string
  default     = "latest"
}

variable "github_config" {
  description = "GitHub Container Registry config: user, repo, token."
  type = object({
    user  = string
    repo  = string
    token = string
  })
}

variable "health_check_path" {
  description = "The HTTP path for the app health check."
  type        = string
  default     = "/"
}

variable "envs" {
  description = "Dynamic list of environment variables (key, value, scope, type)."
  type = list(object({
    key   = string
    value = string
    scope = optional(string, "RUN_TIME")
    type  = optional(string, "GENERAL")
  }))
  default = []
}
