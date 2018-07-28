#!/bin/sh
set -ex

go get -v -d -t ./...
go test -v -race -cover ./...
# go test -v -cover ./tools
# go test -v -cover tools/context_test.go
