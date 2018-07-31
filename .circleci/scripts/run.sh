#!/bin/sh
set -ex

# PROJECT=${1:?'No `project_name` provided '}
PROJECT='geo'

apk add --update --no-cache git

# build base
# fix, lint and build source code of base app
.circleci/scripts/fmt.sh
.circleci/scripts/build.sh

# build microservice related src
# go to micriservice src root
cd ${PROJECT}

# fix, lint and build source code of microservice
../.circleci/scripts/fmt.sh
../.circleci/scripts/build.sh
# copy bin from /go/bin to our dir (conf needs it here)
cp /go/bin/${PROJECT}-$(uname -s)-$(uname -m) ./${PROJECT}

# set application variables for run base
# TODO move to special config in future named boomfunc.yaml
# TODO https://github.com/urfave/cli#values-from-alternate-input-sources-yaml-toml-and-others
export BMP_DEBUG_MODE=true
export BMP_CONFIG='./conf.yml'
export BMP_APPLICATION_LAYER='http'

# run base
/go/bin/base-Linux-x86_64 run tcp
