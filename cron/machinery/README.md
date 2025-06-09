# 简介

go machinery框架类似python中常用celery框架，主要用于 异步任务和定时任务。

```go
go get github.com/RichardKnop/machinery/v2
```

# 特性

* 任务重试机制
* 延迟任务支持
* 任务回调机制
* 任务结果记录
* 支持Workflow模式：Chain，Group，Chord
* 多Brokers支持：Redis, AMQP, AWS SQS(opens new window)
* 多Backends支持：Redis, Memcache, AMQP, MongoDB
