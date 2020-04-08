PLATFORMS := client/darwin/amd64 server/darwin/amd64 client/linux/amd64 server/linux/amd64 client/windows/amd64
VERSION := $(shell git describe --tags --always --dirty="-dev")
GO := go
GOFLAGS = CGO_ENABLED=0
OUTDIR := .

temp = $(subst /, ,$@)
part = $(word 1, $(temp))
os = $(word 2, $(temp))
arch = $(word 3, $(temp))

all: client/linux/amd64 server/linux/amd64
release: $(PLATFORMS)

dockerci:
	(cd ci; docker build -t zehome/sintls:latest . -f Dockerfile --pull)

version:
	@echo ${VERSION}

$(PLATFORMS):
	$(GOFLAGS) GOOS=$(os) GOARCH=$(arch) $(GO) build -ldflags='-s -w -X "main.version=${VERSION}"' -o "$(OUTDIR)/sintls-$(part)_$(os)_$(arch)" ./cmd/sintls-$(part)

.PHONY: $(PLATFORMS)
