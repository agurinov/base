#!/bin/sh
set -ex

go get -d -t ./...
go test -race -run=^$$ -bench=. -benchmem ./...
