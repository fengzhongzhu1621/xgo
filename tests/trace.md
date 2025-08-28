# 1. 生成 trace 数据

## 1.1 代码直接生成
```go
func main() {
  f, _ := os.Create('trace.out')
  defer f.Close()
  trace.Start(f)
  defer trace.Stop()
  ...
}
```

## 1.2 通过 http 接口生成
```go

import (
  'net/http'
  _ 'net/http/pprof'
)

func process() {
  go func() {
    http.ListenAndServe('127.0.0.1:9000', nil)
  }()
    ...
}
```

请求 http 接口生成 trace 数据
```sh
$ curl http://127.0.0.1:9000/debug/pprof/trace?seconds=10 > trace.data
```

# 2. 分析 trace 数据
```sh
$ go tool trace -http 127.0.0.1:9998 trace.out
```

## 2.1 性能影响
因为每一个事件都要执行埋点代码，所以开启 trace 对性能影响非常大，事件越多影响越大，预期至少是30%以上的额外性能开销，trace 在生产环境中请谨慎使用。

