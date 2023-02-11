#!/usr/bin/env bash

set -euo pipefail

PROGNAME=$(basename "$0")

OUTPUT_DIR="output/"
BUILD_SERVER=0
BUILD_CLIENT=0


function usage() {
    cat << EOF
Relique build script

Usage: ${PROGNAME} [flags]
    -h --help: Displays this help
    -o --output-dir: Output directory for generated artefacts
    --server: Build relique-server components
    --client: Build relique-client components
EOF
}

function log() {
    datestring=`date +"%Y-%m-%d %H:%M:%S"`
    echo "${datestring} [${PROGNAME}] ${@}"
}

function log_exit() {
    local exit_code=$?
    echo
    log "Script exited with status code '${exit_code}'"
}

function check_args() {
    if [ -z $OUTPUT_DIR ]; then
        OUTPUT_DIR="output/"
    fi
}

function cmdline() {
    POSITIONAL_ARGS=()

    while [[ $# -gt 0 ]]; do
        case $1 in
            -o|--output-dir)
                OUTPUT_DIR="$2"
                shift # past argument
                shift # past value
                ;;

            --server)
                BUILD_SERVER=1
                shift # past argument
                ;;

            --client)
                BUILD_CLIENT=1
                shift # past argument
                ;;

            -h|--help)
                usage
                exit 0
                shift # past argument
                ;;

            -*|--*)
                echo "Unknown option '$1'"
                usage
                exit 1
                ;;

            *)
                POSITIONAL_ARGS+=("$1") # save positional arg
                shift # past argument
                ;;
        esac
    done

    set -- "${POSITIONAL_ARGS[@]}" # restore positional parameters

    check_args
}

function build_binaries() {
    log "Building binaries to '$OUTPUT_DIR'"
    components=()

    if [ "X${BUILD_CLIENT}X" == "X1X" ]; then
        components+=("relique-client")
    fi

    if [ "X${BUILD_SERVER}X" == "X1X" ]; then
        components+=("relique-server")
    fi

    for component in "${components[@]}"; do
        log "Building $component"
        go env
        #go build -mod=vendor -v -o "${OUTPUT_DIR}/bin/${component}" cmd/${component}/main.go
        go build -v -o "${OUTPUT_DIR}/bin/${component}" cmd/${component}/main.go
        build_result=$?
        if [ "$build_result" -ne "0" ]; then
            log "Binary build failed !!! Aborting build script"
            exit $build_result
        fi
    done
}

function build_webui() {
    log "Building web UI to '$OUTPUT_DIR'"

    if [ "X${BUILD_SERVER}X" != "X1X" ]; then
        log "Not building relique-server. Skipping web UI build"
        return
    fi

    pushd ui
        if [ ! -f "./ui/node_modules" ]; then
            npm install
        fi
        npm run build
    popd
}


function copy_generic_bin_script() {
    log "Copying generic relique script"
    cp build/scripts/relique "${OUTPUT_DIR}/bin/relique"
}


function copy_service_files() {
    log "Copying systemd service files to '$OUTPUT_DIR'"
    mkdir -p "${OUTPUT_DIR}/usr/lib/systemd/system"
    cp -r build/init/*.service "${OUTPUT_DIR}/usr/lib/systemd/system"

    log "Copying freebsd init files to '$OUTPUT_DIR'"
    mkdir -p "${OUTPUT_DIR}/etc/rc.d"
    cp -r build/init/relique-client.freebsd.sh "${OUTPUT_DIR}/etc/rc.d/relique-client"
    cp -r build/init/relique-server.freebsd.sh "${OUTPUT_DIR}/etc/rc.d/relique-server"
}


function copy_config_defaults() {
    log "Copying default configuration files to '$OUTPUT_DIR'"
    cp -r configs/* "$OUTPUT_DIR"
}

function copy_webui_files() {
    log "Copying web UI built files to '$OUTPUT_DIR'"

    if [ "X${BUILD_SERVER}X" != "X1X" ]; then
        log "Not building relique-server. Skipping web UI copy"
        return
    fi

    mkdir -p "${OUTPUT_DIR}/var/lib/relique/ui"
    cp -r ui/build/* "${OUTPUT_DIR}/var/lib/relique/ui/"
}

function package_default_modules() {
    log "Packaging default modules tarballs to '$OUTPUT_DIR'"
    for mod in $(find "${OUTPUT_DIR}/var/lib/relique/default_modules" -mindepth 1 -maxdepth 1 -type d); do
        modname=$(basename $mod)
        log "Packaging default module '${modname}'"
        pushd "${OUTPUT_DIR}/var/lib/relique/default_modules/${modname}" > /dev/null
        tar --exclude-vcs -zcf ../${modname}.tar.gz .
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

function main() {
    trap log_exit EXIT
    cmdline "${@}"

    build_binaries
    build_webui
    copy_generic_bin_script
    copy_service_files
    copy_config_defaults
    copy_webui_files
    package_default_modules
    make_certs
}

main "${@}"
