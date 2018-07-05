#!/bin/sh
set -e

go get -v -d -t ./...
go test -cpuprofile cpu.prof -memprofile mem.prof ./pipeline
