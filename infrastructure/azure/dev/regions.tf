module "instances_we" {
  source                  = "../regions"
  //custom
  region                   = var.primary_region
  vnet_addr_prefix         = var.virtual_network_address_prefix_primary

  //common
  project_shortcut         = var.project_shortcut
  environment              = var.environment_short
  tags                     = var.tags
  client_id                = var.client_id
  client_secret            = var.client_secret
  tfm_name                 = local.tfm_name
  shared_rg_name           = local.shared_rg_name
}

// ENABLE THIS TO ENABLE PAIRED REGION
//module "instances_ne" {
//  source                  = "../regions"
//  //custom
//  region                   = var.paired_region
//  vnet_addr_prefix         = var.virtual_network_address_prefix_paired
//
//  //common
//  project_shortcut         = var.project_shortcut
//  environment              = var.environment_short
//  tags                     = var.tags
//  client_id                = var.client_id
//  client_secret            = var.client_secret
//  tfm_name                 = local.tfm_name
//  shared_rg_name           = local.shared_rg_name
//}