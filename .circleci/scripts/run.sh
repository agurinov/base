#!/bin/sh
set -e

apk add --update git

.circleci/scripts/build.sh

/go/bin/base-Linux-x86_64 run tcp
