GOBIN ?= $(shell go env GOPATH)/bin
PKG = github.com/dotzero/pad
BIN := pad

VERSION := 1.0.0
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
	go test -v ./service

clean:
	if [ -f $(GOBIN)/$(BIN) ] ; then rm -f $(GOBIN)/$(BIN) ; fi

fmt:
	find . -name '*.go' -not -path './.vendor/*' -exec gofmt -w=true {} ';'

vet:
	go vet ./...

.PHONY: build install test clean fmt vet
