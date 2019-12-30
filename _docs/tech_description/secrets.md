## Decrypting ansible secret file

Secrets are stored in `secrets.yaml` and must be  encrypted before push into github

First at all, to decrypt ansible 

```bash
ansible-vault decrypt secret.yaml
```
and enter `.Z...9`


Particular values are base64 decoded. To decode make this:

```bash
echo 'bGludXhoaW50LmNvbQo=' | base64 -decode
```

and to encode 
```bash
echo -n 'linuxhint.com' | base64 -w0
```
It is necessary to run this command in linux bash. 
Running such command within linux power-shell causes an issues by putting `\n` characters.


to encrypt run 


```bash
ansible-vault encrypt secret.yaml
``` 

## generate certificate

### DEV 
generate self sign certificate
```bash
./scripts/cert-generator.sh self-sign ./infrastructure/certificates/dev/dev.onho.cnf ../cert-out/dev
```



