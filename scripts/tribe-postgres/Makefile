.DEFAULT_GOAL := up

ROOT_DIR        := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))
MIGRATIONS_PATH := $(ROOT_DIR)/../../database/sql/migrations

up:
	docker-compose up

up-d:
	docker-compose up -d

down:
	docker-compose down

# https://github.com/golang-migrate/migrate
migrate:
	migrate -path $(MIGRATIONS_PATH)/postgres -database postgres://tribe:secret@localhost:5432/tribe?sslmode=disable up
