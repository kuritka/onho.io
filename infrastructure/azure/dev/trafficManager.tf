
resource "azurerm_resource_group" "rg-global" {
  name = local.shared_rg_name
  location =  var.primary_region
}

// traffic manager lves in the global region, no vnet required
resource "azurerm_traffic_manager_profile" "traffic-manager" {
  name                   = "tfm-${local.suffix}"
  resource_group_name    = local.shared_rg_name
  traffic_routing_method = "Weighted"

  dns_config {
    relative_name = "tfm-${local.suffix}"
    ttl           = 100
  }

  monitor_config {
    protocol = "http"
    port     = 80
    path     = "/"
  }

  tags = var.tags
  depends_on= ["azurerm_resource_group.rg-global"]
}