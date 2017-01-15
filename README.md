# drone-sftp-cache

[![GoDoc](https://godoc.org/github.com/appleboy/drone-sftp-cache?status.svg)](https://godoc.org/github.com/appleboy/drone-sftp-cache) [![Build Status](http://drone.wu-boy.com/api/badges/appleboy/drone-sftp-cache/status.svg)](http://drone.wu-boy.com/appleboy/drone-sftp-cache) [![codecov](https://codecov.io/gh/appleboy/drone-sftp-cache/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-sftp-cache) [![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-sftp-cache)](https://goreportcard.com/report/github.com/appleboy/drone-sftp-cache) [![Docker Pulls](https://img.shields.io/docker/pulls/appleboy/drone-sftp-cache.svg)](https://hub.docker.com/r/appleboy/drone-sftp-cache/) [![](https://images.microbadger.com/badges/image/appleboy/drone-sftp-cache.svg)](https://microbadger.com/images/appleboy/drone-sftp-cache "Get your own image badge on microbadger.com")

This project is forked from [drone-plugins/drone-sftp-cache](https://github.com/drone-plugins/drone-sftp-cache).

Drone plugin for caching artifacts to a central server using rsync. For the
usage information and a listing of the available options please take a look at
[the docs](DOCS.md).

## Build

Build the binary with the following commands:

```
go build
go test
```

## Docker

Build the docker image with the following commands:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build --rm=true -t appleboy/drone-sftp-cache .
```

Please note incorrectly building the image for the correct x64 linux and with
GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-sftp-cache' not found or does not exist..
```

## Usage

Execute from the working directory:

```
docker run --rm \
  -e DRONE_REPO=octocat/hello-world \
  -e DRONE_REPO_BRANCH=master \
  -e DRONE_COMMIT_BRANCH=master \
  -e PLUGIN_MOUNT=node_modules \
  -e PLUGIN_RESTORE=false \
  -e PLUGIN_REBUILD=true \
  -e PLUGIN_IGNORE_BRANCH=false \
  -e SFTP_CACHE_SERVER=1.2.3.4 \
  -e SFTP_CACHE_PORT=22 \
  -e SFTP_CACHE_PATH=/root/cache \
  -e SFTP_CACHE_USERNAME=root \
  -e SFTP_CACHE_PRIVATE_KEY=$(cat ~/.ssh/id_rsa) \
  appleboy/drone-sftp-cache
```
