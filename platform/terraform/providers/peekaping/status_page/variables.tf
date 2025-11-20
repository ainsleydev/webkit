variable "title" {
  description = "Status page title"
  type        = string
}

variable "description" {
  description = "Status page description"
  type        = string
}

variable "slug" {
  description = "URL slug for the status page"
  type        = string
}

variable "published" {
  description = "Whether the status page is published"
  type        = bool
  default     = true
}

variable "theme" {
  description = "Status page theme (light, dark, auto)"
  type        = string
  default     = "auto"
}

variable "icon" {
  description = "Icon URL for the status page"
  type        = string
  default     = null
}
