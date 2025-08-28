# 1. 日志级别
* Debug：调试信息。
* Info：一般信息。
* Warn：警告信息。
* Error：错误信息。
* Fatal：致命错误， 程序将退出。 会记录日志后调用 os.Exit(1) 退出程序。
* Panic：严重错误，程序将 panic。录日志后调用 panic()，使程序进入 panic 状态。

# 2. 日志格式
在默认情况下，logrus 使用 TextFormatter 格式输出日志。

## 2.1 TextFormatter（默认）
人类可读的文本格式。

## 2.2 JSONFormatter
结构化日志输出，适用于机器解析。


# 3. 日志字段（Structured Logging）
支持结构化日志，可以通过 WithField 或 WithFields 方法添加额外的字段（如用户 ID、请求 ID 等）到日志中

```
log.WithField("user", "John").Info("User logged in")
输出：{"level":"info","msg":"User logged in","user":"John","time":"2025-02-07T00:00:00+00
```

```
log.WithFields(logrus.Fields{
"user": "John",
"id":   123,
}).Info("User logged in")
输出：{"level":"info","msg":"User logged in","id":123,"user":"John","time":"2025-02-07T00
```

# 4. 自定义hook
必须实现 logrus.Hook 接口，接口定义了 Levels() 和 Fire() 方法。

