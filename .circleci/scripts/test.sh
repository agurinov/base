#!/bin/sh
set -ex

go get -d -t ./...
go test -race -cover ./...
# go test -race -cover ./tools
# go test -race -cover tools/context_test.go
