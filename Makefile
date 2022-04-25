GOVERSION := $(shell go version | cut -d ' ' -f 3 | cut -d '.' -f 2)

.PHONY: build test bench vet coverage
.DEFAULT_GOAL := help
help:	# Empty target rule

check: test-race vet

build:
	go build -v ./...
test:
	go test -v ./... -cover

test-race:
	go test -v ./... -race

bench:
	go test -v ./... -test.bench -test.benchmem
vet:
	go vet

test-cover-html:
	go test -v ./... -coverprofile=coverage.out -covermode=count
	go tool cover -func=coverage.out

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
