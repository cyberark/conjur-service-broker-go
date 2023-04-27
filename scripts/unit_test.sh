#!/bin/bash

set -e

go test -cover ./... -args -test.gocoverdir=$PWD/coverage/unit

go test -tags=integration -cover ./... -args -test.gocoverdir=$PWD/coverage/integration

go tool covdata textfmt -i=./coverage/unit,./coverage/integration -o coverage/all

go tool cover -func coverage/all

set +e
