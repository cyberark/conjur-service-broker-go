#!/usr/bin/env bash

set -euo pipefail

function cleanup() {
  containers="$(docker network inspect --format '{{range .Containers}}{{.Name}} {{end}}' kind-network)"
  for container in $(echo $containers); do
    docker rm -f -v "$container"
  done

  docker network rm kind-network
}

trap cleanup EXIT ERR SIGINT

# When running in Jenkins, we need to skip the go mod download and go mod tidy
# commands since we already fetch the latest dependencies with
# updatePrivateGoDependencies(). We use the --skip-gomod-download flag for this purpose.
SKIP_GOMOD_DOWNLOAD=false
set +u
while true ; do
  case "$1" in
    --skip-gomod-download ) SKIP_GOMOD_DOWNLOAD=true ; shift ;;
    * ) if [ -z "$1" ]; then break; else echo "$1 is not a valid option"; exit 1; fi;;
  esac
done
set -u

docker network create kind-network

docker build \
  -f Dockerfile.integration \
  -t conjur-service-broker-test-runner:latest \
  .

docker run --rm \
  --network kind-network \
  --env "SKIP_GOMOD_DOWNLOAD=$SKIP_GOMOD_DOWNLOAD" \
  -v "$PWD":/src \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -w /src \
  conjur-service-broker-test-runner:latest \
  bash -c "
    ./scripts/init_kind.sh && tilt ci
  "
exit_code="$?"

exit $exit_code
