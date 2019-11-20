
/*
 create GRS V2 Premium storage with Hot tier manually in Azure.

 Terraform state will be created automatically after `azure plan` command
*/
terraform {
  backend "azurerm" {
    storage_account_name = "saonhosbx" # Use your own unique name here
    container_name       = "tfstate" # Use your own container name here
    key                  = "sbx.terraform.tfstate" # Add a name to the state file
    resource_group_name  = "rg-onho-state" # Use your own container name here
  }
}
