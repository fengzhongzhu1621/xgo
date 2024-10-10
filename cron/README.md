# cron 表达式
https://cloud.tencent.com/document/product/583/9708

# 任务调度原理

维护一个任务列表，其中每个任务包含两部分
* 执行时间计划：通过 cron 表达式定义，表示任务何时需要执行。
* 执行的函数或作业：每个任务对应一个具体的操作，这可以是一个函数（通过 AddFunc）或者一个实现了 Job 接口的作业。

调度器启动后，会进入一个定时循环。在循环中，调度器会定期检查当前时间，遍历任务列表，逐一判断这些任务的执行时间是否已到。若某个任务的时间已到，调度器便会触发该任务的执行。

每次遍历时，cron 会依据任务的计划和当前时间进行匹配。如果匹配成功，任务将被放入调度队列中进行执行。


# 时间计算
调度器将根据这个时间计划生成下一个触发时间

当调度器的定时器触发时，会检查下次执行时间，如果当前时间符合执行条件，则将任务标记为可执行。

# 线程安全
使用了内部锁机制来保证任务添加、删除、修改以及任务执行的并发安全。

锁机制确保了在调度过程中，多个 Goroutine 并发执行任务时，不会出现竞争条件（Race Condition）。这一点非常重要，因为多个任务可能会同时触发并执行，如果不进行适当的锁定，任务状态可能会发生不一致的现象。

* 在任务的增删改时，确保任务列表不被其他 Goroutine 读取或修改。
* 在执行任务时，保证调度器不会因并发修改任务列表而导致崩溃或任务丢失。

# 防止任务堆积
WithChain 的 cron.SkipIfStillRunning 函数来跳过任务。这个功能可以防止任务的堆积，确保每个任务只在上一次执行完毕后才会被再次触发。

```go
c := cron.New(
    cron.WithSeconds(),                // 支持秒级调度
    cron.WithChain(
        cron.SkipIfStillRunning(nil),  // 如果上一次任务未完成，则跳过
    ),
)
```


# 暂停
没有显式的暂停功能，如果需要手动暂停任务，可以通过业务逻辑或状态控制任务的执行。


# 避免 panic 导致调度器崩溃

```go
c := cron.New(
    cron.WithChain(
        cron.Recover(nil), // 捕获任务中的 panic，防止调度器崩溃
    ),
)
```

# 动态添加删除任务

```go
id, _ := c.AddFunc("@daily", dailyTask)
c.Remove(id) // 删除任务
```

```go
c.AddFunc("@every 10s", healthCheck)
```

# 获取任务列表

```go
entries := c.Entries()
for _, entry := range entries {
    fmt.Println("任务ID：", entry.ID)
}
```

# 时区

```go
loc, _ := time.LoadLocation("Asia/Shanghai")
c := cron.New(cron.WithLocation(loc))
```

# 停止调度

```go
defer c.Stop()
```