#!/bin/bash

set -e
# i.e.
# PS /home/michal/workspace/onho.io> ./scripts/build.sh ci ./infrastructure/docker/dev/ "latest"
usage(){
        cat <<EOF
        Usage: $(basename $0) <COMMAND>  <ONHO_DIR_PATH> <TAG>
        Commands:
            ci                run build process with new version and properly tag
            cd                deploy app to container registry.
            cid               ci+cd
            init              you should run this before first deployment starts. Instals istio and certificates renew certificates. not implemented yet...
            drop

        Command arguments:
            <ONHO_DIR_PATH>    root directory of onho

         tag :
            if empty than :latest is used
EOF
}



panic() {
  (>&2 echo "$@")
  exit -1
}


dir_exists(){
	local path="$1"
    	if [[ ! -d "$path" ]]; then
  		panic "$path doesn't exists"
     	fi
}


build(){
    local crt="$1"
}




ci(){
cat <<EOF
***************************************************************
    building docker image
***************************************************************
EOF
   dir_exists ${DOCKER_DIR}

    docker build ${DOCKER_DIR} -t acronhosbx.azurecr.io/frontend:${tag}

    #remove all layers so docker will be build again
    docker rmi $(docker images | grep "^<none>" | awk '{ print $3 }')
}

cd(){
cat <<EOF
***************************************************************
    deploying docker image to container repo
***************************************************************
EOF

    # to be able to push into remote repo we need properly tag. I'm doing this step in ci part
    #docker tag onho.io/frontend:${tag} acronhosbx.azurecr.io/frontend:${tag}

    docker push acronhosbx.azurecr.io/frontend:${tag}

    #TODO: replace to polly manifest

    #frontend
    local secret="${KUBE_DIR}secrets.yaml"
    kubectl apply -f ${KUBE_DIR}namespace.yaml
    kubectl apply -f ${KUBE_DIR}config.yaml
    isEncrypted="$(sed -n '1{/^$ANSIBLE_VAULT;1.1;AES256/p};q' ${secret})"
    if [[ ! "$isEncrypted" ]]; then
      kubectl apply -f ${KUBE_DIR}secrets.yaml
    fi
    kubectl apply -f ${KUBE_DIR}frontend.pod.yaml
    kubectl apply -f ${KUBE_DIR}frontend.service.yaml
    kubectl apply -f ${KUBE_DIR}frontend.gw.yaml


    #rabbit-mq
    kubectl apply -f ${KUBE_DIR}rabbit-mq.yaml
}


drop(){
cat <<EOF
***************************************************************
    cleaning k8s and dockers
***************************************************************
EOF

    #TODO: replace to polly manifest

    #frontend
    kubectl delete -f ${KUBE_DIR}namespace.yaml
}


DOCKER_DIR="${2%/}/infrastructure/docker/dev/"
KUBE_DIR="${2%/}/infrastructure/k8s/dev/"

VERSION=0.1
tag=${3} #"$(date '+%Y%m%d%H%M%S')"

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
    "drop")
        drop
    ;;
      *)
  usage
  exit 0
  ;;
esac

