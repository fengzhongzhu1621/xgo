# Redis Stream
Redis Streams 是 Redis 5.0 引入的一种数据类型，类似于消息队列，但功能更强大
* 支持消息的持久化
* 消费者组 Consumer Group，便于水平扩展；
* 消息确认及重放；
* 高效处理大规模事件。

# 指令
## XAdd
XAdd 是 Redis 的 Stream 添加消息命令。
它会向指定的 Stream 添加一条新消息，并返回消息的 ID（格式如 "1640995200000-0"）
