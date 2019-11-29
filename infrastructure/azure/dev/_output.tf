
//
//output "AKS_SUBNET" {
//  value = module.instances_ne.local.aks_subnet
//}


output "ACR_user" {
  value = azurerm_container_registry.acr.admin_username
}

output "ACR_password" {
  value = azurerm_container_registry.acr.admin_password
}

output "ACR_login_server" {
  value = azurerm_container_registry.acr.login_server
}