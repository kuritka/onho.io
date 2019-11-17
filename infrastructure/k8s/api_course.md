

connecting to minikubes
```bash
kubectl config get-contexts

kubectl config use-context minikube

kubectl cluster-info
```

api-resources documentation
```bash
kubectl api-resources | more

kubectl api-resources --api-group=apps

#contains links to specs, version, meatadata...
kubectl explain pods | more

#more detailed information about spec part
kubectl explain pods.spec | more

#even more detailed info
kubectl explain pods.spec.containers | more


kubectl explain deployment 
#deprecated by apps/v1beta2

kubectl explain deployment --api-version apps/v1beta2
#deprecated by v1 version

kubectl explain deployment --api-version apps/v1
#ok version

kubectl api-versions  |sort |more
#gets all api versions
```


create pod with YAML
```bash
kubectl apply -f pod.yaml

#getting deployment from api-server. Depends on k8s version but maight happen 
# that for some resources you create resource under v1 but it exists under version 
# v1betav2. i.e. deployment in k8s 1.13. Thats why there is further compatiboility
kubectl get deployment hello-world-pod-me -o yaml


```

info about POD request
```bash
kubectl get pod hello-world-pod-me-57f8fd8c76-5g9l4 -v 7
#returns url path to API server and response

#we cannot use curl directly because of authentication
kubectl proxy & 
curl <curl from previous command> | head -n 20
#kubectl proxy uses local kubeconfig to authentiate me agains API server


#I think it will return the sam result as following command : 
kubectl get pod hello-world-pod-me-57f8fd8c76-5g9l4 -o json
#returns detailed information about POD

kubectl get pod hello-world-pod-me-57f8fd8c76-5g9l4 
#returns very basic info about POD


kubectl get pods --watch -v6
#retrieving watch list with verbose info
```



netstat
```bash
#If I run in one terminal window 
kubectl get pods --watch -v6

#and open second terminal window 
netstat -plant | grep kubectl
#I get established connection with minikubes
tcp        0      0 127.0.0.1:8001          0.0.0.0:*               LISTEN      9510/kubectl        
tcp        0      0 192.168.39.1:37016      192.168.39.73:8443      ESTABLISHED 12532/kubectl       

```



logs
```bash

kubectl logs  hello-world-pod-me-57f8fd8c76-5g9l4 
kubectl logs  hello-world-pod-me-57f8fd8c76-5g9l4 -v 6 
#this makes two requests - one to get resource and if exists get logs..


[michal@dhcp-10-20-17-69 k8s]$ kubectl logs  hello-world-pod-me-57f8fd8c76-5g9l4 -v 6
I1106 21:18:45.701973   15266 loader.go:375] Config loaded from file:  /home/michal/.kube/config
I1106 21:18:47.668067   15266 round_trippers.go:443] GET https://192.168.39.73:8443/api/v1/namespaces/default/pods/hello-world-pod-me-57f8fd8c76-5g9l4 200 OK in 1959 milliseconds
I1106 21:18:53.515776   15266 round_trippers.go:443] GET https://192.168.39.73:8443/api/v1/namespaces/default/pods/hello-world-pod-me-57f8fd8c76-5g9l4/log 200 OK in 5839 milliseconds
Running on http://localhost:8080

kubectl  proxy & 
curl <address from previous command>

```



## config

all client certificates, keys and user names are here
`~/.kube/config`

```bash
cp ~/.kube/config ~/.kube/config.ORIG
vim ~/.kube/config  
#and edit user name 
kubectl get pods -v6
#gets an error that cannot authenticate
```


#http statuses
```bash
kubectl apply -f ./hello.yaml -v 6 
#HTTP POST which creates resource

kubectl delete deployment hello-world -v 6 
#HTTP DELETE which deletes resource
```

#namespaces

```bash
kubectl api-resources 
kubectl api-resources --namespaced=false
#resource that cannot be within the namespace
kubectl api-resources --namespaced=false
#resource that can be within the namespace

kubectl describe namespaces
#check if namespace is Active state or Terminated
kubectl describe namespaces kube-system

kubectl get pods --all-namespaces

kubectl get all --all-namespaces

kubectl get pods --namespace kube-system

kubectl get pods --all  --namespace=playground1
kubectl delete pods --all  --namespace=playground1
```


#declaring pod imperatively
```bash
kubectl run hello-world-pod \
    --image=grc.io/google-samples/hello-app:1.0 \
    --generator=run-pod/v1 \
    --namespace playground1
    

```




## DEMO

```bash
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod-1
  labels:
    app: MyWebApp
    deployment: v1
    tier: prod
spec:
  containers:
  - name: nginx
    image: nginx
    ports:
    - containerPort: 80
---
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod-2
  labels:
    app: MyWebApp
    deployment: v1.1
    tier: prod
spec:
  containers:
    - name: nginx
      image: nginx
      ports:
        - containerPort: 80
---
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod-3
  labels:
    app: MyWebApp
    deployment: v1.1
    tier: qa
spec:
  containers:
    - name: nginx
      image: nginx
      ports:
        - containerPort: 80
---
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod-4
  labels:
    app: MyAdminApp
    deployment: v1
    tier: prod
spec:
  containers:
    - name: nginx
      image: nginx
      ports:
        - containerPort: 80

```

and than 

