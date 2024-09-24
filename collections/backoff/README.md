# 退避算法

操作重试有不同的策略（退避算法），同时操作重试在很场景中都使用到了，比如IP数据包发送(CSMA/CD)、网络通信，RPC服务调用等，有的是用于协调网络传输速率，避免网络拥塞，有的是为了考虑网络波动影响，提高服务可用性

常见的是指数退避算法，通常起先是基于一个较低的时间间隔尝试操作，若尝试失败，则按指数级的逐步延长事件发生的间隔时间，直至超过最大尝试机会；大多数指数退避算法会利用抖动（随机延迟）来防止连续的冲突。


# cenkalti/backoff

可以用来对我们需要冥等补偿的业务逻辑进行重试，我们可以设定一个最大间隔时间， 停止时间等重试规则

退避算法有两块，一块是退避策略，即间隔多久操作下一次（退避策略）；另一块是累计可以支持的最大操作次数（重试次数）


## BackOff类型
BackOff是一个接口类型，提供NextBackOff()下次重试的间隔时间，Reset()支持退避恢复到初始状态

```go
Copy
type BackOff interface {
    // NextBackOff returns the duration to wait before retrying the operation,
    // or backoff. Stop to indicate that no more retries should be made.
    //
    // Example usage:
    //
    //  duration := backoff.NextBackOff();
    //  if (duration == backoff.Stop) {
    //      // Do not retry operation.
    //  } else {
    //      // Sleep for duration and retry operation.
    //  }
    //
    NextBackOff() time.Duration

    // Reset to initial state.
    Reset()
}
```

## 退避策略

* NewExponentialBackOff() ：（失败后，基于指数级别的间隔，再次发起操作；考虑到了网络实际情况）
* ZeroBackOff()：零间隔时间退避（失败立马再次发起请求，针对网络抖动，通常会遇到连续失败情况）
* ConstantBackOff()： 相同时间间隔的退避（失败后，间隔一个指定的时间，再次发起操作）

```go
    // exponential back off
    bkf := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), MaxRetry)
    err := backoff.Retry(opfn, bkf)
```

## 重试次数

* WithMaxRetries(b BackOff, max uint64) BackOff：创建一个退避策略，选定一个BackOff类型，并设置最大尝试次数
* Retry(o Operation, b BackOff) error：基于上述设定的退避策略，利用Retry函数，绑定到操作Operation上

```go
// 操作类型
type Operation func() error
// PermanentError表示不继续重试该操作（比如参数错误、数据异常时候，重试也是徒劳）
type PermanentError struct {
    Err error
}
```

## 指数退避算法
* NewExponentialBackOff()：创建指数退避类型
* GetElapsedTime() time.Duration：获取自创建退避实例以来经过的时间，在调用Reset()时重置。
* NextBackOff() time.Duration：使用公式计算下一个退避间隔
* Reset()：可以通过Reset()重置退避算法的耗时时间（最大耗时* * * DefaultMaxElapsedTime=15分钟）

指数退避，下次操作(NextBackOff)触发时间间隔计算公式：

```go
randomized interval = 
    RetryInterval * (random value in range [1 - RandomizationFactor, 1 + RandomizationFactor])
Randomized interval = RetryInterval +/- (RandomizationFactor * RetryInterval)
```

## Ticker定时器

```go
Copy
// 创建一个指定退避策略的定时器，这样当下次退避策略时限到达，将会自动产生一个时间消息到定时器通道
func NewTicker(b BackOff) *Ticker
// 定时器支持停止
func (t *Ticker) Stop()
```

## Context上下文相关（值传递、ctx取消）

```go
type BackOffContext interface {
    BackOff
    Context() context.Context
}
func WithContext(b BackOff, ctx context.Context) BackOffContext
```


## ExponentialBackOff默认的配置
```go
// Default values for ExponentialBackOff.
const (
	DefaultInitialInterval     = 500 * time.Millisecond
	DefaultRandomizationFactor = 0.5
	DefaultMultiplier          = 1.5
	DefaultMaxInterval         = 60 * time.Second
	DefaultMaxElapsedTime      = 15 * time.Minute
```
