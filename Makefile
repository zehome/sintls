GOFLAGS=GOOS=linux GOARCH=amd64 CGO_ENABLED=0
VERSION := $(shell git describe --tags --always --dirty="-dev")
GO=go

all: cmd/sintls-server cmd/sintls-client

dockerci:
	docker build -t registry.clarisys.fr/adm/sintls . -f ci/Dockerfile --pull

version:
	@echo ${VERSION}

cmd/sintls-server:
	$(GOFLAGS) $(GO) build -ldflags='-s -w -X "main.version=${VERSION}"' ./$@
cmd/sintls-client:
	$(GOFLAGS) $(GO) build -ldflags='-s -w -X "main.version=${VERSION}"' ./$@

.PHONY: cmd/sintls-server cmd/sintls-client version
