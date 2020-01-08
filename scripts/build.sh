#!/bin/bash

set -e
# i.e.
# PS /home/michal/workspace/onho.io> ./scripts/build.sh ci ./infrastructure/docker/dev/ "latest"
usage(){
        cat <<EOF
        Usage: $(basename $0) <COMMAND>  <TAG>
        Commands:
            ci                run build process with new version and properly tag
            cd                deploy app to container registry.
            cid               ci+cd
            init              you should run this before first deployment starts. Instals istio and certificates renew certificates. not implemented yet...
            drop              drops onho namespace

        Command arguments:
            <TAG> :   docker tag, if empty than :latest is used
EOF
}



panic() {
  (>&2 echo "$@")
  exit 1
}


dir_exists(){
	local path="$1"
    	if [[ ! -d "$path" ]]; then
  		panic "$path doesn't exists"
     	fi
}

check_kube_cli(){
	KUBECTL=`which kubectl`||true

	if [[ -z "${KUBECTL}" ]]; then
 		panic "Kubectl is not installed"
		exit 1
	fi
}



ci(){
cat <<EOF
***************************************************************
    building docker image
***************************************************************
EOF
   dir_exists ${DOCKER_DIR}

    docker build ${DOCKER_DIR} -t acronhosbx.azurecr.io/onho:${tag}

    #remove all layers so docker will be build again
    docker rmi "$(docker images | grep "^<none>" | awk '{ print $3 }')"
}

cd(){
cat <<EOF
***************************************************************
    deploying docker image to container repo
***************************************************************
EOF

    check_kube_cli
    # to be able to push into remote repo we need properly tag. I'm doing this step in ci part
    #docker tag onho.io/onho:${tag} acronhosbx.azurecr.io/onho:${tag}

    docker push acronhosbx.azurecr.io/onho:${tag}

    #TODO: replace to polly manifest

    kubectl apply -f ${KUBE_DIR}namespace.yaml
    kubectl apply -f ${KUBE_DIR}config.yaml
    local secret="${KUBE_DIR}secrets.yaml"
    isEncrypted="$(sed -n '1{/^$ANSIBLE_VAULT;1.1;AES256/p};q' ${secret})"
    if [[ ! "$isEncrypted" ]]; then
      kubectl apply -f ${KUBE_DIR}secrets.yaml
    fi


    #frontend
    kubectl apply -f ${KUBE_DIR}frontend.pod.yaml
    kubectl apply -f ${KUBE_DIR}frontend.service.yaml
    kubectl apply -f ${KUBE_DIR}frontend.gw.yaml

    #backend
    kubectl apply -f ${KUBE_DIR}backend.yaml


    #rabbit-mq
    kubectl apply -f ${KUBE_DIR}rabbit-mq.yaml
}


shallow_drop(){
cat <<EOF
***************************************************************
    cleaning k8s and dockers
***************************************************************
EOF

    kubectl delete -f ${KUBE_DIR}frontend.pod.yaml
    kubectl delete -f ${KUBE_DIR}backend.yaml
    #waiting until pods are initialised
    sleep 15s
}



drop(){
cat <<EOF
***************************************************************
    cleaning k8s and dockers
***************************************************************
EOF

    #TODO: replace to polly manifest

    kubectl delete -f ${KUBE_DIR}namespace.yaml
}

init(){
  check_kube_cli

  #create certificates here
  kubectl delete secret istio-ingressgateway-certs -n istio-system
  sleep 2s
  kubectl create secret tls istio-ingressgateway-certs -n istio-system --key onho.cz.key --cert onho.cz.crt
  echo "uploading cert..."
  sleep 70s
  kubectl exec -it -n istio-system "$(kubectl -n istio-system get pods -l istio=ingressgateway -o jsonpath='{.items[0].metadata.name}')" -- ls -al /etc/istio/ingressgateway-certs
  #remove crtificates here

}


SCRIPTDIR=$( pwd -P )
ROOTDIR="${SCRIPTDIR}/../"
DOCKER_DIR="${ROOTDIR%/}/infrastructure/docker/dev/"
KUBE_DIR="${ROOTDIR%/}/infrastructure/k8s/dev/"



tag=${2} #"$(date '+%Y%m%d%H%M%S')"

case "$1" in
    "ci")
       ci
    ;;
    "cd")
      cd
    ;;
    "cid")
        ci
        cd
    ;;
    "init")
        init
    ;;
    "drop")
        drop
    ;;
      *)
  usage
  exit 0
  ;;
esac

