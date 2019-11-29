
resource "azurerm_container_registry" "acr" {
  name                     = "acr${local.alphanumeric_suffix}"
  resource_group_name      = azurerm_resource_group.rg-global.name
  location                 = var.primary_region
#https://docs.microsoft.com/en-us/azure/container-registry/container-registry-skus
  sku                      = "Basic"
#true = username and pwd in oputputs
  admin_enabled            = true
  tags                     = var.tags
#TODO: add network_rule_set to secure registry and replicas (both sewttings requires Premium sku)
}