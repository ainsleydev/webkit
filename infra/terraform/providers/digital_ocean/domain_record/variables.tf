variable "domain" {
  description = "The domain to add the record to, can be a domain ID"
  type        = string
}

variable "name" {
  description = "The name of the record"
  type = string
}

variable "value" {
  description = "The value of the record"
  type        = string
}

variable "type" {
  description = "The type of DNS record (A, AAAA, CNAME, etc.)"
  type        = string
  default     = "A"

  validation {
    condition     = contains(["A", "AAAA", "CNAME", "MX", "TXT", "SRV", "NS"], var.type)
    error_message = "Type must be a valid DNS record type."
  }
}
