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

# Pure cost control — shutting a VM down never touches a learner's
# containers or progress/postmortems files, which live on the OS disk,
# not in memory (docs/devops-guide-azure.md §7). Off by default; turn
# on once a cohort runs on a predictable schedule rather than
# speculatively.
variable "auto_shutdown" {
  description = "Nightly auto-shutdown schedule for every learner VM."
  type = object({
    enabled  = bool
    time     = string # "HHmm", 24h, in the VM's own local time
    timezone = string # Azure's accepted tz name, e.g. "UTC" or "Pacific Standard Time"
  })
  default = {
    enabled  = false
    time     = "1900"
    timezone = "UTC"
  }
}
