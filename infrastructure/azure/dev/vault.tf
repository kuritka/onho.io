
resource "azurerm_key_vault" "example" {
  name                        = "vlt${local.alphanumeric_suffix}"
  location                    = var.primary_region
  resource_group_name         = azurerm_resource_group.rg-global.name
  enabled_for_disk_encryption = true
  tenant_id                   = var.tenant_id

  sku_name = "premium"

  access_policy {
    tenant_id = var.tenant_id
    object_id = var.client_id

    key_permissions = [
      "get",
      "create",
      "list",
      "delete",
      "update",
    ]

    secret_permissions = [
      "get",
      "list",
      "delete",
      "set",
    ]

    storage_permissions = [
      "get",
      "list",
      "delete",
      "set",
      "update",
      "regeneratekey",
    ]
  }

  network_acls {
    default_action = "Deny"
    bypass         = "AzureServices"
  }

// for internet access
//  network_acls {
//    default_action = "Allow"
//  }

  tags = var.tags
}