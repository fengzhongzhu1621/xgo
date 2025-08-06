# 插件概述
Hystrix是Netflix开源的容错库，用于隔离对远程系统/服务的访问点，防止级联故障，提升分布式系统的弹性。该插件基于hystrix-go实现（参考https://github.com/afex/hystrix-go）

# 核心功能
熔断机制：当错误率超过阈值时自动熔断，避免系统过载
隔离策略：支持线程隔离（THREAD）和信号量隔离（SEMAPHORE）
监控指标：提供请求量、错误率、超时等实时数据采集
降级处理：支持配置fallback方法返回默认值

# 熔断器工作原理补充

## 1. 熔断器状态
* 关闭(CLOSED)
* 开启(OPEN)
* 半开(HALF_OPEN)

## 2. 状态转换：

当错误率超过阈值时从 CLOSED 转为 OPEN
OPEN 状态持续 SleepWindow 时间后转为 HALF_OPEN
HALF_OPEN 状态下尝试放行部分请求，成功则转回 CLOSED，失败则保持 OPEN


## 3. 关键配置：
* Timeout：执行超时时间
* MaxConcurrentRequests：最大并发请求数
* RequestVolumeThreshold：触发熔断的最小请求数
* ErrorPercentThreshold：错误百分比阈值
* SleepWindow：熔断后尝试恢复的时间窗口

# 使用插件
```go
import （
	_ "trpc.group/trpc-go/trpc-filter/hystrix"
）
```

# 过滤器配置示例
```go
# If you need to fuse the method of the server, configure the server-side filter.
server:
  ...
  filter:
    ...
    - hystrix
# If you need to fuse the method of the client, configure the client-side filter.
client:
  ...
  filter:
    ...
    - hystrix
```

# 插件配置示例
```go
plugins:                                     # Plugin configuration
  circuitbreaker:
    hystrix:
      /trpc.qq_news.user_info.UserInfo/Api1: # Business routing【server】trpc.Message(ctx).ServerRPCName(); 【client】trpc.Message(ctx).ClientRPCName()
        timeout: 1000              # Overtime（ms）
        maxconcurrentrequests: 100 # The maximum number of concurrent requests.
        requestvolumethreshold: 30 # After more than ? requests, the fuse will be turned on according to the error ratio.
        sleepwindow: 2000          # Fuse time（ms）
        errorpercentthreshold: 10  # Turn on the error ratio of fusing.
      /trpc.qq_news.user_info.UserInfo/Api2:
        timeout: 2000 # 超时时间(ms)
        maxconcurrentrequests: 100 # 最大并发数
        requestvolumethreshold: 3 # 触发熔断的最小请求量
        sleepwindow: 5000 # 熔断恢复时间(ms)
        errorpercentthreshold: 10 # 错误百分比阈值(%)
     "*": # Set up wildcards. # 全局默认配置
        timeout: 2000
        maxconcurrentrequests: 100
        requestvolumethreshold: 3
        sleepwindow: 5000
        errorpercentthreshold: 10
     _/trpc.qq_news.user_info.UserInfo/Api3: # Exclude an interface when global configuration is enabled. # 排除特定接口
     _/trpc.qq_news.user_info.UserInfo/Api4: # Exclude an interface when global configuration is enabled. # 排除特定接口
```

注意事项
* 优先级：特定路由配置 > 全局配置（""）
* 特殊符号：_前缀表示排除接口，*表示通配配置
* 错误返回：熔断触发时返回hystrix: circuit open错误

# 注册收集器
```go
type testMetricCollector struct{
 attemptsPrefix          string
 errorsPrefix            string
 successesPrefix         string
 failuresPrefix          string
 rejectsPrefix           string
 shortCircuitsPrefix     string
 timeoutsPrefix          string
 fallbackSuccessesPrefix string
 fallbackFailuresPrefix  string
}

// Update ...
func (m *testMetricCollector) Update(r metricCollector.MetricResult) {
     // Use the metric in trpc-go for data calculation.
}

// Reset ...
func (m *testMetricCollector) Reset() {}

func newTestMetircCollector(name string) metricCollector.MetricCollector {
 return &testMetricCollector{
 	attemptsPrefix:          name + ".attempts",
 	errorsPrefix:            name + ".errors",
 	successesPrefix:         name + ".successes",
 	failuresPrefix:          name + ".failures",
 	rejectsPrefix:           name + ".rejects",
 	shortCircuitsPrefix:     name + ".shortCircuits",
 	timeoutsPrefix:          name + ".timeouts",
 	fallbackSuccessesPrefix: name + ".fallbackSuccesses",
 	fallbackFailuresPrefix:  name + ".fallbackFailures",
 }
}

hystrix.RegisterCollector(newTestMetircCollector)
```
