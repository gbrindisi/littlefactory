.PHONY: build test test-integration install setup

PREFIX ?= /usr/local

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -s -w \
	-X main.version=$(VERSION) \
	-X main.commit=$(COMMIT)

build:
	go build -ldflags '$(LDFLAGS)' -o bin/littlefactory ./cmd/littlefactory

install: build
	install -d $(PREFIX)/bin
	install -m 755 bin/littlefactory $(PREFIX)/bin/littlefactory

test:
	go test ./...

test-integration: build
	go test -tags=integration -v ./cmd/littlefactory/

setup:
	go install github.com/securego/gosec/v2/cmd/gosec@latest
