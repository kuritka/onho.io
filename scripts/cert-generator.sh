#!/bin/bash

set -e

usage(){
        cat <<EOF
        Usage: $(basename $0) <COMMAND> <ARGUMENTS> <INDIR_PATH> <OUTDIR_PATH>
        Commands:
            csrgen  create new certificate requests from configurations
            pfx  create new pfx certificates from .crt in PEM or DR3 format. pfx is created in folder where certificate exists
            self-sign create new self-sign certificate from configuration
            keystore  creates new store for certificates. not implemented yet

        Command arguments:
            csrgen
                <INDIR_PATH>    required	certificate configuration path
                <OUTDIR_PATH>   required    certificate OUTDIR path, create if not exists

             pfx
                <CRT_FILEPATH>  required    .crt filepath.
                <PEM_FILEPATH>  required    .pem filepath.

            self-sign
                <INDIR_PATH>    required	certificate configuration path
                <OUTDIR_PATH>   required    certificate OUTDIR path, create if not exists
                <keystore>      optional    not implemented yet

             keystore
                <NAME>          required    name of keystore
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

# https://stackoverflow.com/questions/10175812/how-to-create-a-self-signed-certificate-with-openssl
create_self_sign(){
    local cnf="$1"
    local outdir="$2"
    baseName="$(basename "$cnf")"
    name="${baseName%.*}"
    pem="$outdir/$name.key"
    csr="$outdir/$name.csr"
    crt="$outdir/$name.crt"
    #signing key on the side of CA. for simplicity I use the same key for csr as for signing
    signkey="${pem}"
    openssl genrsa -out "${pem}" 2048
    #openssl req -sha256 -nodes -newkey rsa:2048 -keyout "${pem}"
    openssl req -new -key "${pem}" -out "${csr}" -extensions v3_req -config "${cnf}" -verbose
    openssl x509 -in "${csr}" -out "${crt}" -req -signkey "${signkey}" -days 365

    #verification
    #openssl verify -CAfile RootCert.pem -untrusted Intermediate.pem UserCert.pem
    openssl x509 -in  "${crt}" -text -noout

    ##cat "${crt}"
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
    "self-sign")
           if [[ "$#" -ne 3 ]]; then
                usage
                exit 0
            fi
            CNF=${2%}
            OUTDIR=${3%/}
            file_exists ${CNF}
            create_dir_if_not_exists ${OUTDIR}
            create_self_sign ${CNF} ${OUTDIR}
            ;;
    *)
          usage
          exit 0
          ;;
esac