```bash
kubectl get pods --selector tier=prod

NAME          READY   STATUS    RESTARTS   AGE
nginx-pod-1   1/1     Running   0          9m40s
nginx-pod-2   1/1     Running   0          9m40s
nginx-pod-4   1/1     Running   0          9m39s


kubectl get pods --show-labels -l 'tier=prod'

NAME          READY   STATUS              RESTARTS   AGE    LABELS
nginx-pod-1   0/1     ContainerCreating   0          106s   app=MyWebApp,deployment=v1,tier=prod
nginx-pod-2   0/1     ContainerCreating   0          106s   app=MyWebApp,deployment=v1.1,tier=prod
nginx-pod-4   0/1     ContainerCreating   0          105s   app=MyAdminApp,deployment=v1,tier=prod


kubectl get pods -l 'tier=prod,app=MyWebApp' --show-labels
NAME          READY   STATUS    RESTARTS   AGE   LABELS
nginx-pod-1   1/1     Running   0          12m   app=MyWebApp,deployment=v1,tier=prod
nginx-pod-2   1/1     Running   0          12m   app=MyWebApp,deployment=v1.1,tier=prod



kubectl get pods -L tier,app
NAME                                  READY   STATUS    RESTARTS   AGE     TIER   APP
hello-world-pod-me-57f8fd8c76-5g9l4   1/1     Running   5          4d18h          hello-world-pod-me
hello-world-pod-me-57f8fd8c76-f6vnw   1/1     Running   5          4d18h          hello-world-pod-me
hello-world-pod-me-57f8fd8c76-fdp4c   1/1     Running   5          4d18h          hello-world-pod-me
hello-world-pod-me-57f8fd8c76-jcvfs   1/1     Running   5          4d17h          hello-world-pod-me
hello-world-pod-me-57f8fd8c76-xnn98   1/1     Running   5          4d18h          hello-world-pod-me
nginx-pod-1                           1/1     Running   0          14m     prod   MyWebApp
nginx-pod-2                           1/1     Running   0          14m     prod   MyWebApp
nginx-pod-3                           1/1     Running   0          13m     qa     MyWebApp
nginx-pod-4                           1/1     Running   0          13m     prod   MyAdminApp

```


## add and remove label
```bash
kubectl label pod nginx-pod-1 another=Label
kubectl get pod nginx-pod-1 --show-labels
kubectl label pod nginx-pod-1 another-
kubectl get pod nginx-pod-1 --show-labels

#change all tier labels to non-prod value
kubectl label pod --all tier=non-prod --overwrite
kubectl delete pod -l tier=non-prod    
```



##service
```bash
#describes service
kubectl describe servicehello-world
#describes pods attached to service
kubectl describe endpoints hello-world
```



# DEMO

to see realtime changes run this command first:
```bash
kubectl get events --watch &
```

after you want to put events back on the foreground run this:
```bash
kubectl get events --watch
```

you can also see all events on POD in describe (last part)
```bash
kubectl describe pod hello-world-pod-me-5f69fbf977-8szwm
```

install deployment with replica set
```bash
kubectl apply -f ./deployment-me.yaml
#scaled rs
#create rs
#create pod
#assign pod
#pulled container
#create pod
#assign
#started pod
```

Change replica set size and see what happens 
```bash
kubectl scale deployment hello-world-pod-me --replicas=2
#killing pod
#deleting pod

kubectl scale deployment hello-world-pod-me --replicas=3
#scaled rs
#create pod
#assign pod
#pull - image already exist 
#create container
#start
```

Remove pod from replica set by changing label
```bash
kubectl label pod hello-world-pod-me-5f69fbf977-2c5dd app=DEBUG --overwrite
kubectl get pods --show-labels
```
Go to container
```bash
kubectl -v 6 exec -it hello-world-pod-me-5f69fbf977-2c5dd -- /bin/sh
hostname 
#hostname is pod name, musi tam byt mezera mezi -- bin/sh
ps 
exit

#3124 loader.go:375] Config loaded from file:  /home/michal/.kube/config
#I1110 13:00:12.043895   13124 round_trippers.go:443] GET https://aks-onho-we-dns-292362c2.hcp.westeurope.azmk8s.io:443/api/v1/namespaces/default/pods/hello-world-pod-me-5f69fbf977-2c5dd 200 OK in 267 milliseconds
#I1110 13:00:12.369634   13124 round_trippers.go:443] POST https://aks-onho-we-dns-292362c2.hcp.westeurope.azmk8s.io:443/api/v1/namespaces/default/pods/hello-world-pod-me-5f69fbf977-2c5dd/exec?command=%2Fbin%2Fsh&container=node-hello&stdin=true&stdout=true&tty=true 101 Switching Protocols in 317 milliseconds



#in case of multi container pod I must specify container as well
kubectl -v 6 exec -it hello-world-pod-me-5f69fbf977-2c5dd --container CONTAINER -- /bin/sh
#the same for logs

```

We can run something else than shell script - we can run any script. i.e. kill app
```bash
kubectl -v 6 exec -it hello-world-pod-me-5f69fbf977-2c5dd -- /usr/bin/killall myapp
```

port forward

```bash
kubectl port-forward hello-world-pod-me-5f69fbf977-2c5dd 80:8080
#localport:containerport
#bind: permission denied unable to create listener 
# because all ports till 1024 are privileged on linux and requires sudo
# instead 8080 as local port must be used 
#THIS WILL BE WORKING ONLY LOCALLY, IF YOU ARE WITHIN NODE. 
# AKS doesnt allow you to directly connect to NODE
# TO ACCESS AKS you need create service
```

