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

docker network create kind-network

docker build \
  -f Dockerfile.integration \
  -t conjur-service-broker-test-runner:latest \
  .

docker run --rm \
  --network kind-network \
  -v "$PWD":/src \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -w /src \
  conjur-service-broker-test-runner:latest \
  bash -c "
    ./scripts/init_kind.sh && tilt ci
  "
exit_code="$?"

exit $exit_code
