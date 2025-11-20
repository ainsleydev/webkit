variable "name" {
  description = "Monitor name"
  type        = string
}

variable "type" {
  description = "Monitor type (http, dns, push)"
  type        = string
}

variable "config" {
  description = "Monitor configuration (JSON-encoded)"
  type        = string
}

variable "interval" {
  description = "Check interval in seconds"
  type        = number
  default     = 60
}

variable "timeout" {
  description = "Timeout in seconds"
  type        = number
  default     = 30
}

variable "max_retries" {
  description = "Maximum number of retries"
  type        = number
  default     = 3
}

variable "retry_interval" {
  description = "Retry interval in seconds"
  type        = number
  default     = 60
}

variable "resend_interval" {
  description = "Resend notification interval in minutes"
  type        = number
  default     = 10
}

variable "active" {
  description = "Whether the monitor is active"
  type        = bool
  default     = true
}

variable "notification_ids" {
  description = "List of notification IDs to use for alerts"
  type        = list(string)
  default     = []
}

variable "tag_ids" {
  description = "List of tag IDs to apply to this monitor"
  type        = list(string)
  default     = []
}
