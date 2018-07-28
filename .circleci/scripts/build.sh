#!/bin/sh
set -ex

go get -v -d ./...

TIMESTAMP=`date +%s`
VERSION="LOCAL (${TIMESTAMP})"

# linux
GOOS=linux GOARCH=amd64 go build \
	-v \
	-ldflags "-X 'main.VERSION=${VERSION}' -X 'main.TIMESTAMP=${TIMESTAMP}'" \
	-o /go/bin/base-Linux-x86_64

# macos
GOOS=darwin GOARCH=amd64 go build \
	-v \
	-ldflags "-X 'main.VERSION=${VERSION}' -X 'main.TIMESTAMP=${TIMESTAMP}'" \
	-o /go/bin/base-Darwin-x86_64
