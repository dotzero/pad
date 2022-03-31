GOBIN ?= $(shell go env GOPATH)/bin
BIN := pad

VERSION := $(shell git rev-parse --abbrev-ref HEAD)
HASH := $(shell git rev-parse --short HEAD)
DATE := $(shell date +%FT%T%z)

LDFLAGS := "-s -w \
	-X main.Version=$(VERSION) \
	-X main.CommitHash=$(HASH) \
	-X main.CompileDate=$(DATE)"

all: build

build: fmt vet
	go build -ldflags=$(LDFLAGS) -o $(GOBIN)/$(BIN)

install:
	go install -ldflags=$(LDFLAGS)

test:
	go test ./...

clean:
	if [ -f $(GOBIN)/$(BIN) ] ; then rm -f $(GOBIN)/$(BIN) ; fi

fmt:
	find . -name '*.go' -not -path './.vendor/*' -exec gofmt -w=true {} ';'

vet:
	go vet ./...

vendor:
	go mod tidy
	go mod vendor

.PHONY: build install test clean fmt vet vendor
