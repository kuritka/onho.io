#!/bin/bash

set -e

# Helper script to setup EDT

logError() {
  (>&2 echo "$@")
}

usage(){
	cat <<EOF
	Usage: $(basename $0) <COMMAND> <ARGUMENTS>

	Commands:
        drop                        removes all containers from local docker
        install                     installing postgre, pgadmin and mq
        start                       starts onho with mock services

	Command arguments:
        drop
        install
        start                       [sensors] [controllers] [param3]
EOF
}


DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT="$DIR/.."


#drops all containers from docker deamon
drop_all_containers(){
	docker container prune
    #    docker image prune
    #    docker network prune
    #    docker volume prune
    docker rmi $(docker images | grep "^<none>" | awk '{ print $3 }')
}


install(){
    docker run -d --hostname my-rabbit --name some-rabbit \
    -e RABBITMQ_DEFAULT_USER=guest \
    -e RABBITMQ_DEFAULT_PASS=guest \
    -p 15672:15672 -p 5671:5671 \
    -p 5672:5672 -p 15671:15671 \
    -p 25672:25672 rabbitmq:3-management

    printf "rabbitmq on localhost:15671 guest/guest \n\n"

    docker run --rm --name pg-docker \
     -e POSTGRES_PASSWORD=password \
     -e POSTGRES_USER=test \
     -e POSTGRES_DB=test_db \
     -p 5432:5432 \
     -v $HOME/docker/volumes/postgres:/var/lib/postgresql/data \
     -d postgres:latest

     printf "postgre on localhost:5432 test@test_db/password  \n\n"

     docker run -p 8085:80 \
     --link pg-docker:pg-docker \
     -e "PGADMIN_DEFAULT_EMAIL=admin@local.com" \
     -e "PGADMIN_DEFAULT_PASSWORD=password" \
     -d dpage/pgadmin4 \
     --name pgadmin

     printf "pgadmin on localhost:8085 admin@local.com/password \n\n\n"
}

start(){
    printf ${1}+${2}+${3}"HELLO\n\n\n"
}

case "$1" in
    "drop")
            printf "dropping all container...\n\n\n"
            drop_all_containers
            ;;
    "install")
                printf "installing local infrastructure...\n\n\n"
                install
                ;;
    "start")

                if [[ "$#" -lt 4 ]]; then
                        usage
                        exit 1
                fi
                shift

                PARAM_1=$1
                PARAM_2=$2
                PARAM_3=$3

                printf "starting application...\n\n\n"
                start ${PARAM_1} ${PARAM_2} ${PARAM_3}
                ;;
             *)
     usage
     exit 0
     ;;
esac