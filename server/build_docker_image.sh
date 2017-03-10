#!/usr/bin/env bash

set -x

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ulink_server

docker build -t hub.bunny-tech.com/prod/ulink_server:git.$1 -f Dockerfile .
docker push hub.bunny-tech.com/prod/ulink_server:git.$1

