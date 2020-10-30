#!/bin/bash

all: format update-mod test build
build-release: build-release-linux build-release-windows build-release-darwin

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

build-release-linux:
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=amd64
	@go build -v -a -installsuffix cgo -o rabbitsky .

build-release-windows:
	export CGO_ENABLED=0
	export GOOS=windows
	export GOARCH=amd64
	@go build -v -a -installsuffix cgo -o rabbitsky.exe .

build-release-darwin:
	export CGO_ENABLED=0
	export GOOS=darwin
	export GOARCH=amd64
	@go build -v -a -installsuffix cgo -o rabbitsky.exe .
