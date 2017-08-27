GO?=go
BINARY_NAME?=pad
BINARY_PATH?=$(GOPATH)/bin/$(BINARY_NAME)

VERSION?=1.0.0
HASH:=$(shell git rev-parse --short HEAD)
DATE:=$(shell date +%FT%T%z)
LDFLAGS:="-s -w \
	-X main.Version=$(VERSION) \
	-X main.CommitHash=$(HASH) \
	-X main.CompileDate=$(DATE)"

all: build

build: fmt
	$(GO) build -ldflags=$(LDFLAGS) -o $(BINARY_PATH)

install:
	$(GO) install -ldflags=$(LDFLAGS)

test:
	$(GO) test -v $(shell $(GO) list ./... | grep -v /vendor/)

clean:
	$(GO) clean
	if [ -f $(GOBIN)/$(BIN) ] ; then rm -f $(GOBIN)/$(BIN) ; fi

fmt:
	find . -name '*.go' -not -path './vendor/*' -exec gofmt -w=true {} ';'

vet:
	$(GO) vet ./...

.PHONY: build install test clean fmt vet
