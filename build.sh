#!/bin/bash
BINARY=pgterminate
VERSION=$(cat VERSION)
BUILD_PATH=/tmp/${BINARY}-${VERSION}
ldflags="-X main.AppVersion=${VERSION}"
GOOS=linux
GOARCH=amd64

export GOOS
export GOARCH

go build -ldflags "$ldflags" -o ${BUILD_PATH}/${BINARY} cmd/${BINARY}/main.go
(cd ${BUILD_PATH} && tar czf ${BINARY}-${VERSION}-${GOOS}-${GOARCH}.tar.gz ${BINARY})

echo "Archive created:"
ls -l ${BUILD_PATH}/${BINARY}-${VERSION}-${GOOS}-${GOARCH}.tar.gz