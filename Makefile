
###############################################################################
# 获得子版本号 go version go1.23.3 darwin/amd64 -> 23
GOVERSION := $(shell go version | cut -d ' ' -f 3 | cut -d '.' -f 2)
# 获得操作系统类型，并转换为小写
OS = $(shell uname | tr A-Z a-z)
# 设置性能分析参数
# ?= 是一个条件赋值操作符，用于设置变量的默认值。如果变量已经被赋值，则不会改变其值；如果变量未被赋值，则会使用 ?= 后面的值作为默认值。
# -cpuprofile=cpu.pprof 启用 CPU 性能分析，并将结果输出到名为 cpu.pprof 的文件中。性能分析可以帮助你找出代码中的 CPU 使用瓶颈，从而优化程序性能。
# -memprofile=mem.pprof 启用内存分配性能分析，并将结果输出到名为 mem.pprof 的文件中。内存分析可以帮助你找出代码中的内存分配和泄漏问题，从而优化程序的内存使用。	
# -benchmem 此标志用于在运行基准测试时显示内存分配统计信息。它会报告每次操作分配的字节数以及总分配的字节数。这有助于了解基准测试中的内存使用情况。
BENCH_FLAGS ?= -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem

###############################################################################
# 设置环境变量
# 用于控制 Go 语言中 CGO（C-Go）的启用状态。CGO 允许 Go 代码调用 C 语言库，从而实现与现有 C 代码库的互操作性
export CGO_ENABLED ?= 0
export GOOS = $(shell go env GOOS)

###############################################################################
# 设置版本
ifdef VERSION
    VERSION=${VERSION}
else
    VERSION=$(shell git describe --always 2>/dev/null || echo "--")
endif
GITCOMMIT=$(shell git rev-parse HEAD 2>/dev/null || echo "--")
BUILDTIME=${shell date +%Y-%m-%dT%H:%M:%S%z}
LDFLAGS="-X github.com/fengzhongzhu1621/xgo/version.AppVersion=${VERSION} \
	-X github.com/fengzhongzhu1621/xgo/version.GitCommit=${GITCOMMIT} \
	-X github.com/fengzhongzhu1621/xgo/version.BuildTime=${BUILDTIME}"

###############################################################################
# Build variables
BUILD_DIR ?= build

ifeq (${VERBOSE}, 1)
# 从 GOARGS 变量中筛选出包含 -v 的部分
ifeq ($(filter -v,${GOARGS}),)
# 如果筛选结果为空
	GOARGS += -v
endif
TEST_FORMAT = short-verbose
endif

TEST_PKGS ?= ./...

###############################################################################
# Dependency versions
# golangci-lint
GOTESTSUM_VERSION = 1.12.0
GOLANGCI_VERSION = 1.62.0
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)
.PHONY: golines swag gofumpt

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

SWAG ?= $(LOCALBIN)/swag
GOLINES ?= $(LOCALBIN)/golines
GOFUMPT ?= $(LOCALBIN)/gofumpt

swag: $(SWAG) ## install swagger
$(SWAG): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/swaggo/swag/cmd/swag@v1.16.3

golines: $(GOLINES) ## install golines
$(GOLINES): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/segmentio/golines@v0.12.2

gofumpt: $(GOFUMPT) ## install gofumpt
$(GOFUMPT): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install mvdan.cc/gofumpt@v0.6.0

gowatch:
	go install github.com/silenceper/gowatch@latest


###############################################################################
# 自定义命令
.PHONY: build test bench vet coverage check test-race test-cover-html help clear lint fix fmt tidy
# 运行 make 命令而没有指定目标时，默认会执行 help 目标，并打印出帮助信息。
.DEFAULT_GOAL := help
help:	# Empty target rule

check: test-race vet

build:
	go build -v ./...

# gotestsum 是一个用于格式化和汇总 Go 测试输出的工具。它可以将 go test 的输出转换为更易读的格式，并且可以生成 JSON、JUnit XML 等格式的报告。
test: TEST_FORMAT ?= short
test: SHELL = /bin/bash
test: export CGO_ENABLED=1
test: go mod tidy
test: bin/gotestsum ## Run tests
	@mkdir -p ${BUILD_DIR}
	bin/gotestsum --no-summary=skipped --junitfile ${BUILD_DIR}/coverage.xml --format ${TEST_FORMAT} -- -coverprofile=${BUILD_DIR}/coverage.txt -covermode=atomic $(filter-out -v,${GOARGS}) $(if ${TEST_PKGS},${TEST_PKGS},./...)

# test:
#   go mod tidy
#	go test -v ./... -cover

# test:
#	go test -mod=vendor -gcflags=all=-l $(shell go list ./... | grep -v mock | grep -v docs) -covermode=count -coverprofile .coverage.cov
#	go tool cover -func=.coverage.cov

# test:
# 	go test ./...

# test_v:
# 	go test -v ./...

# test_short:
# 	go test ./... -short

# test-race:
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
# 	go test -run=nonthingplease -benchmem -bench=. $(shell go list ./... | grep -v /vendor/)
#	go test -v ./... -test.bench -test.benchmem
	go list ./... | xargs -n1 go test -bench=$(BENCH) -run="^$$" $(BENCH_FLAGS)

vet:
	go vet

#test-cover-html:
#	go test -v ./... -coverprofile=coverage.out -covermode=count
#	go tool cover -func=coverage.out

list: ## List all make targets
	@${MAKE} -pRrn : -f $(MAKEFILE_LIST) 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | sort

fmt: golines gofumpt ## format source code
	go fmt ./...
	goimports -l -w .
	$(GOLINES) ./ -m 119 -w --base-formatter gofmt --no-reformat-tags
	$(GOFUMPT) -l -w .

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

# 安装 go install golang.org/x/tools/cmd/godoc
godoc: ## show doc
	echo "http://127.0.0.1:6060"
	godoc -http=127.0.0.1:6060 -goroot="."

docker-build: ## build docker image
#	docker build -f ./Dockerfile -t xgo:${VERSION} .
	docker build -f ./Dockerfile -t xgo:latest .