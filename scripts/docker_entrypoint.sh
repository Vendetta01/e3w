#!/bin/bash

# Exit immmideiatley if a command exits with a non-zero status.
set -e

# Variables
source /usr/bin/environment.sh

# functions
function set_config {
    echo "INFO: Initialize config from ENV"
    if [[ -z "${E3W_AUTH}" ]]; then
	E3W_AUTH=false
    fi
    sed -i 's|auth=.*$|auth='${E3W_AUTH}'|' ${E3W_CONF_URL}

    sed -i 's|addr=.*$|addr='${E3W_ADDR}'|' ${E3W_CONF_URL}
    sed -i 's|root_key=.*$|root_key='${E3W_ROOT_KEY}'|' ${E3W_CONF_URL}
    sed -i 's|dir_value=.*$|dir_value='${E3W_DIR_VALUE}'|' ${E3W_CONF_URL}

    echo "DEBUG: SECRET_SSL_CRT: '${SECRET_SSL_CRT}': $(ls ${SECRET_SSL_CRT})"
    echo "DEBUG: SECRET_SSL_KEY: '${SECRET_SSL_KEY}': $(ls ${SECRET_SSL_KEY})"
    echo "DEBUG: SECRET_ROOT_CRT: '${SECRET_ROOT_CRT}': $(ls ${SECRET_ROOT_CRT})"
    if [[ -f "${SECRET_SSL_CRT}" && -f "${SECRET_SSL_KEY}" && -f ${SECRET_ROOT_CRT} ]]; then
	echo "Found ssl cert, key and ca in /run/secrets: setting conf..."
	sed -i 's|cert_file=.*$|cert_file='${SECRET_SSL_CRT}'|' ${E3W_CONF_URL}
	sed -i 's|key_file=.*$|key_file='${SECRET_SSL_KEY}'|' ${E3W_CONF_URL}
	sed -i 's|ca_file=.*$|ca_file='${SECRET_ROOT_CRT}'|' ${E3W_CONF_URL}
    fi

}


# main
if [[ ! -e "$FIRST_START_FILE_URL" ]]; then
	# Do stuff
	set_config
	touch "$FIRST_START_FILE_URL"
fi


# Start etcd
E3W_ADDR=$(sed -n 's|addr=\(.*\)$|\1|p' ${E3W_CONF_URL})

echo "DEBUG: E3W_ADDR: '${E3W_ADDR}'"

echo "INFO: Waiting for etcd server..."
/usr/bin/wait-for-it.sh ${E3W_ADDR} -s -t ${E3W_STARTUP_TIMEOUT}

if [[ $? -ne 0 ]]; then
    echo "ERROR: ectd server not found, exiting..."
    exit 1
fi

echo "DEBUG: '$(cat ${E3W_CONF_URL})'"

echo "INFO: etcd server found, starting e3w"
exec /app/e3w "$@"
