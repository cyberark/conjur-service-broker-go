#!/bin/bash

set -ex

. scripts/build_utils

CURRENT_DIR="$(abs_path "${BASH_SOURCE[0]}")"
TOPLEVEL_DIR=$(cd "$CURRENT_DIR/.." && pwd)

rm -f "$TOPLEVEL_DIR"/coverage/unit/* "$TOPLEVEL_DIR"/coverage/integration/* "$TOPLEVEL_DIR"/coverage/merged/* "$TOPLEVEL_DIR"/coverage/all "$TOPLEVEL_DIR"/coverage/all_no_gen &>/dev/null || true

# ignore mocks in results
PACKAGES=()
while IFS=$'\n' read -r pkg; do
	PACKAGES+=("$pkg")
done < <(go list ./... | grep -v /mocks)

echo "unit tests"
{ go test -v -covermode=count -cover "${PACKAGES[@]}" -args -test.gocoverdir="$TOPLEVEL_DIR/coverage/unit" >"$TOPLEVEL_DIR"/coverage/results; } 2>&1
echo
go tool covdata percent -i "$TOPLEVEL_DIR"/coverage/unit

echo "integration tests"
{ go test -v -covermode=count -tags=integration -cover "${PACKAGES[@]}" -args -test.gocoverdir="$TOPLEVEL_DIR"/coverage/integration >>"$TOPLEVEL_DIR"/coverage/results; } 2>&1
echo
go tool covdata percent -i "$TOPLEVEL_DIR"/coverage/integration

echo "combined results"
go tool covdata merge -i="$TOPLEVEL_DIR"/coverage/unit,"$TOPLEVEL_DIR"/coverage/integration -o "$TOPLEVEL_DIR"/coverage/merged
go tool covdata percent -i "$TOPLEVEL_DIR"/coverage/merged
echo

go tool covdata textfmt -i="$TOPLEVEL_DIR"/coverage/merged -o "$TOPLEVEL_DIR"/coverage/all

# ignore generated code
grep -v ".gen.go" "$TOPLEVEL_DIR"/coverage/all | grep -v "test_util.go" >"$TOPLEVEL_DIR"/coverage/all_no_gen

go tool cover -func coverage/all_no_gen

go-junit-report -in "$TOPLEVEL_DIR"/coverage/results -out "$TOPLEVEL_DIR"/coverage/junit.xml

gocover-cobertura <"$TOPLEVEL_DIR"/coverage/all_no_gen >"$TOPLEVEL_DIR"/coverage/cobertura.xml

set +ex
