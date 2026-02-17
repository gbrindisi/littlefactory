.PHONY: build test install

PREFIX ?= /usr/local

build:
	go build -o bin/littlefactory ./cmd/littlefactory

install: build
	install -d $(PREFIX)/bin
	install -m 755 bin/littlefactory $(PREFIX)/bin/littlefactory

test:
	go test ./...
