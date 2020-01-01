## HTTPS

Not all functions are allowed on public IP with HTTP. i.e. camera. Visit http://permission.site/ to see more details


### DEV
For securing traffic we use self sign certificate and istio ingress .
To generate `self-sign` certificate run this: 
```bash
./scripts/cert-generator.sh self-sign ./infrastructure/certificates/dev/dev.onho.cnf ../onho-certs/dev
```
