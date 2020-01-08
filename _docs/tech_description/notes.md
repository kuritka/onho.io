# k8s useful commands

Getting service external IP:
```
kubectl get service/frontend -n onho-dev -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
```


getting events for pod (probes must exists)
```bash
kubectl get event -n onho-dev --field-selector involvedObject.name=frontend-57dbc9d845-tv4kg
```



#istio

label istio namespace from commandline. After this value is set you can use istio within namespace
```
kubectl label namespace default istio-injection=enabled
```

## ingress 
https://istio.io/docs/tasks/traffic-management/ingress/ingress-certmgr/
https://istio.io/docs/reference/config/networking/gateway/


Istio ingress gateway (https://istio.io/docs/tasks/traffic-management/ingress/ingress-control/). Don't forget to put resources into correct namespace! 
```bash
kubectl -n istio-system get service istio-ingressgateway
```

Verifying creds were mounted to istio-ingress-gateway pods
```bash
kubectl exec -it -n istio-system $(kubectl -n istio-system get pods -l istio=ingressgateway -o jsonpath='{.items[0].metadata.name}') -- ls -al /etc/istio/ingressgateway-certs
```
