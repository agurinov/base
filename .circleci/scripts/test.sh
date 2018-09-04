#!/bin/sh
set -ex

go get -d -t ./...
go test -race -cover ./...
# go test -v -race -cover ./pipeline
# go test -race -cover tools/context_test.go
