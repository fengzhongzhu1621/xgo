# 安装
```sh
go install github.com/cweill/gotests/...@latest
```
* github.com/cweill/gotests/... 表示安装 gotests 及其所有子命令。
* @latest 表示安装最新版本（也可以指定版本，如 @v1.6.0）。

# 生成测试文件
```sh
gotests -w -all yourfile.go
```

* -w：将生成的测试文件写入当前目录（而不是仅打印到 stdout）。
* -all：为文件中的所有函数生成测试（默认只生成导出的函数）。

## 为包中的所有文件生成测试

```sh
gotests -w -all ./...
```

## 仅为特定函数生成测试
```sh
gotests -w -excl "Test|Example" yourfile.go
```
