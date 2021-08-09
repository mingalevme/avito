#!/usr/bin/env sh

GOOS=darwin GOARCH=amd64 go build -o bin/avito_darwin .
GOOS=linux GOARCH=arm64 go build -o bin/avito_arm64 .
