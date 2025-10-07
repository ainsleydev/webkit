variable "name" {
  description = "The name of the database cluster, should be suffix of _db"
  type = string
}

variable "pg_version" {
  type = string
  default = "17"
  description = "The Postgres version of the database, defaults to 17"
}

variable "size" {
  type = string
  default = "db-s-1vcpu-1gb"
}

variable "region" {
  type = string
  default = "lon1"
}

variable "node_count" {
  type = number
  default = 1
}

variable "allowed_droplet_ips" {
  type = list(string)
  default = []
}

variable "allowed_ips_addr" {
  type = list(string)
  default = []
}

variable "tags" {
  description = "List of tags to apply to the resource"
  type        = list(string)
  default     = []
}