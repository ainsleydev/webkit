terraform {
  required_providers {
    uptimekuma = {
      source  = "kill3r-queen/uptimekuma"
      version = "~> 0.0.12"
    }
  }
}

provider "uptimekuma" {
  base_url       = var.uptime_kuma_url
  username       = var.uptime_kuma_username
  password       = var.uptime_kuma_password
  insecure_https = false
}
