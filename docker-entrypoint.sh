#!/bin/sh
set -e

# run dockerd
/usr/local/bin/dockerd-entrypoint.sh &

# wait for docker daemon to be ready before continuing
while true ; do if [ -S /var/run/docker.sock ] ; then break ; fi ; sleep 1 ; done

# build infinity manifest
infinity build