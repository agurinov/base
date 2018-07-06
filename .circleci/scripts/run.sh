#!/bin/sh
set -e

apk add --update git

./build.sh

/go/bin/base-Linux-x86_64 run tcp
