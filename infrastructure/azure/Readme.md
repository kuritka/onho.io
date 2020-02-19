## Init infrastructure

 create GRS V2 Premium storage with Hot tier manually in Azure.


 look into `_backend.tf` to know resource group, storage account name etc..

```bash
terraform init -backend-config=terraform.tfvars
```

 Terraform state will be created automatically after `terraform plan -out out.plan` command
 
## Deploy infrastructure

You will find terraform.tfvars within directory and cannot be shared by github!

```bash
terraform plan -out out.plan
```

```bash
terraform apply out.plan
```

```bash
terraform output  
```


## after you deploy

```bash 
 docker login acronhosbx.azurecr.io
```


```bash
az account set --subscription S_Sandbox_SBX

# S_Sandbox_SBX must be True
az account list --output table

az aks get-credentials --resource-group rg-onho-sbx-westeurope --name aks-onho-sbx-westeurope

kubectl config use-context aks-onho-sbx-westeurope
```


## keep protected

If you loose terraform state you will loose whole infrastructure. Especially database could be an issue

Never share secrets and state on github. Add these values into `.gitignore`
```bash

**/.terraform/*

# .tfstate files
**/*.tfstate
**/*.tfstate.*
**/terraform.tfvars
**/*.plan

```