locals {
  suffix = "${var.project_shortcut}-${var.environment}-${var.region}"
  vnet_name = "vnet-${local.suffix}"

  resource_group_name = "rg-${local.suffix}"

  tme_name = "tme-${local.suffix}"

  addr_prefix_start = cidrsubnet(var.vnet_addr_prefix,  lookup(var.newbit_size,"size_256"),0)

  aks_name = "aks-${local.suffix}"

  aks_subnet_prefix = cidrsubnet(local.addr_prefix_start,  lookup(var.segments_size,"size_32"),0)

  aks_subnet_name = "sn-aks-${local.suffix}"

  aks_public_endpoint_private_ip = cidrhost(local.aks_subnet_prefix, 5)

  gw_name = "agw-${local.suffix}"

  gw_pip_name = "pip-${local.suffix}"

  //once private subnet range is set, you cannot change it without recreate resources related to that subnet!
  gw_private_subnet = cidrsubnet(local.addr_prefix_start,  lookup(var.segments_size,"size_16"),3)

  gw_private_subnet_name = "sn-agw-private-${local.suffix}"

  backend_address_pool_name      = "${local.vnet_name}-beap"
  frontend_port_name             = "${local.vnet_name}-feport"
  frontend_ip_configuration_name = "${local.vnet_name}-feip"
  http_setting_name              = "${local.vnet_name}-be-htst"
  listener_name                  = "${local.vnet_name}-httpslstn"
  request_routing_rule_name      = "${local.vnet_name}-rqrt"

  //todo: use vault instead providing certificate...
  wildcard_cert1_domain_cert_data = filebase64("../regions/certificates/cert.pfx")
}