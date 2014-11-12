#!/bin/bash -ex

cd /go/src/github.com/joho/shoutcloud_api
go get ./...
rm -f $GOBIN/shoutcloud_api
go install -ldflags '-extldflags "-static"' github.com/joho/shoutcloud_api

