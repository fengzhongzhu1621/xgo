# 简介
* 基于 Redis 的 SETNX 命令实现的分布式锁。SETNX 是一个原子操作，当键不存在时，它会将键设置为指定的值。
* 使用 SETNX 和 EXPIRE 命令组合来实现锁的互斥性和自动释放。
* 使用 RedLock 算法来提高锁的可靠性。（由 Redis 创始人 antirez 提出），该算法旨在在分布式环境中实现安全、健壮的锁。

# 原理
1. 使用多个独立的 Redis 实例。
2. 客户端尝试依次在多数 Redis 节点上设置一个带过期时间的锁。
3. 如果大多数实例成功设置锁，并且总耗时小于锁的过期时间，则认为加锁成功。
4. 释放锁时，会逐个删除这些实例上的锁。

# 特性
## 不支持可重入锁
如果一个协程已经持有了锁，再次尝试加锁将会失败。

## 不支持锁续期
锁会在指定的过期时间后自动释放。

## 性能
* redsync：由于使用 SETNX 和 EXPIRE 命令组合实现锁，性能相对较高。
* redislock：使用 Lua 脚本实现锁，性能相对较低，但在大多数场景下性能仍然可以接受。


# 语法

## 配置
| 配置函数        | 描述                                                         |
|-----------------|--------------------------------------------------------------|
| WithExpiry()    | 设置锁的过期时间。当锁在指定时间内未被释放时，将自动过期，以防止死锁。 |
| WithTries()     | 设置尝试加锁的最大次数。如果在指定次数内未能成功获取锁，将放弃尝试。 |
| WithRetryDelay()| 设置每次重试获取锁的延迟时间。在每次尝试失败后，等待指定的延迟时间再进行下一次尝试。 |
| WithTimeout()   | 设置获取锁的超时时间。在整个获取锁的过程中，如果在指定时间内未能成功获取锁，将放弃尝试。注意，这不是锁本身的过期时间。 |

```go
rs := redsync.New(pools..., redsync.WithRedLockMode())
mutex := rs.NewMutex("lock-name",
    redsync.WithExpiry(10*time.Second),
    redsync.WithTries(5),
    redsync.WithRetryDelay(500*time.Millisecond),
)
```

## 多 Redis 实例支持
Redlock 要求至少多数节点成功加锁才认为锁成功。

```go
pools := []redsync.Pool{
    goredis.NewPool(client1),
    goredis.NewPool(client2),
    goredis.NewPool(client3),
}
rs := redsync.New(pools...)
```
