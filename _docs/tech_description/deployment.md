## Deployment


### initial deployment
- Create self-sign root certificate anf key
    ```bash
    $ ./scripts/cert-generator.sh self-sign  ./infrastructure/certificates/dev/onho.cnf ../onho-certs/dev/
    ```

- Deploy TLS key/cert as secret for the ingress gateway:

	```bash
	pushd ../onho-certs/dev/
	kubectl create -n istio-system secret tls istio-ingressgateway-certs --key onho.key --cert onho.crt
	popd
	```
	
	wait for a while and check that certificates are there. 
	```bash
  $ kubectl exec -it -n istio-system $(kubectl -n istio-system get pods -l istio=ingressgateway -o jsonpath='{.items[0].metadata.name}') -- ls -al /etc/istio/ingressgateway-certs
   
  drwxr-xr-x 2 root root   80 Jan  4 12:32 ..2020_01_04_12_32_52.647689794
  lrwxrwxrwx 1 root root   31 Jan  4 12:32 ..data -> ..2020_01_04_12_32_52.647689794
  lrwxrwxrwx 1 root root   14 Jan  4 12:32 tls.crt -> ..data/tls.crt
  lrwxrwxrwx 1 root root   14 Jan  4 12:32 tls.key -> ..data/tls.key
    ```


To deploy app run following command to redeploy
```bash
$ ./scripts/build.sh cid ./infrastructure/docker/dev/ 0.81  
```


### update
```bash
$ ./scripts/build.sh cid ./infrastructure/docker/dev/ 0.81  
```

