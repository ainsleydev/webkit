variable "name" {
  description = "The name of the Droplet."
  type        = string
}

variable "droplet_size" {
  description = "The instance size slug, defaults to smallest"
  type        = string
  default     = "s-1vcpu-1gb"
}

variable "droplet_region" {
  description = "The region of the droplet, defaults to London"
  type        = string
  default     = "lon1"
}

variable "ssh_keys" {
  description = "List of SSH key names to apply to the droplet"
  type        = list(string)
  default     = []
}

variable "tags" {
  description = "List of tags to apply to the resource"
  type        = list(string)
  default     = []
}

variable "server_user" {
  description = "The SSH user for the server"
  type        = string
  default     = "root"
}
