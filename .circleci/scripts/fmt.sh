#!/bin/sh
set -e

go fix ./...
go vet ./...
go fmt ./...