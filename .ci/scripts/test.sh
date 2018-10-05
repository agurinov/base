#!/bin/sh
set -ex

go get -d -t ./...
go clean -testcache
go test -race -cover ./...
# go test -race -cover ./tools/chronometer
# go test -race -cover tools/context_test.go
