.PHONY: update fmt build test
.DEFAULT_GOAL := build

update:
	go get -u

fmt:
	go fmt ./...

build: fmt
	go build -ldflags="-s -w"

test: fmt
	go test

