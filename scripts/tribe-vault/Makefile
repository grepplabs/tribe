.DEFAULT_GOAL := up

ROOT_DIR 	:= $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

generate-certs:
	bash -c "cd $(ROOT_DIR)/certs && make generate"

up:
	docker-compose up

up-d:
	docker-compose up -d

down:
	docker-compose down

vault-login:
	 export VAULT_ADDR='https://0.0.0.0:8201' && vault login -ca-cert=$(ROOT_DIR)/certs/tribe-ca-bundle.pem token=tribe-root-token

