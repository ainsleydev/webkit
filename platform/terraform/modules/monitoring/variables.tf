variable "project_name" {
  description = "Project name for monitor naming"
  type        = string
}

variable "project_title" {
  description = "Human-readable project title for display"
  type        = string
}

variable "environment" {
  description = "Environment name (production, staging, development)"
  type        = string
}

variable "monitors" {
  description = "List of monitors to create."
  type = list(object({
    name          = string
    type          = string
    url           = optional(string)
    method        = optional(string)
    domain        = optional(string)
    interval      = number
    max_redirects = optional(number)
    variable_name = optional(string)
  }))
  default = []
}

variable "slack_webhook_url" {
  description = "Slack webhook URL for notifications"
  type        = string
  default     = ""
}

variable "brand_primary_color" {
  description = "Primary brand colour for tags (optional)"
  type        = string
  default     = null
}

variable "brand_logo_url" {
  description = "Logo URL for status page (optional)"
  type        = string
  default     = null
}

variable "brand_icon_url" {
  description = "Icon URL for status page favicon (optional)"
  type        = string
  default     = null
}

variable "status_page_domain" {
  description = "Custom domain for status page (e.g., status.example.com)"
  type        = string
  default     = null
}

variable "status_page_slug" {
  description = "Custom slug for status page URL. If not set, defaults to kebab-case project name."
  type        = string
  default     = null
}

variable "status_page_theme" {
  description = "Theme for status page (auto, light, dark). Defaults to auto."
  type        = string
  default     = null
}

variable "peekaping_endpoint" {
  description = "Peekaping instance endpoint URL (without trailing slash)"
  type        = string
}
