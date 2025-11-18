variable "project_name" {
  description = "Project name for monitor naming"
  type        = string
}

variable "monitors" {
  description = "List of monitors to create"
  type = list(object({
    name              = string
    type              = string # "http", "postgres", "push"
    enabled           = bool

    # HTTP fields
    url               = optional(string)
    method            = optional(string)
    expected_status   = optional(list(number))
    health_check_path = optional(string)

    # Database fields
    database_url    = optional(string)
    connection_type = optional(string)

    # Push/Heartbeat fields
    expected_interval = optional(number)

    # Common fields
    interval       = number
    retry_interval = number
    max_retries    = number
    upside_down    = optional(bool, false)
    ignore_tls     = optional(bool, false)
  }))
  default = []
}

variable "notification_ids" {
  description = "List of notification IDs to attach to monitors"
  type        = list(number)
  default     = []
}
