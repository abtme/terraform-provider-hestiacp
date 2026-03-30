default: build

BINARY    := terraform-provider-hestiacp
VERSION   ?= 0.1.0
OS_ARCH   := $(shell go env GOOS)_$(shell go env GOARCH)
INSTALL_PATH := ~/.terraform.d/plugins/registry.terraform.io/abtme/hestiacp/$(VERSION)/$(OS_ARCH)

.PHONY: build install fmt lint test testacc docs clean

build:
	go build -o $(BINARY) .

install: build
	mkdir -p $(INSTALL_PATH)
	cp $(BINARY) $(INSTALL_PATH)/$(BINARY)

fmt:
	gofmt -s -w .
	goimports -w .

lint:
	golangci-lint run ./...

test:
	go test ./... -v -timeout 120s

testacc:
	TF_ACC=1 go test ./internal/provider/... -v -timeout 300s

docs:
	go generate ./...

clean:
	rm -f $(BINARY)
