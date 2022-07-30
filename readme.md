# Relique

Relique is a backup tool based on rsync. It as basically a wrapper around rsync destined to simplify backup management.

Instead of defining manually which files you want to back up foreach client, you use modules that contain file paths, pre/post backup and pre/post restore scripts that do the work for you. For example, you use the `plex` module that contains all the informations needed to backup and restore a Plex Media Server.

# Getting started

## Installation

### FreeBSD system package

You can use the Makefile present in build/package/freebsd/relique-* to build a port and generate .pkg installable file.

### RPM system package

Relique can also be built and used on RPM compatible Linux systems by using the build/package/rpm/relique-*.spec files.

### Docker

Relique server and client can be run by using the associated Docker images:
- macarrie/relique-server on Docker Hub
- macarrie/relique-client on Docker Hub
- Build fresh image from Dockerfiles (build/package/docker/*/Dockerfile)
- Use sample Docker Compose configuration files (build/package/docker/*/docker-compose.yaml)

### Manual build and install

Depencies needed for build and install:
- go 1.18
- Make
- Bash

Dependencies needed for running relique server:
- rsync

To build and install relique client/server, use the following commands:
- Get relique project sources via git clone: `git clone --recurse-submodules github.com/macarrie/relique`
- `make build`
- `make install INSTALL_ARGS="--server"`
- `make install INSTALL_ARGS="--client"`

The `INSTALL_ARGS` parameters can have different values passed directly to the scripts/install.sh script:
- "--server": install relique server
- "--client": install relique client
- "--systemd": install systemd service files
- "--freebsd": install FreeBSD service files

## Configuration

TODO
