# drone-sftp-cache

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-sftp-cache/status.svg)](http://beta.drone.io/drone-plugins/drone-sftp-cache)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-sftp-cache?status.svg)](http://godoc.org/github.com/drone-plugins/drone-sftp-cache)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-sftp-cache)](https://goreportcard.com/report/github.com/drone-plugins/drone-sftp-cache)
[![Join the chat at https://gitter.im/drone/drone](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/drone/drone)

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
docker build --rm=true -t plugins/sftp-cache .
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
  -e SFTP_CACHE_SERVER=1.2.3.4:22 \
  -e SFTP_CACHE_PATH=/root/cache \
  -e SFTP_CACHE_USERNAME=root \
  -e SFTP_CACHE_PRIVATE_KEY="$(cat ~/.ssh/id_rsa)" \
  plugins/sftp-cache
```
