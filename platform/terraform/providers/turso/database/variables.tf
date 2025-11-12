variable "name" {
  description = "The name of the Turso database"
  type        = string
}

variable "organisation" {
  description = "The Turso organisation name"
  type        = string
}

variable "group" {
  description = "The Turso group name (e.g., 'default' or a custom group)"
  type        = string
  default     = "default"
}

variable "size_limit" {
  description = "Optional size limit for the database"
  type        = string
  default     = null
}

variable "tags" {
  description = "List of tags to apply to the resource"
  type        = list(string)
  default     = []
}
