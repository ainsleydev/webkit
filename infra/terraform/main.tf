terraform {
  # backend "s3" {
  #   bucket                      = "bucket"
  #   key                         = "temp/terraform.tfstate"
  #   region                      = "eu-central-003"
  #   skip_credentials_validation = true
  #   skip_region_validation      = true
  #   skip_requesting_account_id  = true
  #   use_path_style              = true
  #
  #   endpoints = {
  #     s3 = "https://s3.eu-central-003.backblazeb2.com"
  #   }
  # }

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
    logtail = {
      source  = "BetterStackHQ/logtail"
      version = ">= 0.1.0"
    }
    github = {
      source  = "integrations/github"
      version = "~> 5.0"
    }
    b2 = {
      source = "Backblaze/b2"
    }
    slack = {
      source  = "pablovarela/slack"
      version = "~> 1.0"
    }
  }
}
