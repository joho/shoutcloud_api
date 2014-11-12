#!/bin/bash -ex

cd "$( dirname "${BASH_SOURCE[0]}" )"

./build.sh

# Push release docker image to registry
DEPLOY_TAG=$(git rev-parse HEAD)

docker tag shoutcloud/shoutcloud_api shoutcloud/shoutcloud_api:$DEPLOY_TAG

docker push shoutcloud/shoutcloud_api:$DEPLOY_TAG
docker push shoutcloud/shoutcloud_api:latest
