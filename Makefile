GOGET=go get -u

all: bin/sintlsserver

bin/sintlsserver: server.go cli.go
	go build -o bin/sintlsserver $^

deps:
	$(GOGET) github.com/gin-gonic/gin
	$(GOGET) github.com/jinzhu/gorm

.PHONY: bin/sintlsserver
