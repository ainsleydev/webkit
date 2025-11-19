terraform {
  required_providers {
    uptimekuma = {
      source  = "ehealth-co-id/uptimekuma"
      version = "~> 0.0.2"
      configuration_aliases = [uptimekuma]
    }
  }
}
