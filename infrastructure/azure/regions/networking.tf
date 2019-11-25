resource "azurerm_resource_group" "rg" {
  name = local.resource_group_name
  location =  var.region
}

resource "azurerm_virtual_network" "vnet" {
  name                = local.vnet_name
  location            = var.region
  resource_group_name = azurerm_resource_group.rg.name
  address_space       = [var.vnet_addr_prefix]

//  subnet {
//    name           = local.gw_private_subnet_name
//    address_prefix = local.gw_private_subnet
//  }

  tags = var.tags
}

resource "azurerm_subnet" "appgwsubnetprivate" {
  name                 = local.gw_private_subnet_name
  virtual_network_name = azurerm_virtual_network.vnet.name
  resource_group_name  = azurerm_resource_group.rg.name
  address_prefix       = local.gw_private_subnet
}




