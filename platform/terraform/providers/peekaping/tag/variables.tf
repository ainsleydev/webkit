variable "name" {
  description = "The name of the tag."
  type        = string
}

variable "colour" {
  description = "The hex colour of the tag (e.g., #3B82F6)."
  type        = string
  default     = "#3B82F6"
}

variable "description" {
  description = "A description for the tag."
  type        = string
  default     = ""
}
