.PHONY: build test clean

VERSION := $(shell git describe --tags --always --dirty)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/git-replicator main.go
	chmod +x bin/git-replicator

lint:
	golangci-lint run

test:
	go test -v ./...

clean:
	rm -rf bin/

deploy: build
	sudo cp bin/git-replicator /usr/local/bin/.
