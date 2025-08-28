# 简介
https://github.com/trpc-ecosystem/go-filter/blob/main/slime/README.zh_CN.md

| 协议 | 重试 | 对冲 | 备注 |
|:-:|:-:|:-:|:-|
|trpc|✓|✓| 原生的 trpc 协议。 |
|trpc SendOnly|✗|✗| 不支持，重试/对冲根据返回的错误码进行判断，而 SendOnly 请求不会回包。 |
|trpc 流式|✗|✗| 暂不支持。 |
|[http](https://github.com/trpc-group/trpc-go/tree/main/http)|✓|✓||
|[Kafaka](https://github.com/trpc-ecosystem/go-database/tree/main/kafka)|✓|✗| 不支持对冲功能。 |
|[MySQL](https://github.com/trpc-ecosystem/go-database/tree/main/mysql)|★|★| <span style="color:red">除 [Query](https://github.com/trpc-ecosystem/go-database/blob/6f75e87fecfc5411e54d93fd1aad5e7afa9a0fcf/mysql/client.go#L40) 和 [Transaction](https://github.com/trpc-ecosystem/go-database/blob/6f75e87fecfc5411e54d93fd1aad5e7afa9a0fcf/mysql/client.go#L42) 两个方法外，其他都支持</span>。这两个方法以函数闭包作为参数，slime 无法保证数据的并发安全性，可以使用 `slime.WithDisabled` 关闭重试/对冲。 |

# 使用
```go
import _ "trpc.group/trpc-go/trpc-filter/slime"
```

# 配置
## 重试
```yaml
--- # retry/hedging strategy
retry1: &retry1 # this is a yaml reference syntax that allows different services to use the same retry strategy
  # use a random name if omitted.
  # if you need to customize backoff or retryable business errors, you must explicitly provide a name, which will be
  # used as the first parameter of the slime.SetXXX method.
  name: retry1 # 重试策略名称
  # default as 2 when omitted.
  # no more than 5, truncate to 5 if exceeded.
  max_attempts: 4 # 最大尝试次数
  # 退避策略配置结构体
  backoff: # must provide one of exponential or linear
    exponential:
      initial: 10ms     # 初始延迟
      maximum: 1s       # 最大延迟
      multiplier: 2     # 乘数
  # 可重试错误码列表
  # when omitted, the following four framework errors are retried by default:
  #  21: RetServerTimeout
  # 111: RetClientConnectFail
  # 131: RetClientRouteErr
  # 141: RetClientNetErr
  # for tRPC-Go framework error codes, please refer to: https://github.com/trpc-group/trpc-go/tree/main/errs
  retryable_error_codes: [ 141 ]

retry2: &retry2
  name: retry2
  max_attempts: 4
  backoff:
    linear: [100ms, 500ms]
  retryable_error_codes: [ 141 ]
  # 是否跳过已访问节点
  skip_visited_nodes: false # omit, false and true correspond to three different cases
```

## 对冲
```yaml
hedging1: &hedging1
  # use a random name if omitted.
  # if you need to customize hedging_delay or non-fatal errors, you must explicitly provide a name, which will be used
  # as the first parameter of the slime.SetHedgingXXX method.
  # 对冲策略名称
  name: hedging1
  # default as 2 when omitted.
  # no more than 5, truncate to 5 if exceeded.
  # 最大尝试次数
  max_attempts: 4
  # 对冲延迟时间
  hedging_delay: 0.5s
  # when omitted, the following four errors default to non-fatal errors:
  # 21: RetServerTimeout
  # 111: RetClientConnectFail
  # 131: RetClientRouteErr
  # 141: RetClientNetErr
  # 非致命错误码列表
  non_fatal_error_codes: [ 141 ]

hedging2: &hedging2
  name: hedging2
  max_attempts: 4
  hedging_delay: 1s
  non_fatal_error_codes: [ 141 ]
  # 是否跳过已访问节点
  skip_visited_nodes: true # omit, false and true correspond to three different cases.
```

```yaml
client: &client
  filter: [slime] # filter must cooperate with plugin, both are indispensable
  service:
    - name: trpc.app.server.Welcome # 服务名称
      # 限流配置
      retry_hedging_throttle: # all retry/hedging strategies under this service will be bound to this rate limit
        max_tokens: 100
        token_ratio: 0.5
      # 重试和对冲配置
      retry_hedging: # service uses policy retry1 by default
        retry: *retry1 # dereference retry1
      # 方法配置
      methods:
        - callee: Hello # use retry policy retry2 instead of retry1 of parent service
          retry_hedging:
            retry: *retry2
        - callee: Hi # use hedging policy hedging1 instead of retry1 of parent service
          retry_hedging:
            hedging: *hedging1
        - callee: Greet # empty retry_hedging means no retry/hedging policy
          retry_hedging: {}
        - callee: Yo # retry_hedging is missing, use retry1 of parent service by default
    - name: trpc.app.server.Greeting
      retry_hedging_throttle: {} # forcibly turn of rate limit
      retry_hedging: # service uses hedging2 by default
        hedging: *hedging2
    - name: trpc.app.server.Bye
      # missing rate limit, use the default one.
      # there's no retry/hedging policy at service level.
      methods:
        - callee: SeeYou # SeeYou use retry1 as its own retry policy
          retry_hedging:
            retry: *retry1
```

```yaml
plugins:
  slime:
    # we reference the entire client here. Of course, you can configure client.service separately under default.
    default: *client
```
