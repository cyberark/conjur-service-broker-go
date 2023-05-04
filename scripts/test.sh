#!/bin/bash

set -e

rm $PWD/coverage/unit/* $PWD/coverage/integration/* $PWD/coverage/all || true

echo "unit test"
go test -cover ./... -args -test.gocoverdir=$PWD/coverage/unit
echo
go tool covdata percent -i ./coverage/unit

echo "integration tests"
go test -tags=integration -cover ./... -args -test.gocoverdir=$PWD/coverage/integration
echo
go tool covdata percent -i ./coverage/integration

echo "combined results"
go tool covdata merge -i=./coverage/unit,./coverage/integration -o coverage/merged
go tool covdata percent -i ./coverage/merged
echo

go tool covdata textfmt -i=./coverage/merged -o coverage/all

go tool cover -func coverage/all

set +e
