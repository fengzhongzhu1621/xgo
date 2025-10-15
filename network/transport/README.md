## 背景

框架间支持多种通信方式，如 tcp、udp 等。对于 udp 协议，一个 udp 包就对应一个 RPC 请求或回包。对于 tcp 这样的流式协议，就需要框架额外做分包处理。为了隔离不同网络协议间的差异，提供了 transport 抽象。

## 原理


- client transport 负责和对端建立连接，并提供 multiplexed 等高级特性；
- server transport 负责建立监听套接字并 accept 新连接请求，并处理连接上到达的请求。


## ClientTransport

## ServerTransport
和 client 的 `RoundTripOptions` 对应，server 也可通过 [`ServerTransportOptions`](server_listenserve_options.go)，设置异步处理、空闲超时、tls 证书等：

```go
st := transport.NewServerTransport(transport.WithServerAsync(true))
```

## ClientStreamTransport

[ClientStreamTransport](transport_stream.go) 用于发送/接收流式请求。因为 stream 是 client 发起创建的，所以，它提供了 `Init` 方法来对流进行初始化，比如与对端建立网络连接。

client stream transport 用了与普通 RPC transport 相同的 `RoundTripOption`，它底层的连接也支持多路复用等。


## ServerStreamTransport

[ServerStreamTransport](transport_stream.go) 用于服务端处理流式请求。当 Server 端收到 client 的 Init 包之后，它会创建一个新协程运行用户业务逻辑，而原始的网络收包协程则负责将收到的包分发给新协程。

注意，ServerStreamTransport embedding 了 `ServerTransport` 用于监听端口并创建对应的网络协程。所以，普通 RPC 的 `ListenServeOption` 对流式 server 也适用。


## 分包

tRPC 的包都由帧头、包头、包体组成。在 server 收到请求和 client 收到回包时（流式请求也适用），需要对原始数据流分割成一个个请求，然后交给对应的处理逻辑。[`codec.FramerBuild`](/codec/framer_builder.go) 和 [`codec.Framer`](/codec/framer_builder.go) 就是用来对数据流进行分包的。

在 client 端，可以通过 [`WithClientFramerBuilder`](client_roundtrip_options.go) 设置 frame builder，在 server 端，可以通过 [`WithServerFramerBuilder`](server_listenserve_options.go) 设置。

