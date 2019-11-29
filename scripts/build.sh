#!/bin/bash

set -e

usage(){
        cat <<EOF
        Usage: $(basename $0) <COMMAND>  <INDIR_PATH> <TAG>
        Commands:
            ci      run build process with new version and properly tag
            cd      deploy app to container registry.
            cid     ci+cd

        Command arguments:
            <INDIR_PATH>    required	directory where dockerfile is placed

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


INDIR=${2%/}
VERSION=0.1
tag=${3} #"$(date '+%Y%m%d%H%M%S')"

case "$1" in
    "ci")
        dir_exists ${INDIR}

        docker build ${INDIR} -t acronhosbx.azurecr.io/frontend:${tag}

        #remove all layers so docker will be build again
        docker rmi $(docker images | grep "^<none>" | awk '{ print $3 }')
    ;;
    "cd")
        # to be able to push into remote repo we need properly tag. I'm doing this step in ci part
        #docker tag onho.io/frontend:${tag} acronhosbx.azurecr.io/frontend:${tag}

        docker push acronhosbx.azurecr.io/frontend:${tag}
    ;;
    "cid")
        dir_exists ${INDIR}
    ;;
      *)
  usage
  exit 0
  ;;
esac

