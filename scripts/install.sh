#!/usr/bin/env bash

set -euo pipefail

PROGNAME=$(basename "$0")

USER=relique
GROUP=relique

ABS_PATH=$(readlink -f "$0")
BASE=$(dirname "${ABS_PATH}")

ROOT_BIN_PATH="/usr/bin"
ROOT_CFG_PATH="/etc/relique"
ROOT_DATA_PATH="/var/lib/relique"

SYSTEMD=0
FREEBSD=0
SKIPMODULEINSTALL=0
SKIPUSERCREATION=0

INSTALL_SERVER=0
INSTALL_CLIENT=0

function usage() {
cat << EOF
Relique install script

Usage: ${PROGNAME} [flags]
    -h --help: Displays this help
    -p --prefix: Install relique to this folder
    -s --src: Get compiled relique package to install from this folder
    --server: Install relique server
    --client: Install relique client
    --systemd: Install systemd service file
    --freebsd: Install freebsd service file
    --skip-user-creation: Skip relique group and user creation
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
    if [ -z $PREFIX ]; then
        log "Missing install prefix"
        usage
        exit 1
    fi

    if [ -z $SRC ]; then
        log "Missing install source directory"
        usage
        exit 1
    fi

    if [ ! -d $SRC ]; then
        log "Source directory '$SRC' does not exist"
        exit 1
    fi

    if [ ! -d $PREFIX ]; then
        mkdir -p "${PREFIX}"
    fi

    if [ "$INSTALL_SERVER" != "1" ] && [ "$INSTALL_CLIENT" != "1" ]; then
        log "Installing neither client or server. Please select at least one option with --client or --server"
        usage
        exit 1
    fi
}

function cmdline() {
    POSITIONAL_ARGS=()

    while [[ $# -gt 0 ]]; do
        case $1 in
            -p|--prefix)
            PREFIX="$2"
            shift # past argument
            shift # past value
            ;;

            -s|--src)
            SRC="$2"
            shift # past argument
            shift # past value
            ;;

            --server)
            INSTALL_SERVER=1
            shift # past argument
            ;;

            --client)
            INSTALL_CLIENT=1
            shift # past argument
            ;;

            --systemd)
            SYSTEMD=1
            shift # past argument
            ;;

            --freebsd)
            FREEBSD=1
            ROOT_BIN_PATH="/usr/local/bin"
            ROOT_CFG_PATH="/usr/local/etc/relique"
            ROOT_DATA_PATH="/usr/local/relique"
            shift # past argument
            ;;

            --skip-user-creation)
            SKIPUSERCREATION=1
            shift # past argument
            ;;

            --skip-module-install)
            SKIPMODULEINSTALL=1
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

function install_file() {
    local src_file=$1
    local dest_file=${2-$1}
    local overwrite=${3-0}

    if [ "X${dest_file}X" == "XX" ]; then
        dest_file="${src_file}"
    fi

    if [ -f "${SRC}/${src_file}" ]; then
        log "--- Found '${SRC}/${src_file}' file to copy"
    fi

    if [ -f "${PREFIX}/${dest_file}" ] && [ "X${overwrite}X" != "X1X" ]; then
        log "--- Skipping ${PREFIX}/${dest_file} copy. File already exists"
        return
    fi

    dest_path=$(dirname $dest_file)
    if [ ! -d "${PREFIX}/${dest_path}" ]; then
        log "--- Creating non-existing directory '${PREFIX}/${dest_path}' before copying file"
        mkdir -p "${PREFIX}/${dest_path}"
    fi

    log "--- Copying ${SRC}/${src_file} to ${PREFIX}/${dest_file}"
    cp "${SRC}/${src_file}" "${PREFIX}/${dest_file}"
}

function install_cfg_file() {
    local src_file=$1

    # Install cfg files into ROOT_CFG_PATH without overwriting
    install_file "etc/relique/${src_file}" "${ROOT_CFG_PATH}/${src_file}" 0
}

function install_binary() {
    local src_file=$1

    # Install cfg files into ROOT_CFG_PATH without overwriting
    install_file "bin/${src_file}" "${ROOT_BIN_PATH}/${src_file}" 1
}

function install_template() {
    local src_file=$1

    # Templating is only done for configuration files
    install_cfg_file "${src_file}"

    log "--- Templating ${PREFIX}/${ROOT_CFG_PATH}/${src_file}"
    sed -i"" -e "s#__CFG__#${ROOT_CFG_PATH}#"   "${PREFIX}/${ROOT_CFG_PATH}/${src_file}"
    sed -i"" -e "s#__DATA__#${ROOT_DATA_PATH}#"  "${PREFIX}/${ROOT_CFG_PATH}/${src_file}"
}

function copy_webui_files() {
    log "Installing web UI files"

    local webui_src="${SRC}/var/lib/relique/ui"

    if [ ! -d "${webui_src}" ]; then
        log "--- Missing web UI files. Skipping copy"
        return
    fi

    log "--- Copying ${webui_src} to ${PREFIX}/${ROOT_DATA_PATH}/ui"
    cp -r "${SRC}/var/lib/relique/ui" "${PREFIX}/${ROOT_DATA_PATH}/ui"
}

