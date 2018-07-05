#!/bin/sh
set -e

go get -v -d ./...
GOOS=linux GOARCH=amd64 go build \
	-v \
	-ldflags "-X main.VERSION=local -X main.TIMESTAMP=`date +%s`" \
	-o /go/bin/base-Linux-x86_64

/go/bin/base-Linux-x86_64 run tcp
