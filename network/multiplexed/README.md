# Multiplexed 多路复用连接池

## 概述

`multiplexed` 是一个高性能的 Go 语言多路复用连接池实现，支持在单个物理连接上复用多个虚拟连接，显著减少网络连接开销，提高并发性能。

## 核心特性

### 🔄 多级连接复用
- **物理连接复用**：多个虚拟连接共享同一个物理连接
- **智能连接管理**：自动维护连接池，支持连接重连和空闲连接清理
- **负载均衡**：支持轮询和基于虚拟连接数的负载均衡策略

### 🛡️ 高可用性
- **自动重连机制**：连接断开时自动重连，支持指数退避策略
- **错误隔离**：单个连接错误不影响其他连接
- **连接健康检查**：自动检测和清理异常连接

### ⚡ 高性能
- **零拷贝设计**：减少内存分配和拷贝开销
- **并发安全**：所有操作都是线程安全的
- **异步 I/O**：读写操作异步执行，不阻塞业务逻辑

### 🔧 灵活配置
- **可配置连接数**：支持设置每个目标地址的连接数量
- **队列管理**：可配置发送队列大小和满队列处理策略
- **TLS 支持**：完整的 TLS 连接配置支持

## 架构设计

### 三级连接结构

```
Multiplexed (连接池)
    ↓
Connections (目标地址连接集合)
    ↓
Connection (物理连接)
    ↓
VirtualConnection (虚拟连接)
```

### 核心组件

1. **Multiplexed** - 顶层连接池管理器
2. **Connections** - 特定目标地址的连接集合
3. **Connection** - 物理连接实现
4. **VirtualConnection** - 虚拟连接接口
5. **IFrameParser** - 帧解析器接口

## 使用场景

### 1. 高并发微服务通信
适用于微服务架构中服务间的频繁通信，减少连接建立开销。

```go
// 创建多路复用连接池
pool := multiplexed.New(
    multiplexed.WithConnectNumber(2),
    multiplexed.WithQueueSize(1024),
)

// 获取虚拟连接
conn, err := pool.GetMuxConn(ctx, "tcp", "service:8080", opts)
```

### 2. 实时数据流处理
适合需要处理大量实时数据流的场景，如消息队列、实时监控等。

### 3. 长连接应用
适用于需要维持长连接的场景，如 WebSocket 服务、实时游戏等。

### 4. 资源受限环境
在连接数受限的环境中，通过多路复用减少实际连接数。

## 快速开始

### 基本用法

```go
package main

import (
    "context"
    "github.com/fengzhongzhu1621/xgo/network/multiplexed"
)

func main() {
    // 使用默认连接池
    pool := multiplexed.DefaultMultiplexedPool
    
    // 创建获取选项
    opts := multiplexed.NewGetOptions()
    opts.WithFrameParser(&MyFrameParser{})
    opts.WithVID(1) // 设置虚拟连接ID
    
    // 获取虚拟连接
    ctx := context.Background()
    conn, err := pool.GetMuxConn(ctx, "tcp", "localhost:8080", opts)
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    
    // 使用连接进行通信
    conn.Write([]byte("hello"))
    data, _ := conn.Read()
}
```

### 自定义帧解析器

```go
type MyFrameParser struct{}

func (p *MyFrameParser) Parse(rc io.Reader) (vid uint32, buf []byte, err error) {
    // 实现自定义帧解析逻辑
    // 返回虚拟连接ID、数据帧和错误
    return 1, []byte("data"), nil
}
```

### 高级配置

```go
// 创建自定义配置的连接池
pool := multiplexed.New(
    multiplexed.WithConnectNumber(4),           // 每个目标4个连接
    multiplexed.WithQueueSize(2048),           // 发送队列大小
    multiplexed.WithDropFull(true),            // 队列满时丢弃
    multiplexed.WithDialTimeout(5*time.Second), // 拨号超时
    multiplexed.WithMaxVirConnsPerConn(100),   // 每个连接最多100个虚拟连接
    multiplexed.WithMaxIdleConnsPerHost(2),    // 最大空闲连接数
)
```

