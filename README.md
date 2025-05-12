# drone-manifest

[![Build Status](http://harness.drone.io/api/badges/drone-plugins/drone-manifest/status.svg)](http://harness.drone.io/drone-plugins/drone-manifest)
[![Slack](https://img.shields.io/badge/slack-drone-orange.svg?logo=slack)](https://join.slack.com/t/harnesscommunity/shared_invite/zt-y4hdqh7p-RVuEQyIl5Hcx4Ck8VCvzBw)
[![Join the discussion at https://community.harness.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://community.harness.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-manifest?status.svg)](http://godoc.org/github.com/drone-plugins/drone-manifest)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-manifest)](https://goreportcard.com/report/github.com/drone-plugins/drone-manifest)

Drone plugin to push Docker manifest to a registry for multi-architecture mappings. For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-manifest/).

## Build

Build the binary with the following command:

```console
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

go build -v -a -tags netgo -o release/linux/amd64/drone-manifest
go build -v -a -tags netgo -o release/linux/amd64/manifest-ecr ./cmd/manifest-ecr
```

## Docker

Build the Docker image with the following command:

```console
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag plugins/manifest .
  
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/ecr/Dockerfile.linux.amd64 --tag plugins/manifest-ecr .
```

## Usage

```console
docker run --rm \
  -e PLUGIN_PLATFORMS=linux/amd64,linux/arm,linux/arm64 \
  -e PLUGIN_TEMPLATE=organization/project-ARCH:1.0.0 \
  -e PLUGIN_TARGET=organization/project:1.0.0 \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/manifest
```
