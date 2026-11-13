# One resource group per learner — teardown is `terraform destroy`, full stop.
resource "azurerm_resource_group" "learner" {
  for_each = var.learners
  name     = "rg-breaklab-${each.key}"
  location = var.location
}

resource "azurerm_virtual_network" "learner" {
  for_each            = var.learners
  name                = "vnet-${each.key}"
  resource_group_name = azurerm_resource_group.learner[each.key].name
  location            = var.location
  address_space       = ["10.42.0.0/24"]
}

resource "azurerm_subnet" "learner" {
  for_each             = var.learners
  name                 = "default"
  resource_group_name  = azurerm_resource_group.learner[each.key].name
  virtual_network_name = azurerm_virtual_network.learner[each.key].name
  address_prefixes     = ["10.42.0.0/28"]
}

# Inbound 8787 (the UI) and 22 (SSH) are both scoped to that one learner's
# IP only — breaklab serve has no authentication, the NSG is the access
# control, not the --host flag on the VM.
resource "azurerm_network_security_group" "learner" {
  for_each            = var.learners
  name                = "nsg-${each.key}"
  resource_group_name = azurerm_resource_group.learner[each.key].name
  location            = var.location

  security_rule {
    name                       = "allow-breaklab-ui"
    priority                   = 100
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "8787"
    source_address_prefix      = each.value
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "allow-ssh"
    priority                   = 110
    direction                  = "Inbound"
    access                     = "Allow"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = each.value
    destination_address_prefix = "*"
  }
}

resource "azurerm_public_ip" "learner" {
  for_each            = var.learners
  name                = "pip-${each.key}"
  resource_group_name = azurerm_resource_group.learner[each.key].name
  location            = var.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_network_interface" "learner" {
  for_each            = var.learners
  name                = "nic-${each.key}"
  resource_group_name = azurerm_resource_group.learner[each.key].name
  location            = var.location

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.learner[each.key].id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.learner[each.key].id
  }
}

resource "azurerm_network_interface_security_group_association" "learner" {
  for_each                  = var.learners
  network_interface_id      = azurerm_network_interface.learner[each.key].id
  network_security_group_id = azurerm_network_security_group.learner[each.key].id
}

resource "azurerm_linux_virtual_machine" "learner" {
  for_each            = var.learners
  name                = "vm-breaklab-${each.key}"
  resource_group_name = azurerm_resource_group.learner[each.key].name
  location            = var.location
  size                = var.vm_size
  admin_username      = var.admin_username

  network_interface_ids = [azurerm_network_interface.learner[each.key].id]

  admin_ssh_key {
    username   = var.admin_username
    public_key = var.ssh_public_key
  }

  os_disk {
    caching              = "ReadWrite"
    storage_account_type = "Standard_LRS"
  }

  source_image_reference {
    publisher = "Canonical"
    offer     = "ubuntu-24_04-lts"
    sku       = "server"
    version   = "latest"
  }

  # SystemAssigned identity is what lets cloud-init pull the build artifact
  # from the shared storage account below with no secret baked into the image.
  identity {
    type = "SystemAssigned"
  }

  custom_data = base64encode(templatefile("${path.module}/cloud-init.yaml.tftpl", {
    admin_username = var.admin_username
  }))
}

data "azurerm_storage_account" "releases" {
  name                = var.releases_storage_account
  resource_group_name = var.releases_resource_group
}

resource "azurerm_role_assignment" "learner_pull" {
  for_each             = var.learners
  scope                = data.azurerm_storage_account.releases.id
  role_definition_name = "Storage Blob Data Reader"
  principal_id         = azurerm_linux_virtual_machine.learner[each.key].identity[0].principal_id
}

output "learner_ips" {
  value = { for k, v in azurerm_public_ip.learner : k => v.ip_address }
}
