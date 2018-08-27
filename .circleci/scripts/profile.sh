#!/bin/sh
set -ex

go get -d -t ./...
go test -cpuprofile cpu.prof -memprofile mem.prof ./pipeline
