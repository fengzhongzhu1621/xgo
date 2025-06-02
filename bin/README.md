# go tool
Go 1.24 引入的 tool 指令为工具依赖管理提供了官方解决方案。它允许你在 go.mod 文件中直接声明工具依赖，就像管理代码依赖一样简单。

## 安装工具
```go
go get -tool github.com/golang/mock/mockgen
```

## 使用工具
```go
go mod tidy
go tool mockgen -source=internal/service.go -destination=internal/mocks/service_mock.go -package=mocks
````

```go
go tool
```

## 本地编译工具加速开发
```go
# 编译依赖的工具到当前目录
go build tool

# 编译到项目的 bin 目录
go build -o bin/ tool
```
