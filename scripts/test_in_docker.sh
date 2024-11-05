#!/bin/bash

#
# executes conjur-service-broker tests in docker
# usage: ./scripts/test_in_docker.sh
set -ex
cd "$(dirname "$0")"
. utils

TOPLEVEL_DIR="$(repo_root)"

DOCKER_FLAGS=("--build-arg" "BUILDER_IMAGE_VERSION=$(go_version)-alpine")

SKIP_GOMOD_DOWNLOAD=false

function main() {
	# When running in Jenkins, we need to skip the go mod download and go mod tidy
	# commands since we already fetch the latest dependencies with
	# updatePrivateGoDependencies(). We use the --skip-gomod-download flag for this purpose.

	set +u
	while true ; do
		case "$1" in
			--skip-gomod-download ) SKIP_GOMOD_DOWNLOAD=true ; shift ;;
			* ) if [ -z "$1" ]; then break; else echo "$1 is not a valid option"; exit 1; fi;;
		esac
	done
	set -u

	test_in_docker
}

function test_in_docker() {
	announce "Executing tests in Docker"
	IMAGE_ID=$(docker build \
		--pull \
		-q \
		"${DOCKER_FLAGS[@]}" \
		--file "Dockerfile.test" \
		"$TOPLEVEL_DIR" || exit 1)
	docker run --env "SKIP_GOMOD_DOWNLOAD=$SKIP_GOMOD_DOWNLOAD" -v "$TOPLEVEL_DIR"/coverage:/src/coverage --rm "$IMAGE_ID"
}

main "$@"
