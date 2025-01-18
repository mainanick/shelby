#!/bin/bash

set -e
set -o pipefail

is_root(){
    if [ "$(id -u)" != "0" ]; then
        echo "This script must be run as root" >&2
        exit 1
    fi
}

has_command(){
    command -v "$@" > /dev/null 2>&1
}

install_docker(){
    curl -fsSL https://get.docker.com | sh
}

install_shellby(){
    # Check if Docker is installed
    if has_command docker; then
        echo "Docker already installed"
    else
        install_docker
    fi


}