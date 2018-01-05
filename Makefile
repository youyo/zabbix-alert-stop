.DEFAULT_GOAL := help

## Setup
setup:
	go get -u -v github.com/golang/dep/cmd/dep

## Install dependencies
deps:
	dep ensure -v

## Vet
vet:
	go tool vet -v *.go

## Lint
lint:
	golint -set_exit_status *.go

## Run tests
test:
	go test -v -cover

## Execute `go run`
run:
	go run *.go

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps vet lint test run help
