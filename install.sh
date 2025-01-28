#!/bin/bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

success_print(){
    printf "${GREEN}$1${NC}\n"
}

warning_print(){
    printf "${YELLOW}$1${NC}\n"
}

error_print(){
    printf "${RED}$1${NC}\n"
}

info_print(){
    printf "${BLUE}$1${NC}\n"
}


is_root(){
    if [ "$(id -u)" != "0" ]; then
        printf "${RED}Please run this script as root${NC}\n" >&2
        exit 1
    fi
}


has_command(){
    command -v "$@" > /dev/null 2>&1
}

install_docker(){
    # Check if Docker is installed
    if has_command docker; then
        echo "Docker already installed"
    else
        curl -fsSL https://get.docker.com | sh
    fi
    
}

is_port_in_use (){
    if ss -tulnp | grep ":$1 " >/dev/null; then
        error_print "Error: something is already running on port $1" >&2
        exit 1
    fi
}

ipv4(){
    local ip=""
    
    # Try IPv4 first
    # First attempt: ifconfig.io
    ip=$(curl -4s --connect-timeout 5 https://ifconfig.io 2>/dev/null)
    
    # Second attempt: ipgrab.io
    if [ -z "$ip" ]; then
        ip=$(curl -4s --connect-timeout 5 https://ipgrab.io 2>/dev/null)
    fi

    # Third attempt: ipify.org
    if [ -z "$ip" ]; then
        ip=$(curl --connect-timeout 5 https://api.ipify.org 2>/dev/null)
    fi

    # Fourth attempt: ipecho.net
    if [ -z "$ip" ]; then
        ip=$(curl -4s --connect-timeout 5 https://ipecho.net/plain 2>/dev/null)
    fi

    echo "$ip"
}

ipv6(){
    local ip=""
    
    # Try IPv6 first
    # First attempt: ifconfig.io
    ip=$(curl -6s --connect-timeout 5 https://ifconfig.io 2>/dev/null)
    
    # Second attempt: ipgrab.io
    if [ -z "$ip" ]; then
        ip=$(curl -6s --connect-timeout 5 https://ipgrab.io 2>/dev/null)
    fi

    # Third attempt: ipify.org
    if [ -z "$ip" ]; then
        ip=$(curl -6s --connect-timeout 5 https://api6.ipify.org 2>/dev/null)
    fi

    # Fourth attempt: ipecho.net
    if [ -z "$ip" ]; then
        ip=$(curl -6s --connect-timeout 5 https://ipecho.net/plain 2>/dev/null)
    fi

    echo "$ip"
}

install_shellby(){

    mkdir -p /etc/shellby
    chmod 777 /etc/shellby

    # Check if something is running on port 80
    is_port_in_use 80

    # Check if something is running on port 443
    is_port_in_use 443

    install_docker

    docker swarm leave --force 2>/dev/null

    
    local ip=$(ipv4)
    if [ -z "$ip" ]; then
        ip=$(ipv6)
    fi

    # Check if the IP is empty
    if [ -z "$ip" ]; then
        error_print "Error: Failed to get the IP address" >&2
        exit 1
    fi

    info_print "IP address: $ip"

    info_print "Initializing Docker Swarm"
    
    if ! docker swarm init --advertise-addr "$ip"; then
        error_print "Error: Failed to initialize Docker Swarm" >&2
        exit 1
    fi

    success_print "Swarm initialized"
    info_print "Advertising address: $ip"

    # Create the shellby network
    docker network rm -f shellby-network 2>/dev/null
    docker network create --driver=overlay --attachable shellby-network

    success_print "Shellby network created"

    # Pull 
    docker pull traefik:v3.3.2
    docker pull postgres:17.2

    # docker service create --name shellby --network shellby-network --mount type=volume,source=shellby-volume,target=/shellby -p 80:80 -p 443:443 shellby
    # success_print "Shellby service created"

    success_print "Shellby is now running at https://$ip"
}

install_shellby