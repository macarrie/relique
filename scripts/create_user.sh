#!/usr/bin/env bash

RELIQUE_USER=relique
RELIQUE_UID=833
RELIQUE_GROUP=relique
RELIQUE_GID=833
RELIQUE_USER_HOME=/var/lib/relique

function usage() {
    echo "\
usage: $0 [options]
    
Options:
    -h --help: Displays this help
    -u --user: user to create (default ${RELIQUE_USER} )
	--uid: user ID (default ${RELIQUE_UID})
	-g --group: group to create (default ${RELIQUE_GROUP})
	--gid: group ID (default ${RELIQUE_GID})
	--home: User homedir (default ${RELIQUE_USER_HOME})
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

    --uid)
    RELIQUE_UID="$2"
    shift # past argument
    shift # past value
    ;;

    -g|--group)
    RELIQUE_GROUP="$2"
    shift # past argument
    shift # past value
    ;;

    --gid)
    RELIQUE_GID="$2"
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

useradd_path=$(which useradd)
groupadd_path=$(which groupadd)
pw_path=$(which pw)

echo "Creating group '${RELIQUE_GROUP}' (${RELIQUE_GID})"
if [ -f "${groupadd_path}" ]; then
	getent group ${RELIQUE_GROUP} >/dev/null || groupadd -r "${RELIQUE_GROUP}" -g $RELIQUE_GID
elif [ -f "${pw_path}" ]; then
	getent group ${RELIQUE_GROUP} >/dev/null || pw groupadd "${RELIQUE_GROUP}" -g $RELIQUE_GID
else
	echo "ERROR: Cannot create relique group (system unsupported)"
	exit 1
fi

echo "Creating user '${RELIQUE_USER}' (${RELIQUE_UID})"
if [ -f "${useradd_path}" ]; then
	getent passwd ${RELIQUE_USER} >/dev/null || \
    useradd -r -g "${RELIQUE_GROUP}" --uid $RELIQUE_UID -d "${RELIQUE_USER_HOME}" -s /sbin/nologin \
    -c "Relique service account" "${RELIQUE_USER}"
elif [ -f "${pw_path}" ]; then
	getent passwd ${RELIQUE_GROUP} >/dev/null || pw useradd "${RELIQUE_USER}" -g "${RELIQUE_GROUP}" -u $RELIQUE_UID
else
	echo "ERROR: Cannot create relique user (system unsupported)"
	exit 1
fi

exit 0
