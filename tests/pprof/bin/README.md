# 1. 运行
```shell
brew install graphviz
dot -V
# dot - graphviz version 12.2.1 (20241206.2353)
```
```shell
go run main.go
```

# 2. 访问页面
http://localhost:6060/debug/pprof/


# 3. profile
探测各函数对 cpu 的占用情况。

cpu 分析是在一段时间内进行打点采样，通过查看采样点在各个函数栈中的分布比例，以此来反映各函数对 cpu 的占用情况.

点击页面上的 profile 后，默认会在停留 30S 后下载一个 cpu profile 文件. 通过交互式指令打开文件后，查看 cpu 使用情况

## 3.1 命令行显示
```go
go tool pprof profile
```

```shell
File: main
Type: cpu
Time: May 21, 2025 at 9:59pm (CST)
Duration: 30.04s, Total samples = 14.35s (47.76%)
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 14320ms, 99.79% of 14350ms total
Dropped 17 nodes (cum <= 71.75ms)
      flat  flat%   sum%        cum   cum%
   13610ms 94.84% 94.84%    14320ms 99.79%  github.com/wolfogre/go-pprof-practice/animal/felidae/tiger.(*Tiger).Eat
     710ms  4.95% 99.79%      710ms  4.95%  runtime.asyncPreempt
         0     0% 99.79%    14320ms 99.79%  github.com/wolfogre/go-pprof-practice/animal/felidae/tiger.(*Tiger).Live
         0     0% 99.79%    14330ms 99.86%  main.main
         0     0% 99.79%    14330ms 99.86%  runtime.main
(pprof)
```

14320ms 采样点大约覆盖的时长

* flat 某个函数执行时长（只聚焦函数本身，剔除子函数部分）,Tiger.Eat 这个方法本身的调用时长
* flat% 某个函数执行时长百分比（只聚焦函数本身，剔除子函数部分）
* sum%：某个函数及其之上父函数的总时长占比
* cum：某个函数及其子函数的总调用时长，Tiger.Eat 加上其调用子函数 runtime.asyncPreempt 的总时长
* cum%：某个函数及其子函数的调用时长在总时长中的占比

```go
func (t *Tiger)  Eat() {
  log.Println(t.Name(),  "eat")
  loop :=  10000000000
  for  i :=  0; i < loop; i++ {
    // do nothing 通过 for 循环大量空转打满 CPU
  }
}
```
runtime.asyncPreempt 花费了大约 710ms 的时间，但是这一点在代码中并没有体现，这又是怎么回事呢？

关于 goroutine 超时抢占机制的设定：

* 监控线程：在 go 进程启动时，会启动一个 monitor 线程，作为第三方观察者角色不断轮询探测各 g 的执行情况，
对于一些执行时间过长的 g 出手干预

  - 协作式抢占：当 g 在运行过程中发生栈扩张时（通常由函数调用引起），则会触发预留的检查点逻辑，
  查看自己若是因为执行过长而被 monitor 标记，则会主动让渡出 m 的执行权在 Tiger.Eat 方法中，
  由于只是简单的 for 循环空转无法走到检查点，因此这种协作式抢占无法生效
  - 非协作式抢占：在 go 1.14 之后，启用了基于信号量实现的非协作抢占机制.
    Monitor 探测到 g 超时会发送抢占信号，g 所属 m 收到信号后，
    会修改 g 的 栈程序计数器 pc 和栈顶指针 sp 为其注入 asyncPreempt 函数.
    这样 g 会调用该函数完成 m 执行权的让渡

## 3.2 图形化显示
```shell
go tool pprof -http=:8082 profile
```

## 3.3 火焰图
VIEW -> Flame Graph

# 4. heap
探测内存分配情况

http://localhost:6060/debug/pprof/heap?debug=1

在页面的路径中能看到 debug 参数，如果 debug = 1，则将数据在页面上呈现；
如果将 debug 设为 0，则会将数据以二进制文件的形式下载，并支持通过交互式指令或者图形化界面对文件内容进行呈现

第一行
```
heap profile: 3: 1409290240 [13: 1961889792] @ heap/1048576
```

* 3—活跃对象个数
* 1409290240—活跃对象大小（单位 byte，总计约 1.4G）
* 13—历史至今所有对象个数
* 1961889792—历史至今所有对象总计大小（byte）
* 1048576—内存采样频率（约每 M 采样一次）

```
1: 1291845632 [1: 1291845632] @ 0x104303b48 0x1043033b8 0x104303cc0 0x10410938c 0x10413ca24
#        0x104303b47        github.com/wolfogre/go-pprof-practice/animal/muridae/mouse.(*Mouse).Steal+0xf7        /Users/xxx/projects/go-pprof-practice/animal/muridae/mouse/mouse.go:60
#        0x1043033b7        github.com/wolfogre/go-pprof-practice/animal/muridae/mouse.(*Mouse).Live+0x47        /Users/xxx/projects/go-pprof-practice/animal/muridae/mouse/mouse.go:25
#        0x104303cbf        main.main+0xbf                                                                        /Users/xxx/projects/go-pprof-practice/main.go:31
#        0x10410938b        runtime.main+0x2bb 
```
对应为某个函数栈中的信息：
* 1-该函数栈上当前存活的对象个数
* 1291845632-当前存活对象总大小（byte）
* [] 内的内容也表示历史至今，不再赘述

