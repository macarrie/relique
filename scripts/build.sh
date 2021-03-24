#!/usr/bin/env bash

function usage() {
    echo "\
usage: $0 [options]
    
Options:
    -h --help: Displays this help
    -o --output-dir: Output directory for generated artefacts
    "
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

echo "Building binaries to '$OUTPUT_DIR'"
components="relique relique-client relique-server"

for component in $components; do
    echo "Building $component"
    go build -o "${OUTPUT_DIR}/usr/bin/${component}" cmd/${component}/main.go
done

echo "Copying default configuration files to '$OUTPUT_DIR'"
cp -r configs/* "$OUTPUT_DIR"

echo "Copying systemd service files to '$OUTPUT_DIR'"
mkdir -p "${OUTPUT_DIR}/usr/lib/systemd/system"
cp -r build/init/*.service "${OUTPUT_DIR}/usr/lib/systemd/system"
