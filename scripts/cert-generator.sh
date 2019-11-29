#!/bin/bash

set -e

usage(){
        cat <<EOF
        Usage: $(basename $0) <COMMAND> <ARGUMENTS> <INDIR_PATH> <OUTDIR_PATH>
        Commands:
            csrgen  create new certificate requests from configurations
            pfx  create new pfx certificates from .crt in PEM or DR3 format. pfx is created in folder where certificate exists

        Command arguments:
            csrgen
                <INDIR_PATH>    required	certificate configuration path
                <OUTDIR_PATH>   required    certificate OUTDIR path, create if not exists
             pfx
                <CRT_FILEPATH>  required    .crt filepath.
                <PEM_FILEPATH>  required    .pem filepath.
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


file_exists(){
	local path="$1"
    	if [[ ! -f "$path" ]]; then
  		panic "$path doesn't exists"
     	fi
}

create_dir_if_not_exists(){
  local path="$1"
  [[ -d "$path" ]] || mkdir -p ${path}
}


generate_requests(){
    local indir="$1"
    local outdir="$2"
	for file in ${indir}/*.cnf; do
           echo "$(basename "$file")"
           filename="$(basename "$file")"
           name="${filename%.*}"
           pem="$outdir/$name.pem"
           csr="$outdir/$name.csr"
           cnf="$indir/$name.cnf"

           openssl genrsa -out "${pem}" 2048
           openssl req -new -key "${pem}" -out "${csr}" -extensions v3_req -config "${cnf}" -verbose
           cp -rf "${cnf}" "$outdir/$name.cnf"
	   done
}


create_pfx(){
    local crt="$1"
    local pem="$2"
    baseName="$(basename "$crt")"
    generalName="${baseName%.*}"
    dirName="$(dirname "$crt")"
    pfx="${dirName}/${generalName}"
    isPem="$(sed -n '1{/^-----BEGIN CERTIFICATE-----/p};q' ${crt})"
    if [[ ! "$isPem" ]]; then
        openssl x509 -outform pem -inform der -in "${crt}" -out "${pfx}.pem.crt"
        openssl pkcs12 -inkey ${pem} -in "${pfx}.pem.crt" -export -out "${pfx}.pfx"  -passout pass:
        exit 0
    fi
    openssl pkcs12 -inkey ${pem} -in ${crt} -export -out "${pfx}.pfx"  -passout pass:
}


case "$1" in
    "csrgen")
            if [[ "$#" -ne 3 ]]; then
                usage
                exit 0
            fi
            INDIR=${2%/}
            OUTDIR=${3%/}
            dir_exists ${INDIR}
            create_dir_if_not_exists ${OUTDIR}
            generate_requests ${INDIR} ${OUTDIR}
            ;;
    "pfx")
            if [[ "$#" -ne 3 ]]; then
                usage
                exit 0
            fi
            CERT=${2%}
            PEM=${3%}
            file_exists ${CERT}
            file_exists ${PEM}
            create_pfx ${CERT} ${PEM}
            ;;
    *)
          usage
          exit 0
          ;;
esac



