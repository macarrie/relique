#!/usr/bin/env bash

function usage() {
    echo "\
usage: $0 [options]
    
Options:
    -h --help: Displays this help
    -p --prefix: Install relique to this folder
    -s --src: Get compiled relique package to install from this folder
    "
}

function install_file() {
    src_file=$1
    overwrite=$2

    if [ -f "${PREFIX}/${src_file}" ] && [ "X${overwrite}X" != "X1X" ]; then
        echo "--- Skipping ${src_file} copy. File already exists"
        return
    fi

    echo "--- Copying ${src_file}"

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
    fi

    if [ "X${INSTALL_CLIENT}X" == "X1X" ]; then
        install_file "etc/relique/client.toml"
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
    echo "TODO: Install systemd service"
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

copy_binaries
copy_default_configuration

echo "TODO: Install"
