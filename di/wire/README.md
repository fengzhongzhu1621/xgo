# 简介
wire是 Google 开源的一个依赖注入工具。它是一个代码生成器，并不是一个框架。我们只需要在一个特殊的go文件中告诉wire类型之间的依赖关系，它会自动帮我们生成代码，帮助我们创建指定类型的对象，并组装它的依赖。

```sh
go install github.com/google/wire/cmd/wire
```

上面的命令会在$GOPATH/bin中生成一个可执行程序wire，这就是代码生成器。我个人习惯把$GOPATH/bin加入系统环境变量$PATH中，所以可直接在命令行中执行wire命令。


# go:build wireinject

```go
//go:build wireinject
// +build wireinject
```

确保该文件仅在生成代码时被编译，运行时不会包含它。

Wire 通过解析标记了 //go:build wireinject 的文件（通常是 _wire.go 或类似名称的文件），生成实际的依赖注入代码。这些生成代码会被写入到另一个文件中（如 wire_gen.go），而原始的 wireinject 文件不会出现在最终的二进制中。

Wire 的注入逻辑（如 wire.Build 调用）仅用于代码生成阶段，运行时不需要这些逻辑，因此通过构建标签将其隔离。

Go 的 go build 或 go generate 命令会自动识别构建标签，确保 wireinject 文件仅在生成代码时被处理。

# go generate ./...
要触发“生成”动作有两种方式：go generate 或 wire 。前者仅在 wire_gen.go 已存在的情况下有效（因为 wire_gen.go 的第三行 //*go:generate* wire），而后者在任何时候都有可以调用。并且后者有更多参数可以对生成动作进行微调， 所以建议始终使用 wire 命令。

```
wire ./...
```
