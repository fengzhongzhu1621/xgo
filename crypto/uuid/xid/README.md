# 简介
xid 是一种轻量级的全局唯一 ID 生成方案。长度为 20 个字符，包含 4 字节的 Unix 时间戳、5 字节的全局唯一机器 ID 和 3 字节的计数器。

id (github.com/rs/xid) 是一个基于 MongoDB ObjectID 算法的更轻量级、更高效的全局唯一 ID 生成库，用于生成 12 字节的唯一 ID，
特性：

* 更短：比 UUID 更短（12 字节 vs 16 字节）。
* 无需配置：不像 Snowflake 需要机器 ID。
* 时间排序：可以按时间顺序存储/查询数据。
* 分布式友好：包含 机器 ID 和 进程 ID，确保在多台服务器上仍然唯一。


# 组成
同一秒内，每个进程最多可生成 16,777,216 个唯一 ID。其字符串表示形式采用 base32hex（无填充），长度为 20 个字符

* 时间戳（4 字节）：Unix 时间，精确到秒。
* 机器 ID（3 字节）：唯一标识运行的机器。
* 进程 ID（2 字节）：标识进程，避免冲突。
* 计数器（3 字节）：从随机数开始，每次递增。

# 特性
* 特性唯一性：通过时间戳、机器 ID 和计数器组合，保证全局唯一性。
* 排序性：ID 具备时间排序性，适合检索和排序。
* 性能：生成 ID 速度非常快。
