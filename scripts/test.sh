#!/bin/bash

set -e

rm -f "$PWD"/coverage/unit/* "$PWD"/coverage/integration/* "$PWD"/coverage/merged/* "$PWD"/coverage/all "$PWD"/coverage/all_no_gen &> /dev/null || true

# ignore mocks in results
PACKAGES=()
while IFS=$'\n' read -r pkg; do
  echo "*$pkg"
  PACKAGES+=("$pkg")
done < <(go list ./... | grep -v /mocks)

echo "unit test"
go test -cover "${PACKAGES[@]}" -args -test.gocoverdir="$PWD/coverage/unit"
echo
go tool covdata percent -i ./coverage/unit

echo "integration tests"
go test -tags=integration -cover "${PACKAGES[@]}" -args -test.gocoverdir="$PWD"/coverage/integration
echo
go tool covdata percent -i ./coverage/integration

echo "combined results"
go tool covdata merge -i=./coverage/unit,./coverage/integration -o coverage/merged
go tool covdata percent -i ./coverage/merged
echo

go tool covdata textfmt -i=./coverage/merged -o coverage/all

# ignore generated code
< coverage/all grep -v ".gen.go" > coverage/all_no_gen

go tool cover -func coverage/all_no_gen

set +e
