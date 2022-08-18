# drone-manifest-ecr

TEST

Drone plugin to push Docker manifest to a AWS ECR for multi-architecture mappings. For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-manifest/).

## Build

Build the binary with the following command:

```console
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-manifest-ecr
```

## Docker

Build the Docker image with the following command:

```console
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag lemontech/drone-manifest-ecr .
```

## Usage

```console
docker run --rm \
  -e PLUGIN_PLATFORMS=linux/amd64,linux/arm,linux/arm64 \
  -e PLUGIN_TEMPLATE=organization/project-ARCH:1.0.0 \
  -e PLUGIN_TARGET=organization/project:1.0.0 \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  lemontech/drone-manifest-ecr
```
