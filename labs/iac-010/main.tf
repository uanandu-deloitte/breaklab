terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "~> 2.5"
    }
  }
}

resource "local_file" "legal_notice" {
  filename = "/var/www/legal/notice.txt"
  content  = "All rights reserved, Acme Corp.\n"
}
