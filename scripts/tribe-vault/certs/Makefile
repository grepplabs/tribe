.DEFAULT_GOAL := generate

ROOT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

generate: clean docker-cfssl

ca:
	cfssl gencert -initca ca.json | cfssljson -bare ca

intermediate:
	cfssl gencert -initca intermediate-ca.json | cfssljson -bare intermediate_ca
	cfssl sign -ca ca.pem -ca-key ca-key.pem -config cfssl.json -profile intermediate_ca intermediate_ca.csr | cfssljson -bare intermediate_ca

host:
	cfssl gencert -ca intermediate_ca.pem -ca-key intermediate_ca-key.pem -config cfssl.json -profile=server vault.json | cfssljson -bare vault-server

bundle:
	mkbundle -f tribe-ca-bundle.pem intermediate_ca.pem ca.pem

clean:
	rm -f *.pem
	rm -f *.csr
	rm -f *.crt

CFSSL_VERSION := 1.5.0
CFSSL := docker run -u $(shell id -u):$(shell id -g) --rm -v $(CURDIR):$(CURDIR) -w $(CURDIR) --entrypoint /bin/bash cfssl/cfssl:$(CFSSL_VERSION)
docker-cfssl:
	# ca
	$(CFSSL) -c 'cfssl gencert -initca ca.json | cfssljson -bare ca'
	# intermediate
	$(CFSSL) -c 'cfssl gencert -initca intermediate-ca.json | cfssljson -bare intermediate_ca'
	$(CFSSL) -c 'cfssl sign -ca ca.pem -ca-key ca-key.pem -config cfssl.json -profile intermediate_ca intermediate_ca.csr | cfssljson -bare intermediate_ca'
	# host
	$(CFSSL) -c 'cfssl gencert -ca intermediate_ca.pem -ca-key intermediate_ca-key.pem -config cfssl.json -profile=server vault.json | cfssljson -bare vault-server'
	# bundle
	$(CFSSL) -c 'mkbundle -f tribe-ca-bundle.pem intermediate_ca.pem ca.pem'
	# read access
	chmod +r *.pem
