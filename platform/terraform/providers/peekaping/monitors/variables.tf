variable "monitors" {
  description = "List of monitors to create."
  type = list(object({
    name           = string
    type           = string           # "http", "http-keyword", "dns", "push"
    url            = optional(string) # For HTTP/HTTP-keyword monitors.
    method         = optional(string) # For HTTP/HTTP-keyword monitors.
    keyword        = optional(string) # For HTTP-keyword monitors.
    invert_keyword = optional(bool)   # For HTTP-keyword monitors (default false).
    domain         = optional(string) # For DNS monitors.
    resolver_type  = optional(string) # For DNS monitors (A, AAAA, etc.).
    interval       = number           # Interval in seconds between checks.
    max_redirects  = optional(number) # For HTTP/HTTP-keyword monitors (default 3).
    variable_name  = optional(string) # Pre-computed GitHub variable name (e.g., PROD_DB_BACKUP_PING_URL).
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
