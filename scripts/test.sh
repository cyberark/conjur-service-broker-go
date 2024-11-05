#!/bin/bash

set -ex
cd "$(dirname "$0")"
. utils

TOPLEVEL_DIR=$(abs_path "./../..")

rm -f "$TOPLEVEL_DIR"/coverage/unit/* "$TOPLEVEL_DIR"/coverage/integration/* "$TOPLEVEL_DIR"/coverage/merged/* "$TOPLEVEL_DIR"/coverage/all "$TOPLEVEL_DIR"/coverage/all_no_gen &>/dev/null || true

SKIP_GOMOD_DOWNLOAD="${SKIP_GOMOD_DOWNLOAD:-false}"
if [ "$SKIP_GOMOD_DOWNLOAD" != "true" ]; then
	go mod tidy
fi


# ignore mocks in results
PACKAGES=()
while IFS=$'\n' read -r pkg; do
	PACKAGES+=("$pkg")
done < <(go list "$TOPLEVEL_DIR"/... | grep -v /mocks)

echo "unit tests"
{ go test -v -covermode=count -cover "${PACKAGES[@]}" -args -test.gocoverdir="$TOPLEVEL_DIR/coverage/unit" | tee "$TOPLEVEL_DIR/coverage/results"; } 2>&1
echo
go tool covdata percent -i "$TOPLEVEL_DIR"/coverage/unit

echo "integration tests"
{ go test -v -covermode=count -tags=integration -cover "${PACKAGES[@]}" -args -test.gocoverdir="$TOPLEVEL_DIR/coverage/integration" | tee -a "$TOPLEVEL_DIR/coverage/results"; } 2>&1
echo
go tool covdata percent -i "$TOPLEVEL_DIR"/coverage/integration

echo "combined results"
go tool covdata merge -i="$TOPLEVEL_DIR"/coverage/unit,"$TOPLEVEL_DIR"/coverage/integration -o "$TOPLEVEL_DIR"/coverage/merged
go tool covdata percent -i "$TOPLEVEL_DIR"/coverage/merged
echo

go tool covdata textfmt -i="$TOPLEVEL_DIR"/coverage/merged -o "$TOPLEVEL_DIR"/coverage/all

# ignore generated code
grep -v ".gen.go" "$TOPLEVEL_DIR"/coverage/all | grep -v "test_util.go" >"$TOPLEVEL_DIR"/coverage/all_no_gen

go tool cover -func "$TOPLEVEL_DIR"/coverage/all_no_gen

go-junit-report -in "$TOPLEVEL_DIR"/coverage/results -out "$TOPLEVEL_DIR"/coverage/junit.xml

gocover-cobertura <"$TOPLEVEL_DIR"/coverage/all_no_gen >"$TOPLEVEL_DIR"/coverage/cobertura.xml

set +ex
