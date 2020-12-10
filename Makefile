.DEFAULT_GOAL := build

.PHONY: clean build fmt test

BUILD_FLAGS   ?=
BINARY        ?= tribe
BRANCH        = $(shell git rev-parse --abbrev-ref HEAD)
REVISION      = $(shell git describe --tags --always --dirty)
BUILD_DATE    = $(shell date +'%Y.%m.%d-%H:%M:%S')
LDFLAGS       ?= -w -s

ROOT_DIR      := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

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

SWAGGER_VERSION := v0.25.0
SWAGGER := docker run -u $(shell id -u):$(shell id -g) --rm -v $(CURDIR):$(CURDIR) -w $(CURDIR) --entrypoint swagger quay.io/goswagger/swagger:$(SWAGGER_VERSION)

delete-api:
	@rm -rf api/v1/client/
	@rm -rf api/v1/models/
	@rm -rf api/v1/server/
	@rm -rf api/v1/cmd/


generate-server-api: api/v1/openapi.yaml
	@echo GEN api/v1/openapi.yaml
	$(SWAGGER) generate server -s server -a restapi \
			-t api/v1 \
			-f api/v1/openapi.yaml \
			--exclude-main \
			--default-scheme=http \
			-C api/v1/server.yml
	@# sort goimports automatically
	@find api/v1/models/ -type f -name "*.go" -print | PATH="$(PWD)/tools:$(PATH)" xargs goimports -w
	@find api/v1/server/ -type f -name "*.go" -print | PATH="$(PWD)/tools:$(PATH)" xargs goimports -w


generate-client-api: api/v1/openapi.yaml
	@echo GEN api/v1/openapi.yaml
	$(SWAGGER) generate client -a restapi \
		-t api/v1 \
		-f api/v1/openapi.yaml
	@# sort goimports automatically
	@find api/v1/client/ -type f -name "*.go" -print | PATH="$(PWD)/tools:$(PATH)" xargs goimports -w
	@find api/v1/models/ -type f -name "*.go" -print | PATH="$(PWD)/tools:$(PATH)" xargs goimports -w


TEST_REALM=main
test-realm:
	curl -X POST -H 'Content-Type: application/json' -d '{"realm_id": "$(TEST_REALM)"}' http://localhost:8080/v1/realms

test-users:
	@for i in `seq 1000`; do \
	echo "{\"username\": \"user-"`uuidgen`"-$${i}\", \"password\": \"hello\"}" |  curl -X POST -H 'Content-Type: application/json' -d @- http://localhost:8080/v1/realms/$(TEST_REALM)/users ; \
	done
