GOGET=go get -u
GOFLAGS=GOOS=linux GOARCH=amd64 CGO_ENABLED=0
VERSION := $(shell git describe --tags --always --dirty="-dev")

all: cmd/sintls-server cmd/sintls-client

version:
	@echo ${VERSION}

cmd/sintls-server:
	$(GOFLAGS) go build -ldflags='-s -w -X "main.version=${VERSION}"' ./$@
cmd/sintls-client:
	$(GOFLAGS) go build -ldflags='-s -w -X "main.version=${VERSION}"' ./$@

.PHONY: cmd/sintls-server cmd/sintls-client version
