terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.5"
    }
  }
}

resource "local_file" "banner" {
  filename = "/etc/motd-app/banner.txt"
  content  = "Welcome to Acme Checkout -- build 2026.09.14\n"
}
