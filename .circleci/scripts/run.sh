#!/bin/sh
set -e

apk add --update --no-cache git python py-pip
pip install --no-cache-dir pygeoip==0.3.2

.circleci/scripts/build.sh

/go/bin/base-Linux-x86_64 run rpc
