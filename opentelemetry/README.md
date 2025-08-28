# 创建嵌套跨度

```go
func parentFunction(ctx context.Context) {
    // 1. 创建一个新的 Span（"parent"），并继承传入的 ctx
    ctx, parentSpan := tracer.Start(ctx, "parent")
    defer parentSpan.End() // 确保函数结束时 Span 被结束

    // 2. 调用 childFunction，传递 ctx（包含 parentSpan 的上下文）
    childFunction(ctx)

    // 3. 更多工作...
    // 当 parentFunction 返回时，parentSpan 自动结束（因为 defer）
}

func childFunction(ctx context.Context) {
    // 1. 创建一个新的 Span（"child"），继承传入的 ctx（包含 parentSpan）
    ctx, childSpan := tracer.Start(ctx, "child")
    defer childSpan.End() // 确保函数结束时 Span 被结束

    // 2. 执行子函数的工作...
    // 当 childFunction 返回时，childSpan 自动结束（因为 defer）
}
```

# 设置跨度属性

属性是作为元数据应用于跨度的键和值，对于聚合、过滤和分组跟踪非常有用。可以在创建跨度时添加属性，也可以在跨度生命周期完成之前的任何其他时间添加属性。

```go
// setting attributes at creation...
ctx, span = tracer.Start(ctx, "attributesAtCreation", trace.WithAttributes(attribute.String("hello", "world")))
// ... and after creation
span.SetAttributes(attribute.Bool("isTrue", true), attribute.String("stringAttr", "hi!"))0

```

# 记录事件
事件是一段人类可读的消息，代表在其生命周期内“发生的事情”。例如，想象一个函数需要对互斥体下的资源进行独占访问。可以在两个点创建事件 - 一次是当我们尝试访问资源时，另一个是当我们获取互斥锁时。

事件的一个有用特征是它们的时间戳显示为从跨度开始的偏移量，使您可以轻松查看它们之间经过了多少时间。

```go
span.AddEvent("Acquiring lock")
mutex.Lock()
span.AddEvent("Got lock, doing work...")
// do stuff
span.AddEvent("Unlocking")
mutex.Unlock()
```

事件也可以有自己的属性
```go
span.AddEvent("Cancelled wait due to external signal", trace.WithAttributes(attribute.Int("pid", 4328), attribute.String("signal", "SIGHUP")))
```

# 设置错误状态及错误日志

可以在跨度上设置状态，通常用于指定跨度正在跟踪的操作中存在错误
```go
result, err := operationThatCouldFail()
if err != nil {
    span.SetStatus(codes.Error, "operationThatCouldFail failed")
}
```

记录错误日志
```go
result, err := operationThatCouldFail()
if err != nil {
    span.SetStatus(codes.Error, "operationThatCouldFail failed")
    span.RecordError(err)
}
```

# 使用baggage统一数据存储和传播

在OpenTelemetry中，Baggage是在跨度之间传递的上下文信息。它是一个键值存储，与跟踪中的跨度上下文一起驻留，使值可用于该跟踪中创建的任何跨度

```go
func setBaggage(ctx context.Context) context.Context {
    // baggage统一数据存储和传播
    m1, _ := baggage.NewMember("data1", "data1-value")
    m2, _ := baggage.NewMember("data2", "data2-value")
    b, _ := baggage.New(m1, m2)
    return baggage.ContextWithBaggage(ctx, b)
}


// baggage统一数据存储和传播
childSpan.SetAttributes(attribute.Key(baggage.FromContext(ctx).Member("data1").Key()).String(baggage.FromContext(ctx).Member("data1").Value()))
// baggage统一数据存储和传播
```
