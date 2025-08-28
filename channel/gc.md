# 1. 简介
自 Go 1.5 起，Go 使用并发标记-清除（concurrent mark-sweep）算法，结合“三色标记”模型与 Yuasa 写屏障。

* Go 1.12：优化了STW阶段的扫描算法
* Go 1.14：实现了真正的并行标记
* Go 1.18：进一步减少内存占用

Go GC 会在后台并发地遍历堆内存，标记可达对象，并逐步清除未被引用的内存块。整个回收过程中，Go 追求低延迟、低停顿。

Go运行时启动的系统监控线程负责检测GC触发条件
```go
// $GOROOT/src/runtime/proc.go
func sysmon(){
    for {
        // 每10ms检查一次
        if t :=(gcTrigger{kind: gcTriggerTime, now: now}); t.test(){
            gcStart(gcTrigger{kind: gcTriggerTime, now: now})
        }
        // ...
    }
}
```

* ✅ 并发标记、并发清除
* ✅ 不会移动对象（即 no compaction）
* ✅ 按 span（内存块）分批清扫，减少单次 STW（Stop-the-World）时长

```go
// 内存分配示例
func createLargeSlice(){
    // 分配后不再引用的内存将成为GC目标
    data :=make([]byte,10*1024*1024)// 分配10MB内存
    _= data // 无实际引用
```

# 2. 缺陷
## 2.1. 内存访问低效
GC 的标记阶段会跨对象跳跃，导致 CPU 频繁 cache miss、等待内存，约 35% 的 GC CPU 周期被耗在“等内存”。
这在 NUMA 架构或多核大内存机器上尤为明显。

## 2.2. 缺乏分代收集
Go GC 没有分代机制，所有对象一视同仁，这在高分配率场景下显得笨重。
Pinterest 工程师曾指出，内存压力一旦增大，GC 就会暴增 CPU 消耗，引发延迟激增。

## 2.3. 频繁 GC 带来的 CPU 占用
Twitch 工程团队曾报告：即便在中小堆内存下（<450 MiB），
系统稳态下每秒会触发 8–10 次 GC， 每分钟累计 400–600 次， GC 占用约 30% 的 CPU 时间。
这直接挤压了业务线程的执行空间。

# 3. 系统自动触发
## 3.1 堆内存阈值触发(gcTriggerHeap)
当堆内存分配达到控制器动态计算的阈值时触发。
```sh
触发阈值 = 上次GC后堆大小 × (1 + 内存增长率)
```

内存增长率由`GOGC`环境变量控制(默认100%)
```sh
# （默认）
export GOGC=100
```

在`mallocgc`函数中检查堆大小是否达到阈值
```go
// $GOROOT/src/runtime/malloc.go
func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {
    // 大对象直接触发GC检查
    if size >= maxSmallSize {
        if shouldhelpgc {
            gcStart(gcTrigger{kind: gcTriggerHeap})
        }
    }
// ...
}
```

## 3.2 定时触发(gcTriggerTime)
```go
// $GOROOT/src/runtime/proc.go
const forcegcperiod = 2*60*1e9// 2分钟(纳秒)
```

## 3.3 启动时触发(gcTriggerCycle)
系统启动时完成初始化GC
```go
// $GOROOT/src/runtime/proc.go
func main(){
    // ...
    systemstack(func(){
        gcStart(gcTrigger{kind: gcTriggerCycle, n:1})
    })
}
```

# 4. 手动触发
```go
func manualTrigger(){
    // 业务代码...

    // 关键点：手动触发GC
    runtime.GC()

    // 调试时可添加跟踪
    debug.FreeOSMemory()
}
```


# 5. 完整GC周期

* 标记准备(Mark Setup) - STW
* 并发标记(Concurrent Marking)
* 标记终止(Mark Termination) - STW
* 并发清除(Concurrent Sweeping)

* STW时间：Go 1.14+通常<1ms
* GC频率：受GOGC影响显著
* CPU占用：通常<25%的CPU核心

```go
// 简化的GC启动流程
func gcStart(trigger gcTrigger){
    // 1. 停止所有运行中的Goroutine(STW)
    stopTheWorld()

    // 2. 执行标记阶段
    gcBgMarkStartWorkers()

    // 3. 恢复运行(结束STW)
    startTheWorld()

    // 4. 后台并发清除
    gcSweep()
}
```

# 6. 监控GC行为

```go
// 输出GC统计信息
func printGCStats() {
    // 创建一个每5秒触发一次的定时器
    t := time.NewTicker(5 * time.Second)
    defer t.Stop() // 确保函数退出时停止定时器

    for range t.C { // 每次定时器触发时执行
        var m runtime.MemStats
        runtime.ReadMemStats(&m) // 读取当前内存统计信息

        // 打印GC统计信息
        fmt.Printf(
            "GC次数:%d STW时间:%.3fms 堆大小:%dMB\n",
            m.NumGC, // GC执行的次数
            float64(m.PauseNs[(m.NumGC+255)%256])/1e6, // 最近一次GC的STW时间(毫秒)
            m.HeapAlloc/1024/1024, // 当前堆分配的内存大小(MB)
        )
    }
}
```
