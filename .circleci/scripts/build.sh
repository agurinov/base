#!/bin/sh
set -ex

go get -v -d ./...

BASE=`basename "$PWD"`
TIMESTAMP=`date +%s`
VERSION="LOCAL (${TIMESTAMP})"

# linux
GOOS=linux GOARCH=amd64 go build \
	-v \
	-ldflags "-X 'main.VERSION=${VERSION}' -X 'main.TIMESTAMP=${TIMESTAMP}'" \
	-o /go/bin/${BASE}-Linux-x86_64

# macos
GOOS=darwin GOARCH=amd64 go build \
	-v \
	-ldflags "-X 'main.VERSION=${VERSION}' -X 'main.TIMESTAMP=${TIMESTAMP}'" \
	-o /go/bin/${BASE}-Darwin-x86_64
