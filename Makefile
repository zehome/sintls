GOGET=go get -u
GOFLAGS=GOOS=linux GOARCH=amd64 CGO_ENABLED=0
VERSION := $(shell git describe --tags --always --dirty="-dev")

all: bin/sintlsserver bin/sintls
dbg: bin/sintlsserver.dbg bin/sintls.dbg

version:
	$(info version: [${VERSION}])

bin/sintlsserver: server.go cli.go
	$(GOFLAGS) go build -ldflags='-s -w -X "main.version=${VERSION}"' -o bin/sintlsserver $^

bin/sintlsserver.dbg: server.go cli.go
	$(GOFLAGS) go build -ldflags='-X "main.version=${VERSION}"' -o bin/sintlsserver.dbg $^

bin/sintls: client.go
	$(GOFLAGS) go build -ldflags='-s -w -X "main.version=${VERSION}"' -o bin/sintls $^

bin/sintls.dbg: client.go
	$(GOFLAGS) go build -ldflags='-X "main.version=${VERSION}"' -o bin/sintls.dbg $^

.PHONY: bin/sintlsserver bin/sintlsserver.dbg bin/sintls bin/sintls.dbg version
