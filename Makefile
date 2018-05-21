
PKGS := $(shell go list ./... | grep -v /vendor)

.PHONY: test
test: install
	go test $(PKGS)

.PHONY: install
install:
	go get -t ./...

BIN_DIR := $(GOPATH)/bin

BINARY := prefsearch
VERSION ?= vlatest

release: test
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/$(BINARY)-v1.0.0
