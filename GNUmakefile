TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
default: build

build:
	go build .

test:
	 go test ./...

testacc:
	KAFKA_CA_CERT=../secrets/ca.crt \
	KAFKA_CLIENT_CERT=../secrets/client.pem \
	KAFKA_CLIENT_KEY=../secrets/client.key \
	KAFKA_CLIENT_KEY_PASSPHRASE=test-pass \
	KAFKA_SKIP_VERIFY=false \
	KAFKA_ENABLE_TLS=true \
	TF_LOG=DEBUG \
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 9m -count=1

BINARY=terraform-provider-schemaregistry
VERSION=$(shell git describe --tags --dirty)
PLATFORMS=darwin linux windows
ARCHITECTURES=amd64

wr-build-all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o $(BINARY)_v$(VERSION)-$(GOOS)-$(GOARCH))))

.PHONY: build test testacc
