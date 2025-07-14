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
