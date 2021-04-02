#!/usr/bin/env bash

USER=relique
GROUP=relique

function usage() {
    echo "\
usage: $0 [options]
    
Options:
    -h --help: Displays this help
    -p --prefix: Install relique to this folder
    -s --src: Get compiled relique package to install from this folder
    --server: Install relique server
    --client: Install relique client
    --systemd: Install systemd service file
    --skip-user-creation: Skip relique group and user creation
    "
}

function install_file() {
    src_file=$1
    overwrite=$2

    if [ -f "${PREFIX}/${src_file}" ] && [ "X${overwrite}X" != "X1X" ]; then
        echo "--- Skipping ${PREFIX}/${src_file} copy. File already exists"
        return
    fi

    echo "--- Copying ${src_file} to ${PREFIX}/${src_file}"

    dest_path=$(dirname $src_file)
    if [ ! -d "${dest_path}" ]; then
        mkdir -p "${PREFIX}/${dest_path}"
    fi

    cp "${SRC}/${src_file}" "${PREFIX}/${src_file}"
}

function copy_binaries() {
    echo -e "\nInstalling binaries"
    install_file "usr/bin/relique" 1

    if [ "X${INSTALL_SERVER}X" == "X1X" ]; then
        install_file "usr/bin/relique-server" 1
    fi

    if [ "X${INSTALL_CLIENT}X" == "X1X" ]; then
        install_file "usr/bin/relique-client" 1
    fi
}

function copy_default_configuration() {
    echo -e "\nInstalling default configuration"

    if [ "X${INSTALL_SERVER}X" == "X1X" ]; then
        install_file "etc/relique/server.toml"
        mkdir -p "${PREFIX}/opt/relique"
        mkdir -p "${PREFIX}/var/lib/relique"
        mkdir -p "${PREFIX}/var/log/relique"
    fi

    if [ "X${INSTALL_CLIENT}X" == "X1X" ]; then
        install_file "etc/relique/client.toml"
    fi
}


function copy_certs() {
    echo -e "\nInstalling self signed quick start certs"

    install_file "etc/relique/certs/cert.pem"
    install_file "etc/relique/certs/key.pem"
}


function create_user {
    echo -e "\nCreating $USER user"

    id -u $USER > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        useradd -M "${USER}"
    fi
}


function create_dir_structure {
    echo -e "\nCreating relique directory structure"

    mkdir -p "${PREFIX}/etc/relique"
    mkdir -p "${PREFIX}/var/lib/relique"
    mkdir -p "${PREFIX}/var/log/relique"
    mkdir -p "${PREFIX}/opt/relique"
}


function setup_files_ownership() {
    echo -e "\nSetting files rights and ownership"

    chown -R $USER:$GROUP "${PREFIX}/etc/relique"
    chown -R $USER:$GROUP "${PREFIX}/var/lib/relique"
    chown -R $USER:$GROUP "${PREFIX}/var/log/relique"
    chown -R $USER:$GROUP "${PREFIX}/opt/relique"
}


function install_systemd_service() {
    echo -e "\nInstalling systemd service files"

    if [ "X${INSTALL_SERVER}X" == "X1X" ]; then
        install_file "usr/lib/systemd/system/relique-server.service"
    fi

    if [ "X${INSTALL_CLIENT}X" == "X1X" ]; then
        install_file "usr/lib/systemd/system/relique-client.service"
    fi
}


POSITIONAL=()
while [[ $# -gt 0 ]]
do
key="$1"

case $key in
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

    --skip-user-creation)
    SKIPUSERCREATION=1
    shift # past argument
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

if [ -z $PREFIX ]; then
    echo "Missing install prefix"
    usage
    exit 1
fi

if [ -z $SRC ]; then
    echo "Missing install source directory"
    usage
    exit 1
fi

if [ ! -d $SRC ]; then
    echo "Source directory '$SRC' does not exist"
    exit 1
fi

if [ ! -d $PREFIX ]; then
    mkdir -p "${PREFIX}"
fi

create_user
copy_binaries
create_dir_structure
copy_default_configuration
copy_certs

if [ "X${SYSTEMD}X" == "X1X" ]; then
    install_systemd_service
fi

if [ "X${SKIPUSERCREATION}X" != "X1X" ]; then
    ./create_user
    setup_files_ownership
fi
