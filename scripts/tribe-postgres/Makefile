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

TEST_REALM=main
test-realm:
	curl -X POST -H 'Content-Type: application/json' -d '{"realm_id": "$(TEST_REALM)"}' http://localhost:8080/v1/realms

test-users:
	@for i in `seq 1000`; do \
	echo "{\"username\": \"user-"`uuidgen`"-$${i}\", \"password\": \"hello\"}" |  curl -X POST -H 'Content-Type: application/json' -d @- http://localhost:8080/v1/realms/$(TEST_REALM)/users ; \
	done