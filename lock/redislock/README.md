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
