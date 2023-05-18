#!/bin/bash

#
# Builds conjur-service-broker binaries
# usage: ./scripts/build.sh
set -ex
cd "$(dirname "$0")"
. utils

TOPLEVEL_DIR="$(repo_root)"

FULL_VERSION_TAG="$(full_version_tag)"

DOCKER_FLAGS=("--build-arg" "TAG=$(git_tag)" "--build-arg" "BUILDER_IMAGE_VERSION=$(go_version)-alpine")

function main() {
	#  retrieve_cyberark_ca_cert
	build_in_docker
}

function build_in_docker() {
	announce "Building conjur-service-broker:$FULL_VERSION_TAG"
	# NOTE: the latest tag is required by downstream pipeline stages
	# (we want the flags to be word split here)
	docker build --tag "conjur-service-broker:${FULL_VERSION_TAG}" \
		--tag "conjur-service-broker:latest" \
		--pull \
		"${DOCKER_FLAGS[@]}" \
		--file "$TOPLEVEL_DIR/Dockerfile" \
		"$TOPLEVEL_DIR"
	CONTAINER_ID=$(docker create conjur-service-broker:latest)
	docker cp "$CONTAINER_ID":/opt/conjur_service_broker/conjur_service_broker "$TOPLEVEL_DIR"
	docker rm "$CONTAINER_ID"
}

main
