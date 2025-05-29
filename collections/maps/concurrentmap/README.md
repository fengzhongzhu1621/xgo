# orcaman/concurrent-map/v2
https://github.com/orcaman/concurrent-map

## 特性
* 分片锁机制:将数据分散到多个分片（shard），每个分片独立加锁，减少全局锁竞争。
* 默认分片数为32，可通过构造函数调整。
* 键的分配使用哈希函数（如 FNV-1a）取模分片数，确保均匀分布。

## 与sync.Map的对比

| 特性 | concurrent-map/v2 | sync.Map |
| --- | --- | --- |
| 锁机制 | 分片锁(细粒度) | 全局锁+无锁读 |
| 键值类型 | 泛型支持(任意类型) | 固定为 interface{} |
| API友好性 | 提供 Set/Get等方法 | 需使用 Load/ Store |
| 适用场景 | 高频写操作 | 读多写少 |
| 性能 | 高并发下更优 | 简单场景下足够 |

* 在并发读写场景下，性能显著优于sync.Map（尤其在写操作多时）。
* 提供基准测试（见仓库benchmarks/目录）。

## 语法
* Set(key K, value V): 插入或更新键值对。
* Get(key K) (V, bool): 获取值，返回是否存在。
* Remove(key K): 删除键值对。
* Count() int: 返回总元素数。
* Iter() <-chan Tuple: 返回通道遍历所有键值对。
* Keys() []K: 获取所有键的切片。

# sync.map
设计哲学就是“优化读性能”。它内部用两层 map 管理数据，一层是读 map，一层是脏 map。只要数据没被改，它就直接从读 map 取数据

```go
var m sync.Map

func main() {
    m.Store("foo", "bar")       // 写入
    val, ok := m.Load("foo")    // 读取
    if ok {
        fmt.Println(val) // 输出：bar
    }

    m.Delete("foo")             // 删除
}
```
## 场景
* 配置缓存、字典数据、路由表等读远大于写的场景。
* 多 goroutine 同时读取共享数据，偶尔有写操作。
* 高并发下为了避免锁竞争，希望提升读效率。

## 不适合
* 写操作频繁，比如需要频繁更新、删除的缓存（比如 TTL 缓存），会触发频繁的 map 替换和锁竞争，性能反而可能拉跨。
* 对数据结构操作需要复杂逻辑（比如顺序遍历、批量更新）的时候，sync.Map 的 API 不支持这些复杂操作。

## 遍历

key 和 value 类型都是 any，需要类型断言
```go
m.Range(func(key, value any) bool {
    fmt.Printf("%v: %v\n", key, value)
    return true // 返回 false 会终止遍历
})
```
