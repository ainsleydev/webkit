variable "project_name" {
  description = "Project name for monitor naming"
  type        = string
}

variable "monitors" {
  description = "List of monitors to create. Defaults are applied based on monitor type."
  type = list(object({
    name   = string
    type   = string # "http", "postgres", "push"
    url    = optional(string)
    method = optional(string)
  }))
  default = []
}

variable "uptime_kuma_url" {
  type        = string
  description = "Uptime Kuma Web API adapter URL"
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
