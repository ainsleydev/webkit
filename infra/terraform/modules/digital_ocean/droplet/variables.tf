variable "name" {
  description = "The name of the Droplet."
  type        = string
}

variable "user_ssh_key_name" {
  description = "The name of the existing SSH key in DigitalOcean"
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

variable "tags" {
  description = "List of tags to apply to the resource"
  type        = list(string)
  default     = []
}