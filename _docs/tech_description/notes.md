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

Istio ingress gateway
```bash
kubectl -n istio-system get service istio-ingressgateway
```