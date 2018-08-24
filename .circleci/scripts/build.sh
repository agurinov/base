#!/bin/sh

<< ////

TODO:
use cross-compiling
$go get github.com/mitchellh/gox
$gox \
	-verbose \
	-osarch="linux/amd64 darwin/amd64" \
	-ldflags="-X main.VERSION=$CIRCLE_TAG -X main.TIMESTAMP=`date +%s`" \
	-output="/go/bin/{{.Dir}}-{{.OS}}-{{.Arch}}"

////

set -ex

go get -d ./...

# calculate base variables
BASE=`basename "$PWD"`
TIMESTAMP=`date +%s`
VERSION="${CIRCLE_TAG:=LOCAL}"

# calculate build/compile specific variables
ldflags="-X 'main.VERSION=${VERSION}' -X 'main.TIMESTAMP=${TIMESTAMP}'"

# linux 64 bit
# - alpine
# - ubuntu
GOOS=linux GOARCH=amd64 go build \
	-v \
	-ldflags "${ldflags}" \
	-o /go/bin/${BASE}-Linux-x86_64

# macos
GOOS=darwin GOARCH=amd64 go build \
	-v \
	-ldflags "${ldflags}" \
	-o /go/bin/${BASE}-Darwin-x86_64
