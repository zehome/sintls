PLATFORMS := client/darwin/amd64 server/darwin/amd64 client/linux/amd64 server/linux/amd64 client/windows/amd64
VERSION := $(shell git describe --tags --always --dirty="-dev")
GO := go
GOFLAGS = CGO_ENABLED=0

temp = $(subst /, ,$@)
part = $(word 1, $(temp))
os = $(word 2, $(temp))
arch = $(word 3, $(temp))

all: client/linux/amd64 server/linux/amd64
release: $(PLATFORMS)

dockerci:
	docker build -t zehome/sintls:latest . -f ci/Dockerfile --pull

version:
	@echo ${VERSION}

$(PLATFORMS):
	$(GOFLAGS) GOOS=$(os) GOARCH=$(arch) $(GO) build -ldflags='-s -w -X "main.version=${VERSION}"' -o 'sintls-$(part)_$(os)_$(arch)' ./cmd/sintls-$(part)

.PHONY: $(PLATFORMS)
