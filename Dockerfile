# Docker image for Drone's slack notification plugin
#
#     docker build --rm=true -t plugins/drone-sftp-cache .

FROM alpine:3.2
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ADD drone-sftp-cache /bin/
ENTRYPOINT ["/bin/drone-sftp-cache"]
