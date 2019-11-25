resource "azurerm_kubernetes_cluster" "edt-aks-cluster" {

  name                = local.aks_name
  location            = var.region
  resource_group_name = azurerm_resource_group.rg.name
  dns_prefix          = "aks"

  addon_profile {
    http_application_routing {
      enabled = false
    }
  }

  role_based_access_control {
    enabled = false
  }

  agent_pool_profile {
    name           = "agentpool1"
    os_type        = "Linux"
    count           = "1"
    vm_size         = "Standard_DS2_v2"
    os_disk_size_gb = 30
   // vnet_subnet_id  = azurerm_subnet.aks-subnet.id
  }

//  use kubenet if you dont have enough IP's than MUST BE SPECIFIED
//  agent_pool_profile {
//    ...
//    vnet_subnet_id  = azurerm_subnet.aks-subnet.id
//  }
//  network_profile {
//    must be kubenet  vnet_subnet_id  = azurerm_subnet.aks-subnet.id
//    network_plugin = "kubenet"
//  }


//  agent_pool_profile {
//    name           = "agentpool2"
//    os_type        = "Linux"
//    count           = "1"
//    vm_size         = "Standard_DS2_v2"
//    os_disk_size_gb = 30
//
//    vnet_subnet_id  = azurerm_subnet.aks-subnet.id
//  }



  service_principal {
    client_id     = var.client_id
    client_secret = var.client_secret
  }

  network_profile {
    network_plugin = "kubenet"
  }

  tags = merge(var.tags, map(
            "Name", local.aks_name,
          ))
}


//
//this must be specified if kubenet will be used .
// use kubenet only if you have no enough free IP's
//
//resource "azurerm_subnet" "aks-subnet" {
//  name = local.aks_subnet_name
//  resource_group_name = azurerm_resource_group.rg.name
//  virtual_network_name = local.vnet_name
//  address_prefix = local.aks_subnet_prefix
//  service_endpoints = [
//    "Microsoft.Storage",
//    "Microsoft.Sql",
//    "Microsoft.AzureCosmosDB",
//    "Microsoft.EventHub",
//    "Microsoft.KeyVault",
//    "Microsoft.ServiceBus"
//  ]
//}