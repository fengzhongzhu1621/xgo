# 默认重试策略
指数退避增长（BackOffDelay）加上一个 0~100ms 的随机抖动（RandomDelay）。默认最多重试 10 次。
```
attempts:         uint(10),  // 默认最多重试 10 次
delay:            100 * time.Millisecond, // 初始基础间隔 100ms
maxJitter:        100 * time.Millisecond, // 最大抖动 100ms
delayType:        CombineDelay(BackOffDelay, RandomDelay), // 指数退避 + 随机抖动
```

# 最大重试次数

retry.Attempts(0)，无限重试，直到成功

```go
retry.Do(
  task,
  retry.Attempts(5), // 最大重试 5 次
 )
```

# 重试间隔
```go
retry.Do(
  task,
  retry.Delay(200*time.Millisecond), // 每次重试间隔 200 毫秒
 )
```

# 指数回退
```go
 retry.Do(
  task,
  retry.DelayType(retry.BackOffDelay), // 指数回退策略
 )
```

# 随机延迟
```go
 retry.Do(
  task,
  retry.DelayType(retry.RandomDelay), // 随机回退策略
 )
```

# 固定延迟
```go
retry.Do(
  task,
  retry.DelayType(retry.FixedDelay),
 )
```

# 仅在特定错误时重试
```go
retry.Do(
    func() error {
        return errors.New("special error")
    },
    retry.RetryIf(func(err error) bool {
    return err.Error() == "special error"
    }),
)
```

# 仅在特定错误时终止重试
当使用 retry.Unrecoverable() 包裹错误值后，返回这个错误时不会再进行重试。

```go
retry.Do(
    func() error {
        return retry.Unrecoverable(errors.New("special error"))
    },
)
```
