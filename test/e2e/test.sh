#!/bin/bash

set -e
cd "$(dirname "$0")"
. ../../scripts/utils
. tanzucli
. ip_manager

export HAMMER_TARGET_CONFIG="${HAMMER_TARGET_CONFIG:-"./hammerfile.json"}"

trap remove_ip EXIT

function main() {
	set -e
	announce "E2E Tests starting"
	check_env
	tanzucli ./test/e2e/tanzu_data.sh
	. tanzu_data
	allow_ip
	./test_in_docker.sh
}

function check_env() {
	banner "checking required environment variables"
	required_env_vars=(
		IPMANAGER_TOKEN
		PCF_CONJUR_ACCOUNT
		PCF_CONJUR_APPLIANCE_URL
		PCF_CONJUR_USERNAME
		PCF_CONJUR_API_KEY
	)
	for env_var in "${required_env_vars[@]}"; do
		if [[ -z "${env_var}" ]]; then
			echo "need to set $env_var"
			exit 1
		fi
	done
	echo "required env variables are set"
}

main
