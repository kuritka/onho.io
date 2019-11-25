

You will find terraform.tfvars within eml and cannot be shared by github!

```bash
terraform init -backend-config=terraform.tfvars
```

```bash
terraform plan -out out.plan
```

```bash
terraform apply out.plan
```
