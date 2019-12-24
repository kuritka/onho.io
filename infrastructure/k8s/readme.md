

## integrating Azure Vault with AKS

kubectl create -f https://raw.githubusercontent.com/Azure/kubernetes-keyvault-flexvol/master/deployment/kv-flexvol-installer.yaml

kubectl create secret generic kvcreds --from-literal clientid=<CLIENTID> --from-literal clientsecret=<CLIENTSECRET> --type=azure/kv

kubectl describe secret/kvcreds -v 7

## secrets

secrets are stored within secrets.yaml . Secrets must be encrypted before push into github.
Use following command sto encrypt/decrypt secrets file. To perform this task we wil use ansible but 
encrypt file by zipping is also possible.

```bash

ansible-vault encrypt ./secrets.yaml

#make changes within secrets file
kubectl apply -f ./secrets.yaml

kubectl describe  secrets/onho-secrets -n onho-dev

ansible-vault decrypt ./secrets.yaml
#(.Z....9)

#now you can push back to github
```