## 配置选项

### 连接池配置 (PoolOptions)

| 选项 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `connectNumberPerHost` | `int` | 2 | 每个目标地址的连接数量 |
| `sendQueueSize` | `int` | 1024 | 每个连接的发送队列大小 |
| `dropFull` | `bool` | `false` | 队列满时是否丢弃请求 |
| `dialTimeout` | `time.Duration` | 1秒 | 连接超时时间 |
| `maxVirConnsPerConn` | `int` | 0 | 每个物理连接的最大虚拟连接数（0=无限制） |
| `maxIdleConnsPerHost` | `int` | 0 | 每个目标的最大空闲连接数 |

### 连接获取配置 (GetOptions)

| 选项 | 类型 | 说明 |
|------|------|------|
| `FP` | `IFrameParser` | 帧解析器接口 |
| `VID` | `uint32` | 虚拟连接ID |
| `CACertFile` | `string` | CA证书文件路径 |
| `TLSCertFile` | `string` | 客户端证书文件路径 |
| `TLSKeyFile` | `string` | 客户端私钥文件路径 |
| `TLSServerName` | `string` | TLS服务器名称验证 |
| `LocalAddr` | `string` | 本地绑定地址 |

## 错误处理

包定义了丰富的错误类型，便于调试和错误处理：

```go
var (
    ErrFrameParserNil              = errors.New("frame parser is nil")
    ErrRecvQueueFull               = errors.New("virtual connection's recv queue is full")
    ErrSendQueueFull               = errors.New("connection's send queue is full")
    ErrChanClose                   = errors.New("unexpected recv chan close")
    ErrAssertFail                  = errors.New("type assert fail")
    ErrDupRequestID                = errors.New("duplicated Request ID")
    ErrInitPoolFail                = errors.New("init pool for specific node fail")
    ErrWriteNotFinished            = errors.New("write not finished")
    ErrNetworkNotSupport           = errors.New("network not support")
    ErrConnectionsHaveBeenExpelled = errors.New("connections have been expelled")
)
```

## 性能优化建议

1. **合理设置连接数**：根据业务并发量调整 `connectNumberPerHost`
2. **优化队列大小**：根据数据包大小和频率设置 `sendQueueSize`
3. **实现高效帧解析**：自定义 `IFrameParser` 实现零拷贝解析
4. **监控连接状态**：定期检查连接池健康状态

## 文件说明

| 文件 | 说明 |
|------|------|
| `multiplexed.go` | 主入口文件，连接池管理器 |
| `connections.go` | 目标地址连接集合管理 |
| `connection.go` | 物理连接实现 |
| `mux_conn.go` | 虚拟连接实现 |
| `define.go` | 接口定义和错误常量 |
| `get_options.go` | 连接获取配置 |
| `pool_options.go` | 连接池配置 |
| `utils.go` | 工具函数 |

## 限制和注意事项

1. **帧解析器要求**：必须实现 `IFrameParser` 接口来处理自定义协议
2. **连接生命周期**：虚拟连接需要手动关闭以避免资源泄漏
3. **错误处理**：需要正确处理连接错误和重连逻辑
4. **内存使用**：大量虚拟连接可能增加内存开销

## 扩展性

包设计具有良好的扩展性，可以通过实现以下接口来自定义行为：

- `IFrameParser` - 自定义帧解析逻辑
- `IPool` - 自定义连接池实现
- `IMuxConn` - 自定义虚拟连接行为

## 总结

`multiplexed` 包提供了一个强大而灵活的多路复用连接池解决方案，特别适合需要高性能、高并发网络通信的场景。通过合理的配置和使用，可以显著提升应用程序的网络性能和资源利用率。