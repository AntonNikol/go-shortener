# syntax=docker/dockerfile:1

ARG UBUNTU_IMAGE_TAG="22.04"

# base image
FROM ubuntu:$UBUNTU_IMAGE_TAG as base
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update \
    && apt-get install --no-install-recommends --autoremove --purge -y curl tar gzip gccgo ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# prepare base image with golang
FROM base as go
ARG TARGETOS
ARG TARGETARCH
ARG GO_VERSION="1.19.5"
ENV GOPATH=/go
ENV PATH=/usr/local/bin/go/bin:$PATH
WORKDIR /usr/local/bin/go
RUN curl -L "https://go.dev/dl/go${GO_VERSION}.${TARGETOS}-${TARGETARCH}.tar.gz" | tar --strip=1 -xzf - \
    && rm -rf /var/lib/apt/lists/*

# compile all necesssarlily autotests for both tracks from sources \
# executables: devopstest, shortenertest (+goimports lib), statictest, random
# (this step is need to support linux-arm64 docker hosts)
FROM go as go-autotests
ARG GO_AUTOTESTS_VERSION="v0.7.10"
WORKDIR /opt/devops_autotests
RUN curl -L "https://github.com/Yandex-Practicum/go-autotests/archive/refs/tags/$GO_AUTOTESTS_VERSION.tar.gz" \
    | tar --strip=1 -xzf -
RUN go test -c -o=/usr/local/bin/devopstest ./cmd/devopstest/... \
    && go test -c -o=/usr/local/bin/shortenertest ./cmd/shortenertest/... \
    && go build -o=/usr/local/bin/statictest ./cmd/statictest/... \
    && go build -o=/usr/local/bin/random ./cmd/random/... \
    && go install golang.org/x/tools/cmd/goimports@latest

# cache project dependencies
FROM go as project-go-mod-cache
COPY go.* .
RUN go mod download

# create base image for all tracks (devops, shortener)
FROM base as track
ENV GOPATH=/go
ENV PATH=/usr/local/bin/go/bin:$GOPATH/bin:$PATH
COPY --from=go /usr/local/bin/go /usr/local/bin/go
COPY --from=go-autotests /usr/local/bin/statictest /usr/local/bin/statictest
COPY --from=go-autotests /usr/local/bin/random /usr/local/bin/random
COPY --from=project-go-mod-cache $GOPATH $GOPATH
COPY . /app
WORKDIR /app

# image for devops-track
FROM track as devops-track
ARG YA_AGENT_BINARY_PATH
ARG YA_SERVER_BINARY_PATH
COPY --from=go-autotests /usr/local/bin/devopstest /usr/local/bin/devopstest
RUN go build -o=$YA_AGENT_BINARY_PATH ./cmd/agent/... \
    && go build -o=$YA_SERVER_BINARY_PATH ./cmd/server/...

# image for shortener-track
FROM track as shortener-track
ARG YA_BINARY_PATH
COPY --from=go-autotests /usr/local/bin/shortenertest /usr/local/bin/shortenertest
RUN go build -o=$YA_BINARY_PATH ./cmd/shortener/...