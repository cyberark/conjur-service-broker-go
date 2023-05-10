#!/bin/bash

set -e

# gosec - inspects source code for security problems by scanning the Go AST
go install github.com/securego/gosec/v2/cmd/gosec@latest

# staticcheck - state of the art linter for the Go programming language. Using static analysis, it finds bugs and performance issues, offers simplifications, and enforces style rules
go install honnef.co/go/tools/cmd/staticcheck@latest

# goimports - updates your Go import lines, adding missing ones and removing unreferenced ones
go install golang.org/x/tools/cmd/goimports@latest

# golint
go install golang.org/x/lint/golint@latest

# revive ~6x faster, stricter, configurable, extensible, and beautiful drop-in replacement for golint
go install github.com/mgechev/revive@latest

# critic - the most opinionated Go source code linter for code audit
go install github.com/go-critic/go-critic/cmd/gocritic@latest

# golangci - FAST linter aggregator
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# install shellcheck
brew install shellcheck

# shfmt
brew install shfmt

pre-commit install

set +e
