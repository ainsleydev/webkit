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
