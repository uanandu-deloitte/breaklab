# Terraform configuration from registry.terraform.io
terraform {
  required_version = ">= 1.9"
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.90" # floor only — latest is 2.106.0 as of writing, run
      # `terraform init -upgrade` to pick up newer 2.x releases
    }
  }
}

# Provider configuration for DigitalOcean
provider "digitalocean" {
  # Configuration options
  token = var.do_token
}