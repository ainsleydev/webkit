#
# Shared Peekaping Tags Setup
#
# This is a one-time setup file to create shared tags that are used across
# all webkit repos. Run this once to create the shared tags, then comment out
# or delete this file.
#
# Usage:
#   1. Ensure you have Peekaping credentials configured
#   2. Run: terraform apply -target=peekaping_tag.shared_webkit -target=peekaping_tag.shared_production -target=peekaping_tag.shared_staging -target=peekaping_tag.shared_development
#   3. After tags are created, you can safely delete or comment out this file
#

terraform {
  required_providers {
    peekaping = {
      source  = "tafaust/peekaping"
      version = "~> 0.1.1"
    }
  }
}

# NOTE: Configure the provider with your Peekaping credentials
# provider "peekaping" {
#   endpoint = "https://your-peekaping-instance.com"
#   api_key  = var.peekaping_api_key
# }

#
# Shared WebKit Tag
#
resource "peekaping_tag" "shared_webkit" {
  name        = "WebKit"
  color       = "#10B981" # Green
  description = "Managed by WebKit infrastructure"
}

#
# Shared Environment Tags
#
resource "peekaping_tag" "shared_production" {
  name        = "Production"
  color       = "#EF4444" # Red
  description = "Production environment"
}

resource "peekaping_tag" "shared_staging" {
  name        = "Staging"
  color       = "#F59E0B" # Orange
  description = "Staging environment"
}

resource "peekaping_tag" "shared_development" {
  name        = "Development"
  color       = "#3B82F6" # Blue
  description = "Development environment"
}

resource "peekaping_tag" "shared_test" {
  name        = "Test"
  color       = "#8B5CF6" # Purple
  description = "Test environment"
}
