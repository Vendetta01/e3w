#!/bin/bash

# Variables
export FIRST_START_FILE_URL=/tmp/first_start_done
export DEFAULT_LOG_LEVEL="info"

declare -A LOG_LEVELS
LOG_LEVELS=([none]=0 [error]=10 [info]=20 [warn]=30 [debug]=40)


if [[ ! ${LOG_LEVEL+x} ]]; then
    LOG_LEVEL=${DEFAULT_LOG_LEVEL}
fi
LOG_LEVEL_I=${LOG_LEVELS[$LOG_LEVEL]}

# functions
logit () {
    param_log_level=$(echo "$1" | tr '[:upper:]' '[:lower:]')
    param_log_level_i=${LOG_LEVELS[$param_log_level]}
    if [[ ${param_log_level_i} -le ${LOG_LEVEL_I} ]]; then
	echo "$(date -Iseconds) | ${1^^} | ${BASHPID} | ${@:2}"
    fi
}
