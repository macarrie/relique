#!/usr/bin/env bash

RELIQUE_USER=relique
RELIQUE_GROUP=relique
RELIQUE_USER_HOME=/var/lib/relique

function usage() {
    echo "\
usage: $0 [options]
    
Options:
    -h --help: Displays this help
    -u --user: user to create
    -g --group: group to create
    --home: User homedir
    "
}


POSITIONAL=()
while [[ $# -gt 0 ]]
do
key="$1"

case $key in
    -u|--user)
    RELIQUE_USER="$2"
    shift # past argument
    shift # past value
    ;;

    -g|--group)
    RELIQUE_GROUP="$2"
    shift # past argument
    shift # past value
    ;;

    --home)
    RELIQUE_USER_HOME="$2"
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

echo "Creating group '${RELIQUE_GROUP}'"
getent group ${GROUP} >/dev/null || groupadd -r relique

echo "Creating user '${RELIQUE_USER}'"
getent passwd ${RELIQUE_USER} >/dev/null || \
    useradd -r -g "${RELIQUE_GROUP}" -d "${RELIQUE_USER_HOME}" -s /sbin/nologin \
    -c "Relique service account" "${RELIQUE_USER}"

exit 0