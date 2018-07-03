#!/bin/sh
go get -v -d ./...
GOOS=darwin GOARCH=amd64 go build \
	-v \
	-ldflags "-X main.VERSION=local -X main.TIMESTAMP=`date +%s`" \
	-o /go/bin/base-Darwin-x86_64
