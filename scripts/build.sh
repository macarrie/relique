#!/usr/bin/env bash


function usage() {
    echo "\
        usage: $0 [options]

    Options:
    -h --help: Displays this help
    -o --output-dir: Output directory for generated artefacts
    "
}


function build_binaries() {
    echo "Building binaries to '$OUTPUT_DIR'"
    components="relique-client relique-server"

    for component in $components; do
        echo "Building $component"
        go build -o "${OUTPUT_DIR}/usr/bin/${component}" cmd/${component}/main.go
    done
}


function copy_service_files() {
    echo "Copying systemd service files to '$OUTPUT_DIR'"
    mkdir -p "${OUTPUT_DIR}/usr/lib/systemd/system"
    cp -r build/init/*.service "${OUTPUT_DIR}/usr/lib/systemd/system"
}


function copy_config_defaults() {
    echo "Copying default configuration files to '$OUTPUT_DIR'"
    cp -r configs/* "$OUTPUT_DIR"
}

function package_default_modules() {
    # TODO: Remove .git folders if found with tar --exclude-vcs --exclude-vcs-ignore
    echo "Packaging default modules tarballs to '$OUTPUT_DIR'"
    for mod in $(ls -1 ${OUTPUT_DIR}/var/lib/relique/default_modules); do
        pushd "${OUTPUT_DIR}/var/lib/relique/default_modules/${mod}" > /dev/null
            tar -zcf ../${mod}.tar.gz .
        popd > /dev/null
    done
}


# Create self signed certs for quick first setup
function make_certs() {
    mkdir -p "${OUTPUT_DIR}/etc/relique/certs"
    echo  -e "[req]\ndistinguished_name=req\n[san]\nsubjectAltName=DNS.1:localhost,DNS.2:relique" > "${OUTPUT_DIR}/tmp.certs"
    openssl req \
        -x509 \
        -newkey rsa:4096 \
        -sha256 \
        -days 3650 \
        -nodes \
        -keyout "${OUTPUT_DIR}/etc/relique/certs/key.pem" \
        -out "${OUTPUT_DIR}/etc/relique/certs/cert.pem" \
        -subj '/CN=relique' \
        -extensions san \
        -config "${OUTPUT_DIR}/tmp.certs"
    rm "${OUTPUT_DIR}/tmp.certs"
}


POSITIONAL=()
while [[ $# -gt 0 ]]
do
    key="$1"

    case $key in
        -o|--output-dir)
            OUTPUT_DIR="$2"
            shift # past argument
            shift # past value
            ;;

        -h|--help)
            usage
            exit 0
            shift # past argument
            ;;

        *)    # unknown option
            POSITIONAL+=("$1") # save it in an array for later
            shift # past argument
            ;;
    esac
done
set -- "${POSITIONAL[@]}" # restore positional parameters


if [ -z $OUTPUT_DIR ]; then
    OUTPUT_DIR="output/"
fi

build_binaries
copy_service_files
copy_config_defaults
package_default_modules
make_certs
