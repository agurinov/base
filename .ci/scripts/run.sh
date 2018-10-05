#!/bin/sh
set -ex

# PROJECT=${1:?'No `project_name` provided '}
# PROJECT='geo'

apk add --update --no-cache \
	git \
	imagemagick \
	python \
	# py-pip
	# py-lxml
	# chromium \

# build base
# fix, lint and build source code of base app
.ci/scripts/fmt.sh
.ci/scripts/build.sh

# build microservice related src
# go to micriservice src root
# cd ${PROJECT}

# fix, lint and build source code of microservice
# pip install --no-cache-dir -r requirements.txt
# pip freeze > requirements.txt
# ../.ci/scripts/fmt.sh
# ../.ci/scripts/build.sh
# # copy bin from /go/bin to our dir (conf needs it here)
# cp /go/bin/${PROJECT}-$(uname -s)-$(uname -m) ./${PROJECT}

# set application variables for run base
# TODO move to special config in future named boomfunc.yaml
# TODO https://github.com/urfave/cli#values-from-alternate-input-sources-yaml-toml-and-others
export BMP_DEBUG_MODE=true
# export BMP_CONFIG='./conf/example.yml'
export BMP_CONFIG='https://bitbucket.org/!api/2.0/snippets/letsnetdevinternal/eazARA/files/example.yml'
export BMP_APPLICATION_LAYER='http'
# export BMP_WORKER_NUM=8

# run base
# /go/bin/base-Linux-x86_64 run tick
/go/bin/base-Linux-x86_64 run tcp
