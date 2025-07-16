# 生产者
| 参数名称         | 含义                          | 可选值&说明 |
|------------------|-------------------------------|-------------|
| `ip:port`        | 地址列表                      | 复数地址使用逗号分割，支持 域名:port，ip:port。暂不支持 cl5 |
| `clientid`       | 生产者 ID                     | 若为虫洞 kafka 需要到管理页面注册 |
| `topic`          | 生产的 topic                  | 调用 Produce 必需。若为虫洞 kafka 需要到管理页面注册 |
| `version`        | 客户端版本                    | 支持以下两种格式的版本号 0.1.1.0 / 1.1.0 |
| `partitioner`    | 消息分片方式                  | `random`: sarama.NewRandomPartitioner(默认值);<br>`roundrobin`: sarama.NewRoundRobinPartitioner;<br>`hash`: sarama.NewHashPartitioner 暂未自定义 hash 方式; |
| `compression`    | 压缩方式                      | `none`: sarama.CompressionNone;<br>`gzip`: sarama.CompressionGZIP(默认值);<br>`snappy`: sarama.CompressionSnappy;<br>`lz4`: sarama.CompressionLZ4;<br>`zstd`: sarama.CompressionZSTD; |
| `maxMessageBytes`| Msg 最大长度                  | 默认 131072 |
| `requiredAcks`   | 是否需要回执                  | 生产消息时，broker 返回消息回执（ack）的模式，支持以下 3 个值 (0/1/-1):<br>0: NoResponse，不需要等待 broker 响应。<br>1: WaitForLocal，等待本地（leader 节点）的响应即可返回。<br>-1: WaitForAll，等待所有节点（leader 节点和所有 In Sync Replication 的 follower 节点）均响应后返回。（默认值） |
| `maxRetry`       | 失败最大重试次数              | 生产消息最大重试次数，默认为 3 次，注意：必须大于等于 0，否则会报错 |
| `retryInterval`  | 失败重试间隔                  | 单位毫秒，默认 100ms |
| `trpcMeta`       | 透传 trpc 元数据到 sarama header | `true` 表示开启透传，`false` 表示不透传，默认 false |
| `discover`       | 用于服务发现的 discovery 类型   | 例如：polaris |
| `namespace`      | 服务命名空间                  | 例如：Development |
| `idempotent`     | 是否开启生产者幂等            | `true/false`，默认 false |

# 消费者
参数名称	含义	可选值&说明
ip:port 列表	地址	复数地址使用逗号分割，支持域名:port，ip:port。暂不支持 cl5
group	消费者组	若为虫洞 kafka 需要到管理页面注册
clientid	客户端 id	连接 kafka 的客户端 id
topics	消费的的 toipc	复数用逗号分割
compression	压缩方式
strategy	策略
sticky: sarama.BalanceStrategySticky;
range : sarama.BalanceStrategyRange;
roundrobin: sarama.BalanceStrategyRoundRobin;
fetchDefault	拉取消息的默认大小（字节），如果消息实际大小大于此值，需要重新分配内存空间，会影响性能，等价于 sarama.Config.Consumer.Fetch.Default	默认 524288
fetchMax	拉取消息的最大大小（字节），如果消息实际大小大于此值，会直接报错，等价于 sarama.Config.Consumer.Fetch.Max	默认 1048576
batch	每批次数目	使用批量消费时必填，注册批量消费函数时，batch 不填会发生参数不匹配导致消费失败，使用参考 examples/batchconsumer
batchFlush	批次消费间隔	默认 2 秒，单位 ms, 表示当批量消费不满足最大条数时，强制消费的间隔
initial	初始消费位置
新消费者组第一次连到集群消费的位置
newest: 最新位置
oldest: 最老位置
maxWaitTime	单次消费拉取请求最长等待时间	最长等待时间仅在没有最新数据时才会等待，默认 1s
maxRetry	失败最大重试次数	超过后直接确认并继续消费下一条消息，默认 0:没有次数限制，一直重试、负数表示不重试，直接确认并继续消费下一条消息，正数表示如果一直错误最终执行次数为 maxRetry+1
netMaxOpenRequests	最大同时请求数	网络层配置，最大同时请求数，默认 5
maxProcessingTime	消费者单条最大请求时间	单位 ms，默认 1000ms
netDailTimeout	链接超时时间	网络层配置，链接超时时间，单位 ms，默认 30000ms
netReadTimeout	读超时时间	网络层配置，读超时时间，单位 ms，默认 30000ms
netWriteTimeout	写超时时间	网络层配置，写超时时间，单位 ms，默认 30000ms
groupSessionTimeout	消费组 session 超时时间	单位 ms，默认 10000ms
groupRebalanceTimeout	消费者 rebalance 超时时间	单位 ms，默认 60000ms
mechanism	使用密码时加密方式	可选值 SCRAM-SHA-512/SCRAM-SHA-256
user	用户名
password	密码
retryInterval	重试间隔	单位毫秒，默认 3000ms
isolationLevel	隔离级别	可选值 ReadUncommitted/ReadCommitted
trpcMeta	透传 trpc 元数据，读取 sarama header 设置 trpc meta	true 表示开启透传，false 表示不透传，默认 false
discover	用于服务发现的 discovery 类型	例如：polaris
namespace	服务命名空间	例如：Development
