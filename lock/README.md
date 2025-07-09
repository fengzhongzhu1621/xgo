# redis锁比较
| 特性           | redsync                                      | redislock                                    |
| -------------- | -------------------------------------------- | -------------------------------------------- |
| 实现方式       | 基于 Redis 的 SETNX 命令实现分布式锁         | 基于 Redis 的 Lua 脚本实现分布式锁           |
| 可重入性       | 不支持可重入锁                               | 支持可重入锁                                 |
| 锁续期         | 不支持锁续期                                 | 支持锁续期                                   |
| 性能           | 较高                                         | 相对较低                                     |
| 社区支持       | 较少                                         | 较好                                         |
| 文档和示例     | 较少                                         | 详细                                         |
| 适用场景       | 需要对 RedLock 算法的支持和较高性能的场景   | 需要可重入锁和锁续期功能的场景               |

# sync.Mutex（互斥锁）
同一时间只允许一个 goroutine 获取锁，其他 goroutine 必须等待。
读操作无法并行：即使多个 goroutine 只是读取数据，也必须排队等待锁，导致性能下降（读多写少时瓶颈明显）。
写多读少时性能更好：如果写操作频繁，Mutex 不会比 RWMutex 更差，甚至可能更高效（因为 RWMutex 有额外的读写状态管理开销）。

* 写操作频繁（如计数器、状态标志位）。
* 读写操作都需要严格互斥（如修改共享结构体）。

# sync.RWMutex（读写互斥锁）

## 互斥性，在持有写锁时，任何读/写行为都会被阻塞；且未持有写锁时，任何读行为都不会被阻塞
* 读锁和读锁 - 不互斥
* 读锁和写锁 - 互斥
* 写锁和写锁 - 互斥

## 读写锁获取锁的优先级
* 写优先的策略，可以保证即便在读密集的场景下，写锁也不会饥饿
* 只要有一个写锁申请加锁，那么就会阻塞后续的所有读锁加锁行为（已经获取到读锁的reader不受影响，写锁仍然要等待这些读锁释放之后才能加锁）

## 优点
* 读操作可并行：多个 goroutine 可以同时读数据，提高并发性能（适合读多写少场景）。
* 写操作仍然互斥：保证数据一致性。

## 缺点
* 写操作性能较差：写锁会阻塞所有读锁和写锁，如果写操作频繁，可能导致性能下降。
* 实现复杂度稍高：需要正确管理 RLock() 和 RUnlock()，否则可能死锁或数据竞争。

## 适用场景
* 读多写少（如缓存、配置热更新）。
* 读操作远多于写操作（如数据库查询缓存）。

## 实现
采用了writer-preferring的策略，使用Mutex实现写-写互斥，通过信号量等待和唤醒实现读-写互斥

```go
type RWMutex struct {
  w           Mutex        // 互斥锁，用于写锁互斥
  writerSem   uint32       // writer信号量，读锁RUnlock时释放，可以唤醒等待写加锁的线程
  readerSem   uint32       // reader信号量，写锁Unlock时释放，可以唤醒等待读加锁的线程
  readerCount atomic.Int32 // 所有reader的数量(包括等待读锁和已经获得读锁)
  readerWait  atomic.Int32 // 已经获取到读锁的reader数量，writer需要等待这个变量归0后才可以获得写锁
}
```

读锁加锁-释放行为，开销非常小，仅仅是原子更新readerCount；
当有writer请求写锁时，写锁的加锁-释放行为会重一些，已经获得读锁的reader也不受影响；
后续再请求读锁的reader将被阻塞，直到写锁释放

## 方法
* RLock()：加读锁，允许多个 goroutine 同时读。
* RUnlock()：释放读锁。
* Lock()：加写锁，独占访问，其他 goroutine（包括读和写）必须等待。
* Unlock()：释放写锁。

```go
var (
    rwmu   sync.RWMutex
    cache  map[string]string
)

func getFromCache(key string) (string, bool) {
    rwmu.RLock()
    defer rwmu.RUnlock()
    val, ok := cache[key]
    return val, ok
}

func setToCache(key, value string) {
    rwmu.Lock()
    defer rwmu.Unlock()
    cache[key] = value
}
```

* getFromCache() 是读操作，多个 goroutine 可以同时读取缓存，性能高。
* setToCache() 是写操作，独占锁，保证数据一致性。

## 性能对比
在 读多写少 场景下，RWMutex 比 Mutex 快 10 倍以上。
但在 写多读少 场景下，RWMutex 可能比 Mutex 慢（因为额外的读写状态管理开销）。
