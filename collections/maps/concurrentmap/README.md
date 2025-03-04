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

