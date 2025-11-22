variable "title" {
  description = "The title of the status page."
  type        = string
}

variable "description" {
  description = "The description of the status page."
  type        = string
  default     = ""
}

variable "slug" {
  description = "The URL slug for the status page."
  type        = string
}

variable "published" {
  description = "Whether the status page is publicly visible."
  type        = bool
  default     = true
}

variable "theme" {
  description = "The theme of the status page (auto, light, dark)."
  type        = string
  default     = "auto"
}

variable "icon_url" {
  description = "URL to the favicon for the status page."
  type        = string
  default     = null
}

variable "domains" {
  description = "List of custom domains for the status page."
  type        = list(string)
  default     = []
}

variable "monitor_ids" {
  description = "List of monitor IDs to display on the status page."
  type        = list(string)
  default     = []
}
