#!/bin/sh
set -ex

apk add --update --no-cache git

# build base
.circleci/scripts/build.sh

# build microservice related src
cd parser
go get -v -d
go fmt
go build main.go
cd ../

# set application variables for run base
export BMP_DEBUG_MODE=true
export BMP_CONFIG='./parser/conf.yml'
export BMP_APPLICATION_LAYER='http'

# run base
/go/bin/base-Linux-x86_64 run tcp
