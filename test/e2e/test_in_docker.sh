#!/bin/bash

#
# executes conjur-service-broker E2E tests in docker
# usage: ./test/e2e/test_in_docker.sh
set -ex
cd "$(dirname "$0")"
. ../../scripts/utils

# with that we would use the docker image compatible with go.mod go version
DOCKER_FLAGS=("--build-arg" "BUILDER_IMAGE_VERSION=$(go_version)-alpine")

function main() {
	test_in_docker
}

function test_in_docker() {
	set -e
	announce "Executing E2E tests in Docker"
	# build the test image
	IMAGE_ID=$(docker build \
		--pull \
		-q \
		"${DOCKER_FLAGS[@]}" \
		--file "Dockerfile.test" \
		.)
	# copy the artifact to test, the goreleaser is using --clean flag so there would be exactly one artifact in ./dist
	cp ../../dist/goreleaser/*.zip cyberark-conjur-service-broker.zip
	# execute the tests, inject needed configuration, mount the artifact and results report
	summon docker run \
		-v "$PWD"/reports:/opt/e2e/reports \
		-v "$PWD"/cyberark-conjur-service-broker.zip:/opt/e2e/cyberark-conjur-service-broker.zip:ro \
		-e CF_USERNAME \
		-e CF_PASSWORD \
		-e CF_API_URL \
		-e PCF_CONJUR_ACCOUNT \
		-e PCF_CONJUR_APPLIANCE_URL \
		-e PCF_CONJUR_USERNAME \
		-e PCF_CONJUR_API_KEY \
		--rm \
		"$IMAGE_ID"
}

main
