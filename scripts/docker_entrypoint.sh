#!/bin/bash

# Exit immediatley if a command exits with a non-zero status.
set -e

# Variables
source /usr/bin/environment.sh

# functions
function initialize {
    logit "INFO" "Initializing..."

    # check if server key and cert files are provided and exist
    if [[ ${CONFD__E3W__APP__CERT_FILE+x} && ! -f ${CONFD__E3W__APP__CERT_FILE} ]]; then
        logit "DEBUG" "CONFD__E3W__APP__CERT_FILE: ${CONFD__E3W__APP__CERT_FILE}: file not found!"
    fi
    if [[ ${CONFD__E3W__APP__KEY_FILE+x} && ! -f ${CONFD__E3W__APP__KEY_FILE} ]]; then
        logit "DEBUG" "CONFD__E3W__APP__KEY_FILE: ${CONFD__E3W__APP__KEY_FILE}: file not found!"
    fi

    # initial set up of confd from environemnt
    logit "INFO" "Initial configuration of confd from environment..."
    confd -onetime -sync-only -backend env -confdir /tmp/etc/confd/

    # wait for etcd service to start
    if [[ ! ${CONF__CONFD__E3W__NODES__1+x} || ${CONF__CONFD__E3W__NODES__1} == "" ]]; then
        logit "ERROR" "No ETCD node for configuration specified (CONF__CONFD__E3W__NODES__1)"
        exit 1
    fi
    /usr/bin/wait-for-it.sh ${CONF__CONFD__E3W__NODES__1} -s -t ${E3W_STARTUP_TIMEOUT}

    #sleep 30

    # inital configuration of e3w from etcd
    logit "INFO" "Inital configuration of e3w..."
    confd -onetime -sync-only -log-level debug

    logit "INFO" "Initialization successful"
}


# main
if [[ ! -e "$FIRST_START_FILE_URL" ]]; then
	initialize
fi


logit "INFO" "Starting supervisord..."
#while [[ 1 ]]; do
#    sleep 10
#done
exec /usr/bin/supervisord -c /etc/supervisor/conf.d/supervisord.conf
