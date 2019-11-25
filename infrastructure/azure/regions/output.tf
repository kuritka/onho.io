output "resource_group" {
  value = local.resource_group_name
}

output "public_ip_address" {
  value = azurerm_public_ip.pip.ip_address
}

