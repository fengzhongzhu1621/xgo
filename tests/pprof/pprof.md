# 1. profile CPU

## 1.1 代码直接生成
```go
func main() {
    f, _ := os.Create('CPU.out')
    defer f.Close()
    pprof.StartCPUProfile(f)
    defer   pprof.StopCPUProfile()
    ...
}
```

## 1.2 go test 参数生成
```go
go test -CPUprofile CPU.out . -run=TestFunc
```

## 1.3 通过 http 请求生成

引入 pprof 包后，会在默认处理器 DefaultServeMux 上注册 /debug/pprof/profile 接口的路由；调用 ListenAndServe 启动 http 服务，第二个参数传nil使用默认处理器处理请求。

```go
import (
  'net/http'
  _ 'net/http/pprof'
)

func pprofServerStart() {
  go func() {
    http.ListenAndServe('127.0.0.1:9000', nil)
  }()
}
```

## 1.4 采集数据
使用 go tool pprof 命令访问本地服务的 /debug/pprof/profile 接口，CPU 采样数据会自动保存到 $HOME/pprof/ 目录下，同时也会直接进入分析命令行。

```sh
$ go tool pprof http://127.0.0.1:9000/debug/pprof/profile?seconds=30

Entering interactive mode (type 'help' for commands, 'o' for options)
(pprof) 
```

## 1.5 命令行分析
```sh
go tool pprof $HOME/pprof/pprof.xgo.samples.cpu.001.pb.gz
```

## 1.6 可视化界面分析
```sh
go tool pprof -http=127.0.0.1:9888 $HOME/pprof/pprof.xgo.samples.cpu.001.pb.gz
```

# 2. allocs/heap 分析内存

## 2.1 代码直接生成

```go
func main() {
    f, _ := os.Create('mem.out')
    defer f.Close()
    runtime.GC() // 手动执行一次GC垃圾回收
    if err := pprof.WriteHeapProfile(f); err != nil {
        log.Fatal('could not write memory profile: ', err)
    }
    ...
}
```

## 2.2 go test 参数生成
```sh
go test -memprofile mem.out . -run=TestFunc
```

## 2.3 通过 http 接口生成
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

## 2.4 采集数据
使用 go tool pprof 命令访问本地服务的 /debug/pprof/heap 接口，获得 heap 采样数据保存到 $HOME/pprof/ 目录下，同时也会直接进入分析命令行

```sh
$ go tool pprof http://127.0.0.1:9000/debug/pprof/heap?seconds=30

Entering interactive mode (type 'help' for commands, 'o' for options)
(pprof) 
```

# 3. goroutine 分析 Golang 协程
## 3.1 代码直接生成
```go
func main() {
  f, _ := os.Create('goroutine.out')
  defer f.Close()
  err := pprof.Lookup('goroutine').WriteTo(f, 1)
  if err != nil {
    log.Fatal(err)
  }
}
```

## 3.2 通过 http 接口生成
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

## 3.3 采集数据
```

$ go tool pprof http://127.0.0.1:9000/debug/pprof/goroutine?seconds=30

(pprof)
```

## 3.4 性能影响
Go1.18及之前版本在循环遍历获取 goroutine 数据的整个过程中，会 stopTheWorld 停止进程的执行，生产环境需谨慎使用。

Go1.19版本之后，在获取 goroutine 数据过程中，只会有两次短暂的 stopTheWorld 停止整个进程，实测对程序整体的影响不大，生产环境对性能要求不高的场景仍然可以使用。
