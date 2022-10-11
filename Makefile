GOVERSION := $(shell go version | cut -d ' ' -f 3 | cut -d '.' -f 2)
OS = $(shell uname | tr A-Z a-z)
BENCH_FLAGS ?= -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem

# Build variables
BUILD_DIR ?= build
export CGO_ENABLED ?= 0
export GOOS = $(shell go env GOOS)
ifeq (${VERBOSE}, 1)
ifeq ($(filter -v,${GOARGS}),)
	GOARGS += -v
endif
TEST_FORMAT = short-verbose
endif

TEST_PKGS ?= ./...

# Dependency versions
GOTESTSUM_VERSION = 1.8.0
GOLANGCI_VERSION = 1.45.2

bin/gotestsum-${GOTESTSUM_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_${OS}_amd64.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum-${GOTESTSUM_VERSION} && chmod +x ./bin/gotestsum-${GOTESTSUM_VERSION}

bin/gotestsum: bin/gotestsum-${GOTESTSUM_VERSION}
	@ln -sf gotestsum-${GOTESTSUM_VERSION} bin/gotestsum

bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b ./bin/ v${GOLANGCI_VERSION}
	@mv bin/golangci-lint "$@"

bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint

.PHONY: build test bench vet coverage check test-race test-cover-html help clear lint fix fmt
.DEFAULT_GOAL := help
help:	# Empty target rule

check: test-race vet

build:
	go build -v ./...

#test:
#	go test -v ./... -cover
test: TEST_FORMAT ?= short
test: SHELL = /bin/bash
test: export CGO_ENABLED=1
test: bin/gotestsum ## Run tests
	@mkdir -p ${BUILD_DIR}
	bin/gotestsum --no-summary=skipped --junitfile ${BUILD_DIR}/coverage.xml --format ${TEST_FORMAT} -- -coverprofile=${BUILD_DIR}/coverage.txt -covermode=atomic $(filter-out -v,${GOARGS}) $(if ${TEST_PKGS},${TEST_PKGS},./...)

# test:
# 	go test ./...

# test_v:
# 	go test -v ./...

# test_short:
# 	go test ./... -short

#test-race:
#	go test -v ./... -race
test-race: bin/gotestsum ## Run tests with race
	@mkdir -p ${BUILD_DIR}
	bin/gotestsum --no-summary=skipped --junitfile ${BUILD_DIR}/coverage.xml --format ${TEST_FORMAT} -- -race -coverprofile=${BUILD_DIR}/coverage.txt -covermode=atomic $(filter-out -v,${GOARGS}) $(if ${TEST_PKGS},${TEST_PKGS},./...)

lint: bin/golangci-lint ## Run linter
	bin/golangci-lint run

fix: bin/golangci-lint ## Fix lint violations
	bin/golangci-lint run --fix

BENCH ?= .
bench:
#	go test -v ./... -test.bench -test.benchmem
	go list ./... | xargs -n1 go test -bench=$(BENCH) -run="^$$" $(BENCH_FLAGS)

vet:
	go vet

#test-cover-html:
#	go test -v ./... -coverprofile=coverage.out -covermode=count
#	go tool cover -func=coverage.out

list: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

fmt:
	go fmt ./...
	goimports -l -w .

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

clear: ## Clear the working area and the project
	rm -rf bin/

generate_gomod:
	rm go.mod go.sum || true
	go mod init github.com/fengzhongzhu1621/xgo

	go install ./...
	sed -i '\|go |d' go.mod
	go mod edit -fmt

validate_examples:
	go run dev/update-examples-deps/main.go
	go run dev/validate-examples/main.go
