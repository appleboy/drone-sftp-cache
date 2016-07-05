# drone-sftp-cache

Drone plugin for caching artifacts to a central server using rsync

## Build

Build the binary with the following commands:

```
export GO15VENDOREXPERIMENT=1
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

go build -a -tags netgo
```

## Docker

Build the docker image with the following commands:

```
docker build --rm=true -t plugins/sftp-cache .
```

Please note incorrectly building the image for the correct x64 linux and with CGO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-sftp-cache' not found or does not exist..
```

## Usage

Build and publish from your current working directory:

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
  -e SFTP_CACHE_PRIVATE_KEY=$(cat ~/.ssh/id_rsa) \
  plugins/sftp-cache
```
