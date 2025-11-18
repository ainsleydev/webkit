variable "name" {
  type        = string
  description = "Volume name"
}

variable "size" {
  type        = number
  description = "Size of the volume in GB (min: 10)"
  default     = 10

  validation {
    condition     = var.size >= 10
    error_message = "Volume size must be at least 10 GB"
  }
}

variable "location" {
  type        = string
  description = "Hetzner location (region)"
  default     = "nbg1"
}

variable "server_id" {
  type        = number
  description = "Hetzner server ID to attach the volume to"
}

variable "format" {
  type        = string
  description = "Filesystem format for the volume"
  default     = "ext4"

  validation {
    condition     = contains(["ext4", "xfs"], var.format)
    error_message = "Format must be either ext4 or xfs"
  }
}

variable "automount" {
  type        = bool
  description = "Automount the volume when attaching it"
  default     = true
}

variable "tags" {
  type        = list(string)
  description = "Tags to apply to the volume (converted to labels)"
  default     = []
}
