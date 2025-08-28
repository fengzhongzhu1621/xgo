# 测试

http

```sh
curl -X POST -d '{"msg":"hello"}' -H "Content-Type:application/json" "http://127.0.0.1:8002/trpc.examples.helloworld.GreeterHttp/SayHello" | jq
curl -X POST -d '{"msg":""}' -H "Content-Type:application/json" "http://127.0.0.1:8002/trpc.examples.helloworld.GreeterHttp/SayHello" | jq
curl -X POST -d '{"msg":"hello"}' -H "Content-Type:application/json" "http://127.0.0.1:8002/cgi-bin/hello" | jq
```

trpc
```sh
trpc-cli -func /trpc.examples.helloworld.Greeter/SayHello -target ip://127.0.0.1:8001 -body '{"msg":"hello"}'
trpc-cli -func /cgi-bin/hello -target ip://127.0.0.1:8001 -body '{"msg":"hello"}'
```

# 管理端

## 查看管理命令列表

```sh
// 查看所有命令
curl "http://127.0.0.1:11014/cmds" | jq
// 框架版本
curl "http://127.0.0.1:11014/version" | jq
// 日志级别
curl "http://127.0.0.1:11014/cmds/loglevel?logger=default&output=0" | jq
// 框架配置
curl "http://127.0.0.1:11014/cmds/config" | jq

// 自定义命令
curl "http://127.0.0.1:11014/cmds/custom" | jq
```

```
// 健康检查
curl "http://127.0.0.1:11014/is_healthy/"
```
| HTTP 状态码 | 服务状态 |
|-------------|----------|
| 200         | 健康     |
| 404         | 未知     |
| 503         | 不健康   |


## pprof

使用 go tool pprof 命令访问本地服务的 /debug/pprof/profile 接口，CPU 采样数据会自动保存到 $HOME/pprof/ 目录下，同时也会直接进入分析命令行。
```sh
go tool pprof http://127.0.0.1:11014/debug/pprof/profile?seconds=20
```

```sh
curl "http://127.0.0.1:11014/debug/pprof/profile?seconds=20" > profile.out
go tool pprof profile.out

curl "http://127.0.0.1:11014/debug/pprof/trace?seconds=20" > trace.out
go tool trace trace.out
```
