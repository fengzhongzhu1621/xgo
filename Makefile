
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
GOTESTSUM_VERSION = 1.12.0
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)


###############################################################################
# 安装工具
.PHONY: install_golangcui-lint golines_tool swag_tool gofumpt_tool subfinder_tool gowatch_tool wire_tool nilaway_tool govulncheck_tool install_oapi-codegen tools

bin/gotestsum:
	@mkdir -p bin
	curl -L https://github.com/gotestyourself/gotestsum/releases/download/v${GOTESTSUM_VERSION}/gotestsum_${GOTESTSUM_VERSION}_${OS}_amd64.tar.gz | tar -zOxf - gotestsum > ./bin/gotestsum-${GOTESTSUM_VERSION} && chmod +x ./bin/gotestsum-${GOTESTSUM_VERSION}

SWAG ?= $(LOCALBIN)/swag
GOLINES ?= $(LOCALBIN)/golines
GOFUMPT ?= $(LOCALBIN)/gofumpt

install_swag: $(SWAG) ## install swagger
$(SWAG): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/swaggo/swag/cmd/swag@latest

install_golines: $(GOLINES) ## install golines
$(GOLINES): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install github.com/segmentio/golines@latest

install_gofumpt: $(GOFUMPT) ## install gofumpt
$(GOFUMPT): $(LOCALBIN)
	GOBIN=$(LOCALBIN) go install mvdan.cc/gofumpt@latest

install_golangcui-lint:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.3.1

install_gowatch:
	go install github.com/silenceper/gowatch@latest

install_subfinder:
	go install github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest

# 安装依赖注入生成器
install_wire:
	go get -tool github.com/google/wire/cmd/wire@latest

# 安装空引用静态分析器
install_nilaway:
	go get -tool go.uber.org/nilaway/cmd/nilaway@latest

# 安装第三方包漏洞检测工具
install_govulncheck:
	go get -tool golang.org/x/vuln/cmd/govulncheck@latest

# go语言模版代码生成器
install_oapi-codegen:
	go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

install_trpc:
	go install trpc.group/trpc-go/trpc-cmdline/trpc@latest

tools: bin/gotestsum install_golangcui-lint install_swag install_golines install_gofumpt install_gowatch install_subfinder install_wire install_nilaway install_govulncheck install_oapi-codegen

###############################################################################
# 工具使用
.PHONY: mocks sqlc wire audit docs govulncheck

mocks:
	go tool github.com/vektra/mockery/v2 --all --output ./mocks

sqlc:
	go tool github.com/kyleconroy/sqlc/cmd/sqlc generate

wire:
	go tool github.com/google/wire/cmd/wire ./...

docs: $(SWAG)
	# make gin_server
	# http://127.0.0.1:8000/swagger/index.html
	cd ginx && go tool swag init --parseDependency --parseDepth=6 --generalInfo ../main/main.go

# 空引用静态分析器
nilaway:
	go tool nilaway ./... || true

# 第三方包漏洞检测工具
govulncheck:
	go tool govulncheck ./... || true

# 审计
audit:
	go tool govulncheck go.uber.org/mock/mockgen
	go tool govulncheck github.com/swaggo/swag/cmd/swag
	go tool govulncheck github.com/vektra/mockery/v2
	go tool govulncheck github.com/google/wire/cmd/wire

###############################################################################
# 自定义命令
.PHONY: build test bench vet coverage check test-race test-cover-html help clear lint fix fmt tidy mocks generate
# 运行 make 命令而没有指定目标时，默认会执行 help 目标，并打印出帮助信息。
.DEFAULT_GOAL := help
help:	# Empty target rule

check: test-race vet

build_main:
	go build -o main main/main.go
	upx -9 main/main

build:
	go build -v ./...

build_small:
	# -ldflags="-s -w" 会去除符号表和调试信息
	go build -v -ldflags="-s -w" ./...

# gotestsum 是一个用于格式化和汇总 Go 测试输出的工具。它可以将 go test 的输出转换为更易读的格式，并且可以生成 JSON、JUnit XML 等格式的报告。
test: TEST_FORMAT ?= short
test: SHELL = /bin/bash
test: export CGO_ENABLED=1
test: go mod tidy
test: bin/gotestsum ## Run tests
	@mkdir -p ${BUILD_DIR}
	bin/gotestsum --no-summary=skipped --junitfile ${BUILD_DIR}/coverage.xml --format ${TEST_FORMAT} -- -coverprofile=${BUILD_DIR}/coverage.txt -covermode=atomic $(filter-out -v,${GOARGS}) $(if ${TEST_PKGS},${TEST_PKGS},./...)

test-race: bin/gotestsum ## Run tests with race
	@mkdir -p ${BUILD_DIR}
	bin/gotestsum --no-summary=skipped --junitfile ${BUILD_DIR}/coverage.xml --format ${TEST_FORMAT} -- -race -coverprofile=${BUILD_DIR}/coverage.txt -covermode=atomic $(filter-out -v,${GOARGS}) $(if ${TEST_PKGS},${TEST_PKGS},./...)

test_all: lint vet govulncheck nilaway
	go test ./... -gcflags=all=-l -cover

test2file:
	go test ./... -v -gcflags=all=-l -json > sn_report_test.json
	go test ./... -gcflags=all=-l -coverprofile=sn_report_covprofile
	go vet -json ./... 2> sn_report_vet_report.out
	golangci-lint run --out-format checkstyle ./... > sn_report_report.xml || true
	nilaway ./... > sn_report_nilaway.out || true

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

lint: ## Run linter
	# 会以彩色终端输出（如果支持）显示检查结果，包括错误、警告和建议的修复方式。
	# 这是开发时最常用的模式，适合在本地快速检查代码问题。
	# 默认情况下，golangci-lint 会一直运行直到完成所有检查（适合本地开发，代码量通常可控）。
	# --fast：跳过耗时检查（如 unused），加速本地 lint。
	golangci-lint run ./... || true

lint-ci:
	# CI/CD 环境使用
	# 强制将结果输出为纯文本格式到标准输出（stdout）
	# 这是为了兼容 CI 系统（如 GitHub Actions、GitLab CI）的日志收集机制，避免因彩色输出或特殊格式导致日志解析失败。
	# 超时能防止因 lint 耗时过长阻塞流水线。
	# 禁用缓存，确保 CI 中每次都是全新检查。
	golangci-lint run ./... --no-cache --output.text.path=stdout --timeout=5m

fix: ## Fix lint violations
	golangci-lint run --fix ./... || true

BENCH ?= .
bench:
# 	go test -run=nonthingplease -benchmem -bench=. $(shell go list ./... | grep -v /vendor/)
#	go test -v ./... -test.bench -test.benchmem
#	go test -bench=. -benchmem
	go list ./... | xargs -n1 go test -bench=$(BENCH) -run="^$$" $(BENCH_FLAGS)

vet:
	go vet ./... || true

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

generate:
	go generate ./...

# 安装 go install golang.org/x/tools/cmd/godoc
godoc:
	echo "http://127.0.0.1:6060"
	# http://127.0.0.1:6060/ginx/
	godoc -http=127.0.0.1:6060 -goroot="."

docker-build: ## build docker image
#	docker build -f ./Dockerfile -t xgo:${VERSION} .
	docker build -f ./Dockerfile -t xgo:latest .

gin_server:
	go run main/main.go
