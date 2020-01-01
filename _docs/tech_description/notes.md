# k8s useful commands

Getting service external IP:
```
kubectl get service/frontend -n onho-dev -o jsonpath='{.status.loadBalancer.ingress[0].ip}'
```
