# Map of learner name => allowed source IP (CIDR, e.g. "203.0.113.5/32").
# One resource group, VNet, NSG, and VM get created per entry.
variable "learners" {
  description = "Map of learner name => allowed source IP (CIDR)"
  type        = map(string)
}

variable "location" {
  description = "Azure region for every learner's resource group"
  type        = string
  default     = "eastus"
}

variable "vm_size" {
  description = "VM size — enough to build + run one lab container at a time"
  type        = string
  default     = "Standard_B2s"
}

variable "admin_username" {
  description = "SSH admin user created on every learner VM"
  type        = string
  default     = "breaklab"
}

variable "ssh_public_key" {
  description = "Public key installed for admin_username on every learner VM"
  type        = string
}

variable "releases_storage_account" {
  description = "Name of the pre-existing shared storage account holding build artifacts"
  type        = string
  default     = "breaklabreleases"
}

variable "releases_resource_group" {
  description = "Resource group of the pre-existing shared storage account"
  type        = string
  default     = "rg-breaklab-releases"
}
