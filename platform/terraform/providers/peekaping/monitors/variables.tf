variable "monitors" {
  description = "List of monitors to create."
  type = list(object({
    name     = string
    type     = string           # "http", "dns", "push"
    url      = optional(string) # For HTTP monitors.
    method   = optional(string) # For HTTP monitors.
    domain   = optional(string) # For DNS monitors.
    interval = number           # Interval in seconds between checks.
  }))
  default = []
}

variable "tag_ids" {
  description = "List of tag IDs to apply to monitors."
  type        = list(string)
  default     = []
}

variable "notification_ids" {
  description = "List of notification IDs for alerting."
  type        = list(string)
  default     = []
}

variable "peekaping_endpoint" {
  description = "Peekaping instance endpoint URL (without trailing slash)."
  type        = string
}

variable "defaults" {
  description = "Override default monitor settings."
  type = object({
    timeout          = optional(number)
    http_max_retries = optional(number)
    dns_max_retries  = optional(number)
    push_max_retries = optional(number)
    retry_interval   = optional(number)
    resend_interval  = optional(number)
  })
  default = {}
}
