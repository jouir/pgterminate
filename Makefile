BINARY := pgterminate
GOOS := linux
GOARCH := amd64
APPVERSION := $(shell cat ./VERSION)
GOVERSION := $(shell go version | awk '{print $$3}')
GITCOMMIT := $(shell git log -1 --oneline | awk '{print $$1}')
LDFLAGS = -X main.AppVersion=${APPVERSION} -X main.GoVersion=${GOVERSION} -X main.GitCommit=${GITCOMMIT}

.PHONY: clean

build:
	go build -ldflags "${LDFLAGS}" -o bin/${BINARY} cmd/${BINARY}/main.go

release:
	go build -ldflags "${LDFLAGS}" -o bin/${BINARY} cmd/${BINARY}/main.go
	(cd bin && tar czf ${BINARY}-${APPVERSION}-${GOOS}-${GOARCH}.tar.gz ${BINARY})

test:
	go test -cover base/*
	go test -cover terminator/*

clean:
	rm -rf bin
