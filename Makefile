.DEFAULT_GOAL := build

.PHONY: clean build fmt test

ROOT_DIR      := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

BUILD_FLAGS   ?=
BRANCH        = $(shell git rev-parse --abbrev-ref HEAD)
REVISION      = $(shell git describe --tags --always --dirty)
BUILD_DATE    = $(shell date +'%Y.%m.%d-%H:%M:%S')
LDFLAGS       ?= -w -s

BINARY        = tribe
LOCAL_IMAGE   ?= local/$(BINARY)

default: build

check:
	go vet ./...
	golint $$(go list ./...) 2>&1
	gosec ./... 2>&1

test:
	GO111MODULE=on go test -mod=vendor -v ./...

build:
	CGO_ENABLED=0 GO111MODULE=on go build -mod=vendor -o $(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" .


fmt:
	go fmt ./...

clean:
	@rm -rf $(BINARY)

.PHONY: deps
deps:
	GO111MODULE=on go get ./...

.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor

.PHONY: tidy
tidy:
	GO111MODULE=on go mod tidy

SWAGGER_VERSION := v0.26.1
SWAGGER := docker run -u $(shell id -u):$(shell id -g) --rm -v $(CURDIR):$(CURDIR) -w $(CURDIR) -e GOCACHE=/tmp/.cache --entrypoint swagger quay.io/goswagger/swagger:$(SWAGGER_VERSION)
