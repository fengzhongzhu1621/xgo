# 延时队列
队列是存储消息的载体，延时队列存储的对象是延时消息。
延时消息，是指消息被发送以后，并不想让消费者立刻获取，而是等待特定的时间后，消费者才能获取这个消息进行消费。

## 实现
### Redis zset
* 将延时消息的score 设置为到期的时间戳，消息内容序列化为value，调用zdd命令将此条延迟消息保存在key为delayqueue 的zset中
* 另起线程，循环从 delayqueue 中获取score小于等于当前时间戳的消息元素（zrangebyscore命令）
* zrem删除获取到的元素

### Kafka实现延时队列
* 在发送延时消息时，先将消息投递到延时队列（delay_topic）中（headers中设置延时时间，timestamp存消息发送初始发送时间戳）
* 定义一个服务去消费延时队列中的消息，将满足条件的消息再投递到目标队列（target_topic）中。按照延时等级来划分 delay_topic，如设定5s，10s，30s，1min，5min，30min，1h，2h这些递增的延时等级，延时消息只支持这些等级内的延时，然后延时的消息按照延时时间投递到不同等级的topic中
