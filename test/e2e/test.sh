#!/bin/bash

set -e
cd "$(dirname "$0")"
. ../../scripts/utils
. tanzucli
. ip_manager

export HAMMER_TARGET_CONFIG="${HAMMER_TARGET_CONFIG:-"$(repo_root)/hammerfile.json"}"

# cleanup after the tests - remove the tanzu worker ip allowed to access aws conjur test instance
trap remove_ip EXIT

function main() {
	set -e
	if [ ! -f "${HAMMER_TARGET_CONFIG}" ]; then
		echo "hammerfile not found in ${HAMMER_TARGET_CONFIG}!"
		exit 1
	fi
	announce "E2E Tests starting"
	# get tanzu credentials and worker ip address
	tanzucli ./test/e2e/tanzu_data.sh
	# load tanzu env data
	. tanzu_data
	# allow tanzu worker ip to access our aws conjur test instance
	allow_ip
	# execute tests in docker
	./test_in_docker.sh
}

main
