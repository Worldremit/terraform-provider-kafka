FROM golang:1.13

WORKDIR /go/src/github.com/Worldremit/terraform-provider-kafka/

COPY go.mod go.sum main.go GNUmakefile ./
COPY kafka kafka
COPY secrets secrets
