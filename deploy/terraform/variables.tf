# Token for DigitalOcean API access
variable "do_token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}

# SSH key fingerprint for the droplet
variable "ssh_key_fingerprint" {
  description = "SSH key fingerprint for the droplet"
  type        = string
}

# Region for the DigitalOcean droplet
variable "region" {
  description = "DigitalOcean region for the droplet"
  type        = string
  default     = "nyc3"
}

# Name of the DigitalOcean droplet
variable "droplet_name" {
  description = "Name of the DigitalOcean droplet"
  type        = string
  default     = "breaklab-droplet"
}

# Droplet size for the DigitalOcean droplet
variable "droplet_size" {
  description = "Size of the DigitalOcean droplet"
  type        = string
  default = "s-2vcpu-2gb" # ~$18/mo — the whole budget
}

# IP CIDR for SSH access to the droplet
variable "my_ip_cidr" {
    description = "Your public IP address in CIDR notation"
    type        = string # "203.0.113.5/32" — restricts SSH access to a single IP address
}