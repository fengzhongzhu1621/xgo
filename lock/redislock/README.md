# 简介
```go
go get github.com/bsm/redislock

```
bsm/redislock 是一个用于 在 Go 语言中实现基于 Redis 的分布式锁（distributed lock） 的开源库，由 bsm 组织开发和维护。

* 基于 Redis 单实例实现分布式锁
* 遵循 Redis 官方推荐的 SET resource-name value NX PX max-lock-time 模式
* 简洁、高性能、线程安全
* 支持锁自动过期
* 支持锁续约（可选）
* 没有依赖复杂的 Redlock 多主节点算法，更适合大多数实际应用

# 缺点
* 不是 Redlock 多 Redis 节点实现，不适合对高可用性强依赖的分布式系统。
* 更适合内部服务、微服务之间的协调、任务调度等不要求高强一致性的场景。


# 语法

## Valid()
锁是否仍然有效ok, err := lock.Valid(ctx)
如果 ok == false，表示锁已经过期或被其他客户端释放/抢占。

## 锁续约
```go
err := lock.Refresh(ctx, 10*time.Second, nil)
```

* RetryStrategy：设置重试策略（如使用 backoff 重试机制）
* Metadata：设置锁的元数据（可选）

```go
// LinearBackoff 是一种​​线性退避策略​​，表示：
// 如果第一次获取锁失败，等待 500ms 后重试；
// 如果第二次仍然失败，再等待 500ms 后重试；
// 依此类推，直到达到最大重试次数（默认无限制，除非设置 Tries）。

opts := &redislock.Options{
    RetryStrategy: redislock.LinearBackoff(500 * time.Millisecond),
}
lock, err := locker.Obtain(ctx, "my-key", 10*time.Second, opts)
```
