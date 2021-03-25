.DEFAULT_GOAL := help

.PHONY: clean build fmt test

ROOT_DIR      := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

BUILD_FLAGS   ?=
BRANCH        = $(shell git rev-parse --abbrev-ref HEAD)
REVISION      = $(shell git describe --tags --always --dirty)
BUILD_DATE    = $(shell date +'%Y.%m.%d-%H:%M:%S')
LDFLAGS       ?= -w -s

BINARY        = tribe
LOCAL_IMAGE   ?= local/$(BINARY)

default: help

.PHONY: help
help:
	@grep -E '^[a-zA-Z%_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

check: ## Vet and lint
	go vet ./...
	golint $$(go list ./...) 2>&1
	gosec ./... 2>&1

test: ## Test
	GO111MODULE=on go test -mod=vendor -v ./...

build: ## Build executables
	CGO_ENABLED=0 GO111MODULE=on go build -mod=vendor -o $(BINARY) $(BUILD_FLAGS) -ldflags "$(LDFLAGS)" .

fmt: ## Go format
	go fmt ./...

clean: ## Clean
	@rm -rf $(BINARY)

.PHONY: deps
deps: ## Get dependencies
	GO111MODULE=on go get ./...

.PHONY: vendor
vendor: ## Go vendor
	GO111MODULE=on go mod vendor

.PHONY: tidy
tidy: ## Go tidy
	GO111MODULE=on go mod tidy

SWAGGER_VERSION := v0.26.1
SWAGGER := docker run -u $(shell id -u):$(shell id -g) --rm -v $(CURDIR):$(CURDIR) -w $(CURDIR) -e GOCACHE=/tmp/.cache --entrypoint swagger quay.io/goswagger/swagger:$(SWAGGER_VERSION)

postgres-up: ## Start test postgres
	cd $(ROOT_DIR)/scripts/tribe-postgres && docker-compose up

# https://github.com/golang-migrate/migrate
MIGRATIONS_PATH := $(ROOT_DIR)/database/migrations/migrate
migrate: ## Database migration using migrate tool
	migrate -path $(MIGRATIONS_PATH)/postgres -database postgres://tribe:secret@localhost:5432/tribe?sslmode=disable up

LIQUIBASE_PATH  := $(ROOT_DIR)/database/migrations/liquibase
liquibase: ## Database migration using liquibase
	docker run --network host --rm -u $(shell id -u):$(shell id -g) \
		-v $(LIQUIBASE_PATH):/liquibase/changelog liquibase/liquibase:4.2 \
		--driver=org.postgresql.Driver --url="jdbc:postgresql://localhost:5432/tribe" \
		--changeLogFile=changelog/changelog.yaml --username=tribe --password=secret \
		update
