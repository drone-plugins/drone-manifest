# drone-manifest

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-manifest/status.svg)](http://beta.drone.io/drone-plugins/drone-manifest)
[![Join the discussion at https://discourse.drone.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://discourse.drone.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-manifest?status.svg)](http://godoc.org/github.com/drone-plugins/drone-manifest)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-manifest)](https://goreportcard.com/report/github.com/drone-plugins/drone-manifest)
[![](https://images.microbadger.com/badges/image/plugins/manifest.svg)](https://microbadger.com/images/plugins/manifest "Get your own image badge on microbadger.com")

Drone plugin to push Docker manifest to a registry for multi-architecture mappings. For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-manifest/).

## Build

Build the binary with the following commands:

```
go build
```

## Docker

Build the Docker image with the following commands:

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -o release/linux/amd64/drone-manifest
docker build --rm -t plugins/manifest .
```

### Usage

```
docker run --rm \
  -e PLUGIN_PLATFORMS=linux/amd64,linux/arm,linux/arm64 \
  -e PLUGIN_TEMPLATE=organization/project-ARCH:1.0.0 \
  -e PLUGIN_TARGET=organization/project:1.0.0 \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/manifest
```
