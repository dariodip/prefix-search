PKGS := $(shell go list ./... | grep -v /vendor)

.PHONY: test
test: install
	go test $(PKGS)

.PHONY: install
install:
	go get -t ./...

BIN_DIR := $(GOPATH)/bin

BINARY := prefix-search
VERSION ?= vlatest

.PHONY: build
build: install test
	go build -o $(BIN_DIR)/$(BINARY)

PLATFORMS := windows linux darwin
os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$(os)-amd64

.PHONY: release
release: windows linux darwin

