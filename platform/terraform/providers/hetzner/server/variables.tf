variable "name" {
  type        = string
  description = "Server name"
}

variable "server_type" {
  type        = string
  description = "Hetzner server type (size)"
  default     = "cx22"
}

variable "location" {
  type        = string
  description = "Hetzner location (region)"
  default     = "nbg1"
}

variable "ssh_key_ids" {
  type        = list(string)
  description = "List of Hetzner SSH key IDs or names to apply to the server"
  default     = []
}

variable "tags" {
  type        = list(string)
  description = "Tags to apply to the server (converted to labels)"
  default     = []
}

variable "server_user" {
  type        = string
  description = "SSH user for server access"
  default     = "root"
}
