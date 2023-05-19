#!/bin/bash

#
# executes conjur-service-broker tests in docker
# usage: ./scripts/test_in_docker.sh
set -ex
cd "$(dirname "$0")"
. utils

TOPLEVEL_DIR="$(repo_root)"

DOCKER_FLAGS=("--build-arg" "BUILDER_IMAGE_VERSION=$(go_version)-alpine")

function main() {
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
	docker run -v "$TOPLEVEL_DIR"/coverage:/src/coverage --rm "$IMAGE_ID"
}

main
