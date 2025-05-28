.PHONY: build test clean

VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/git-replicator main.go
	chmod +x bin/git-replicator

test:
	go test -v ./...

clean:
	rm -rf bin/ 