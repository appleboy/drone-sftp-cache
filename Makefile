.PHONY: build

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o drone-sftp-cache .

docker:
	docker build --rm=true -t plugins/drone-sftp-cache .

docker-test:
	docker build --rm=true -t plugins/drone-sftp-cache:test .
