#!/bin/bash

all: format update-mod test build

format:
	@echo "Formatting all golang file."
	@gofmt -w .

update-mod:
	export GO111MODULE=on
	@echo "Go Module Download"
	@go mod download

test:
	@echo "Testing Golang"
	@go test -race -cover ./...
	@go vet .

build:
	@echo "Building Golang"
	@go build -v -o rabbitsky .

build-release:
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=amd64
	@go build -v -a -installsuffix cgo -o rabbitsky .
