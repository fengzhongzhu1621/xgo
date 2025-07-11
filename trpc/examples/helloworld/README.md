# 初始化项目
```sh
go mod init github.com/fengzhongzhu1621/xgo/trpc/examples/helloworld
go get github.com/envoyproxy/protoc-gen-validate/validate
```

# 生成工程
```sh
# 生成完整工程
trpc create -p helloworld.proto
go mod tidy
# 只生成 rpcstub，常用于已经创建工程以后更新协议字段时，重新生成桩代码
trpc create -p helloworld.proto --rpconly
# 使用 http 协议
trpc create -p helloworld.proto --protocol=http
```
