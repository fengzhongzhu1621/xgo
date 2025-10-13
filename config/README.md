https://github.com/trpc-group/trpc-go/blob/main/config/README.zh_CN.md

## 如何管理业务配置
对于业务配置的管理，我们建议最佳实践是使用配置中心来管理业务配置，使用配置中心有以下优点：
- 避免源代码泄露敏感信息
- 服务动态更新配置
- 多服务共享配置，避免一份配置拥有多个副本
- 支持灰度发布，配置回滚，拥有完善的权限管理和操作日志

业务配置也支持本地文件。对于本地文件，大部分使用场景是客户端作为独立的工具使用，或者程序在开发调试阶段使用。好处在于不需要依赖外部系统就能工作。

## 什么是多数据源
数据源就获取配置的来源，配置存储的地方。常见的数据源包括：file，etcd，configmap，env 等。tRPC 框架支持对不同业务配置设定不同的数据源。框架采用插件化方式来扩展对更多数据源的支持。在后面的实现原理章节，我们会详细介绍框架是如何实现对多数据源的支持的。

## 什么是 Codec
业务配置中的 Codec 是指从配置源获取到的配置的格式，常见的配置文件格式为：yaml，json，toml 等。框架采用插件化方式来扩展对更多解码格式的支持。

config 接口为业务代码提供了获取配置项的标准接口，每种数据类型都有一个独立的接口，接口支持返回 default 值。

Codec 和 DataProvider 这两个模块都提供了标准接口和注册函数以支持编解码和数据源的插件化。以实现多数据源为例，DataProvider 提供了以下三个标准接口，其中 Read 函数提供了如何读取配置的原始数据（未解码），而 Watch 函数提供了 callback 函数，当数据源的数据发生变化时，框架会执行此 callback 函数。
```go
type DataProvider interface {
    Name() string
    Read(string) ([]byte, error)
    Watch(ProviderCallback)
}
```

最后我们来看看，如何通过指定数据源，解码器来获取一个业务配置项：
```go
// 加载 etcd 配置文件：config.WithProvider("etcd")
c, _ := config.Load("test.yaml", config.WithCodec("yaml"), config.WithProvider("etcd"))
// 读取 String 类型配置
c.GetString("auth.user", "admin")
```
在这个示例中，数据源为 etcd 配置中心，数据源中的业务配置文件为“test.yaml”。当 ConfigLoader 获取到"test.yaml"业务配置时，指定使用 yaml 格式对数据内容进行解码。最后通过`c.GetString("server.app", "default")`函数来获取 test.yaml 文件中`auth.user`这个配置型的值。


```go
// 加载配置文件：path 为配置文件路径
func Load(path string, opts ...LoadOption) (Config, error)
// 更改编解码类型，默认为“yaml”格式
func WithCodec(name string) LoadOption
// 更改数据源，默认为“file”
func WithProvider(name string) LoadOption
```

示例代码为：
```go
// 加载 etcd 配置文件：config.WithProvider("etcd")
c, _ := config.Load("test1.yaml", config.WithCodec("yaml"), config.WithProvider("etcd"))

// 加载本地配置文件，codec 为 json，数据源为 file
c, _ := config.Load("../testdata/auth.yaml", config.WithCodec("json"), config.WithProvider("file"))

// 加载本地配置文件，默认为 codec 为 yaml，数据源为 file
c, _ := config.Load("../testdata/auth.yaml")
```


从 config 数据结构中获取指定配置项值。支持设置默认值，框架提供以下标准接口：
```go
// Config 配置通用接口
type Config interface {
    Load() error
    Reload()
    Get(string, interface{}) interface{}
    Unmarshal(interface{}) error
    IsSet(string) bool
    GetInt(string, int) int
    GetInt32(string, int32) int32
    GetInt64(string, int64) int64
    GetUint(string, uint) uint
    GetUint32(string, uint32) uint32
    GetUint64(string, uint64) uint64
    GetFloat32(string, float32) float32
    GetFloat64(string, float64) float64
    GetString(string, string) string
    GetBool(string, bool) bool
    Bytes() []byte
}
```

示例代码为：
```go
// 读取 bool 类型配置
c.GetBool("server.debug", false)

// 读取 String 类型配置
c.GetString("server.app", "default")
```

## 监听配置项
对于 KV 型配置中心，框架提供了 Watch 机制供业务程序根据接收的配置项变更事件，自行定义和执行业务逻辑。监控接口设计如下：
```go

// Get 根据名字使用 kvconfig
func Get(name string) KVConfig

// KVConfig kv 配置
type KVConfig interface {
    KV
    Watcher
    Name() string
}

// 监控接口定义
type Watcher interface {
    // Watch 监听配置项 key 的变更事件
    Watch(ctx context.Context, key string, opts ...Option) (<-chan Response, error)
}

// Response 配置中心响应
type Response interface {
    // Value 获取配置项对应的值
    Value() string
    // MetaData 额外元数据信息
    // 配置 Option 选项，可用于承载不同配置中心的额外功能实现，例如 namespace,group, 租约等概念
    MetaData() map[string]string
    // Event 获取 Watch 事件类型
    Event() EventType
}

// EventType 监听配置变更的事件类型
type EventType uint8
const (
    // EventTypeNull 空事件
    EventTypeNull EventType = 0
    // EventTypePut 设置或更新配置事件
    EventTypePut EventType = 1
    // EventTypeDel 删除配置项事件
    EventTypeDel EventType = 2
)
```

下面示例展示了业务程序监控 etcd 上的“test.yaml”文件，打印配置项变更事件并更新配置。
```go
import (
    "sync/atomic"
    // ...
)

type yamlFile struct {
    Server struct {
        App string
    }
}

var cfg atomic.Value // 并发安全的 Value

// 使用 trpc-go/config 中 Watch 接口监听 etcd 远程配置变化
c, _ := config.Get("etcd").Watch(context.TODO(), "test.yaml")

go func() {
    for r := range c {
        yf := &yamlFile{}
        fmt.Printf("event: %d, value: %s", r.Event(), r.Value())

        if err := yaml.Unmarshal([]byte(r.Value()), yf); err == nil {
            cfg.Store(yf)
        }
    }
}()

// 当配置初始化完成后，可以通过 atomic.Value 的 Load 方法获得最新的配置对象
cfg.Load().(*yamlFile)
```

# 数据源实现

参考：[trpc-ecosystem/go-config-etcd](https://github.com/trpc-ecosystem/go-config-etcd)