# 5. block
探测阻塞情况 （包括 mutex、chan、cond 等）

http://localhost:6060/debug/pprof/block?debug=1

查看某个 goroutine 陷入 waiting 状态（被动阻塞，通常因 gopark 操作触发，比如因加锁、读chan条件不满足而陷入阻塞）的触发次数和持续时长.

pprof 默认不启用 block 分析，若要开启则需要进行如下设置：
```go
runtime.SetBlockProfileRate(1)
```

此处的入参能够控制 block 采样频率：
* 1：始终采用
* <=0：不采样
* 若 > 1：当阻塞时长(ns)大于该值则采样，否则有阻塞时长/rate的概率被采样

```go
--- contention:
cycles/second=999999999
                                                           /usr/local/go/src/runtime/proc.go:267
3002910915 3 @ 0x100052224 0x10027e9e4 0x10027e5d8 0x10027fb00 0x10008538c 0x1000b8a24
#        0x100052223        runtime.chanrecv1+0x13                                                                /usr/local/go/src/runtime/chan.go:442
#        0x10027e9e3        github.com/wolfogre/go-pprof-practice/animal/felidae/cat.(*Cat).Pee+0xa3        /Users/xxx/projects/go-pprof-practice/animal/felidae/cat/cat.go:39
#        0x10027e5d7        github.com/wolfogre/go-pprof-practice/animal/felidae/cat.(*Cat).Live+0x37        /Users/xxx/projects/go-pprof-practice/animal/felidae/cat/cat.go:19
#        0x10027faff        main.main+0xbf                                                                        /Users/xxx/projects/go-pprof-practice/main.go:31
#        0x10008538b        runtime.main+0x2bb    • c
```

* cycles/second=999999999——是每秒钟对应的cpu周期数. pprof在反映block时长时，以cycle为单位
* 3002910915——阻塞的cycle数. 可以换算成秒：3003307959/1000002977 ≈ 3S
* 3——发生的阻塞次数

# 6. mutex
探测互斥锁占用情况

http://localhost:6060/debug/pprof/mutex?debug=1

某个 goroutine 持有锁的时长（mutex.Lock -> mutex.Unlock 之间这段时间），且只有在存在锁竞争关系时才会上报这部分数据.
pprof 默认不开启 mutex 分析，需要显式打开开关：

```go
runtime.SetMutexProfileFraction(1)
```

入参控制的是 mutex 采样频率：
* 1——始终进行采样
* 0——关闭不进行采样
* <0——不更新这个值，只是把之前设的值结果读出来
* 若 >1 ——有 1/rate 的概率下的事件会被采样

```go
--- mutex:
cycles/second=999999999
sampling period=1
4007486874 4 @ 0x1024e24d4 0x1024e2495 0x10231ca24
#        0x1024e24d3        sync.(*Mutex).Unlock+0x73                                                                /usr/local/go/src/sync/mutex.go:223
#        0x1024e2494        github.com/wolfogre/go-pprof-practice/animal/canidae/wolf.(*Wolf).Howl.func1+0x34        /Users/xxx/projects/go-pprof-practice/animal/canidae/wolf/wolf.go:58

```
* 1000002767——每秒下的 cycle 数
* 4007486874——持有锁的 cycle 总数
* 4——采样了 4 次

# 7. goroutine：探测协程使用情况

http://localhost:6060/debug/pprof/goroutine?debug=1

```go
goroutine profile: total 173
150 @ 0x100a017e8 0x100a315ac 0x100bfa578 0x100a34a24
#        0x100a315ab        time.Sleep+0x10b                                                                        /usr/local/go/src/runtime/time.go:195
#        0x100bfa577        github.com/wolfogre/go-pprof-practice/animal/canidae/wolf.(*Wolf).Drink.func1+0x27        /Users/xxx/projects/go-pprof-practice/animal/canidae/wolf/wolf.go:34

15 @ 0x100a017e8 0x100a315ac 0x100bfb6f0 0x100a34a24
#        0x100a315ab        time.Sleep+0x10b                                                                        /usr/local/go/src/runtime/time.go:195
#        0x100bfb6ef        github.com/wolfogre/go-pprof-practice/animal/muridae/mouse.(*Mouse).Pee.func1+0x2f        /Users/xxx/projects/go-pprof-practice/animal/muridae/mouse/mouse.go:43先
```

* total 173——总计有 173 个 goroutine
* 然后找到创建goroutine 数量较大的方法，分别是150 和 15 个
