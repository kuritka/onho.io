//#https from internal network
//# test like curl -kv https://10.121.36.101

resource "azurerm_public_ip" "pip" {
  name                = local.gw_pip_name
  location            = var.region
  resource_group_name = azurerm_resource_group.rg.name
  allocation_method   = "Dynamic"
  tags = merge(var.tags, map(
    "Name", local.gw_pip_name,
  ))
}

data "azurerm_public_ip" "pip" {
  name                = azurerm_public_ip.pip.name
  resource_group_name = azurerm_resource_group.rg.name
}


resource "azurerm_application_gateway" "network" {
  name                = local.gw_name
  resource_group_name = azurerm_resource_group.rg.name
  location            = var.region

  sku {
    //standard supports private ip only
    name     = "Standard_Medium"
    tier     = "Standard"
    capacity = 1
  }

  gateway_ip_configuration {
    name      = "appGatewayIpConfig"
    subnet_id = azurerm_subnet.appgwsubnetprivate.id
  }

  frontend_port {
    name = local.frontend_port_name
    port = 443
  }

  frontend_ip_configuration {
    name                 = local.frontend_ip_configuration_name
//    subnet_id            = azurerm_subnet.appgwsubnetprivate.id
//    private_ip_address   = cidrhost(local.gw_private_subnet,5)
//    private_ip_address_allocation = "Static"
    public_ip_address_id = azurerm_public_ip.pip.id
  }

  backend_address_pool {
    name = local.backend_address_pool_name
    ip_addresses = [local.aks_public_endpoint_private_ip]
  }

  backend_http_settings {
    name                  = local.http_setting_name
    cookie_based_affinity = "Disabled"
    port                  = 80
    protocol              = "Http"
    request_timeout       = 1
  }

  http_listener {
    name                           = local.listener_name
    frontend_ip_configuration_name = local.frontend_ip_configuration_name
    frontend_port_name             = local.frontend_port_name
    protocol                       = "https"
    ssl_certificate_name           = "ssl-certificate-1"
  }

  ssl_certificate {
    name                            = "ssl-certificate-1"
    data                            = local.wildcard_cert1_domain_cert_data
    password                        = ""
  }

  request_routing_rule {
    name                       = local.request_routing_rule_name
    rule_type                  = "Basic"
    http_listener_name         = local.listener_name
    backend_address_pool_name  = local.backend_address_pool_name
    backend_http_settings_name = local.http_setting_name
  }

  tags = merge(var.tags, map(
            "Name", local.gw_name,
          ))
}
