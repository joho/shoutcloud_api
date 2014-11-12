#!/bin/bash -ex

WORKING_DIR="$( cd "$(dirname "$0")" ; pwd -P )"
cd $WORKING_DIR

docker build -t shoutcloud/shoutcloud_api_build build

# You could probably use $GOPATH/src
GO_SRC_PATH=$HOME/go/src
BIN_VOLUME_PATH=$WORKING_DIR/.bin

docker run \
    --volume $HOME/Projects/go/src:/go/src \
    --volume $BIN_VOLUME_PATH:/go/bin \
    --tty --interactive \
    shoutcloud/shoutcloud_api_build \
    /provision.sh

cp .bin/shoutcloud_api release

docker build -t shoutcloud/shoutcloud_api release