function copy_binaries() {
    log "Installing binaries"

    if [ "X${INSTALL_SERVER}X" == "X1X" ]; then
        install_binary "relique-server"
    fi

    if [ "X${INSTALL_CLIENT}X" == "X1X" ]; then
        install_binary "relique-client"
    fi

    # Install generic relique binary
    install_binary "relique"
}

function copy_default_configuration() {
    log "Installing default configuration"

    if [ "X${INSTALL_SERVER}X" == "X1X" ]; then
        install_template "server.toml.sample"
        install_cfg_file "schedules/daily.toml"
        install_cfg_file "schedules/weekly.toml"
        install_cfg_file "schedules/manual.toml"
        install_cfg_file "clients/example.toml.disabled"
        mkdir -p "${PREFIX}/${ROOT_DATA_PATH}"
        mkdir -p "${PREFIX}/${ROOT_DATA_PATH}/storage"
    fi

    if [ "X${INSTALL_CLIENT}X" == "X1X" ]; then
        install_template "client.toml.sample"
    fi
}


function copy_certs() {
    log "Installing self signed quick start certs"

    install_cfg_file "certs/cert.pem"
    install_cfg_file "certs/key.pem"
}


function create_user {
    log "Creating $USER user"

    id -u $USER > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        useradd -M "${USER}"
    fi
}


function create_dir_structure {
    log "Creating relique directory structure"

    mkdir -p "${PREFIX}/${ROOT_CFG_PATH}"
    mkdir -p "${PREFIX}/${ROOT_DATA_PATH}"
    mkdir -p "${PREFIX}/${ROOT_DATA_PATH}/modules"

    if [ "X${INSTALL_SERVER}X" == "X1X" ]; then
        mkdir -p "${PREFIX}/${ROOT_DATA_PATH}/db"
        mkdir -p "${PREFIX}/${ROOT_DATA_PATH}/storage"
    fi
}


function setup_files_ownership() {
    log "Setting files rights and ownership"

    chown -R $USER:$GROUP "${PREFIX}/${ROOT_CFG_PATH}"
    chown -R $USER:$GROUP "${PREFIX}/${ROOT_DATA_PATH}"
}


function install_systemd_service() {
    log "Installing systemd service files"

    if [ "X${INSTALL_SERVER}X" == "X1X" ]; then
        install_file "usr/lib/systemd/system/relique-server.service"
    fi

    if [ "X${INSTALL_CLIENT}X" == "X1X" ]; then
        install_file "usr/lib/systemd/system/relique-client.service"
    fi
}

function install_freebsd_service() {
    log "Installing freebsd service files"

    if [ "X${INSTALL_SERVER}X" == "X1X" ]; then
        install_file "etc/rc.d/relique-server" "usr/local/etc/rc.d/relique-server" 1
    fi

    if [ "X${INSTALL_CLIENT}X" == "X1X" ]; then
        install_file "etc/rc.d/relique-client" "usr/local/etc/rc.d/relique-client" 1
    fi
}


function install_default_modules() {
    log "Installing default relique modules"

    log "--- Looking for relique binaries in '${SRC}'"

    if [ "X${INSTALL_SERVER}X" == "X1X" ]; then
        RELIQUE_BINARY="${SRC}/bin/relique-server"
    fi

    if [ "X${INSTALL_CLIENT}X" == "X1X" ]; then
        RELIQUE_BINARY="${SRC}/bin/relique-client"
    fi

    if [ "X${RELIQUE_BINARY}X" == "XX" ]; then
        log "ERROR: Cannot find relique binary to install default modules"
        return
    fi

    log "--- Using '${RELIQUE_BINARY}' as relique binary to install default modules"
    for mod in $(ls "${SRC}"/var/lib/relique/default_modules/*.tar.gz); do
        log "--- Install relique module '$(basename ${mod})'"
        ${RELIQUE_BINARY} module install --local --archive -p "${PREFIX}/${ROOT_DATA_PATH}/modules/" --force --skip-chown $mod
    done
}

function main() {
    trap log_exit EXIT
    cmdline "${@}"

    log "Using '${ROOT_CFG_PATH}' as root configuration path"
    log "Using '${ROOT_DATA_PATH}' as root data path"

    copy_binaries
    create_dir_structure
    copy_default_configuration
    copy_certs
    copy_webui_files

    if [ "X${SYSTEMD}X" == "X1X" ]; then
        install_systemd_service
    fi

    if [ "X${FREEBSD}X" == "X1X" ]; then
        install_freebsd_service
    fi

    if [ "X${SKIPUSERCREATION}X" != "X1X" ]; then
        ${BASE}/create_user.sh
        setup_files_ownership
    fi

    if [ "X${SKIPMODULEINSTALL}X" != "X1X" ]; then
        install_default_modules
    fi

    log "Relique distribution installed in '${PREFIX}'. Please check logs for any errors"
}

main "${@}"
