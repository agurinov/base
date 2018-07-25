#!/bin/sh
set -e

apk add --update --no-cache git

.circleci/scripts/build.sh

cd geo
go get -v -d
go fmt
go build geoip.go
cd ../

/go/bin/base-Linux-x86_64 run --app=http tcp
