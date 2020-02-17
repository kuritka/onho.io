variable "subscription_id" {
  description = "The Azure subscription ID."
}

variable "client_id" {
  description = "The Azure Service Principal app ID."
}

variable "client_secret" {
  description = "The Azure Service Principal password."
}

variable "tenant_id" {
  description = "The Azure Tenant ID."
}

variable "project_shortcut" {
  description = "project shortcut"
  default = "onho"
}


variable environment_short {
  description="Environment shortut used for naming"
  default="sbx"
}

variable "primary_region" {
  description = "The Azure region to create things in."
  default = "westeurope"
}

variable "paired_region" {
  description = "The Azure region to create things in."
  default = "northeurope"
}


variable "tags" {
  type = "map"
  default = {
    Product = "Corporate Systems"
    Environment = "Development"
    ApplicationId = "0"
    Owner = "edt"
  }
}


variable "virtual_network_address_prefix_primary" {
  description = "vnet prefix."
  default="172.17.0.0/26"
}

variable "virtual_network_address_prefix_paired" {
  description = "vnet prefix."
  default="172.17.1.0/26"
}


variable "db_admin_username" {
  description = "Postgre admin name in .tfvars"
}


variable "db_admin_password" {
  description = "Postgre admin password in .tfvars"
}

