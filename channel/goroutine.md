
# 调度对象：进程和线程

时间片轮转：多个程序快速切换使用 CPU，这种切换策略被称为 调度。

调度器调度的是什么？

* 执行中的不同程序 (进程)
* 进程的子集，CPU 利用的基本单位 (线程)

然而，线程切换会带来一定的开销。

# goroutine

用户级线程由运行时系统 (用户级库) 管理，其切换开销几乎与函数调用相同。内核对用户级线程一无所知，将其视为单线程进程。
* 在 Go 中，用户级线程被称为 goroutine。
* goroutine 是 Go 运行时管理的轻量级线程，逻辑上代表着执行线程。

# 协作式的抢占式调度

在1.14 版本之前，程序只能依靠 Goroutine 主动让出 CPU 资源才能触发调度

* 某些 Goroutine 可以⻓时间占用线程，造成其它 Goroutine 的饥饿
* 垃圾回收需要暂停整个程序(Stop-the-world，STW)，最⻓可能需要几分钟的时间，导致整个程序无法工作。

# GMP

G 是 Goroutine 执行的实体，M 是 Goroutine 的承载者，P 是调度器。

## P
* P 的数量肯定不是越小越好，所以官方默认值就是 CPU 的核心数; 逻辑处理器 P 的数量始终固定 (默认为当前进程可用的逻辑 CPU 数量)
* 在任何情况下，Go运行时并行执行(注意，不是并发)的 goroutines 数量是小于等于 P 的数量的
* P 代表处理器，可以看作是在线程上运行的本地调度器。
* Go 运行时将首先根据机器的逻辑 CPU 数量 (或根据请求) 创建固定数量的逻辑处理器 P。

每个 Goroutine 都是由 Go 运行时调度器（Scheduler）进行调度的。调度器负责将 Goroutine 转换成线程上的执行上下文，并在多个线程之间分配 Goroutine 的执行。

Go 调度器的调度策略是基于协作式调度的。 也就是说，调度器会在 Goroutine 主动让出执行权（例如在 I/O 操作、channel 操作、time.Sleep() 等操作中）时，将 CPU 的执行权转交给其他 Goroutine。这种调度策略可以保证 Goroutine 之间的调度是非常轻量级的。

Goroutine 的调度是非确定性的，也就是说，Goroutine 之间的调度是不可预测的。这种调度策略可以保证 Goroutine 的执行具有随机性，可以充分利用多核 CPU 的性能。在实现上采用的是用户态调度，不需要进行内核态和用户态之间的切换，从而可以更快地切换和调度多个 Goroutine。相比之下，传统的线程需要占用更多的资源和时间，因此在多并发的情况下，Go 的 Goroutine 会更加高效。

## G
* 每个 goroutine (G) 都将在分配给逻辑 CPU (P) 的操作系统线程 (M) 上运行。
* 每个协程至少需要消耗 2KB 的空间，那么假设计算机的内存是 2GB，那么至多允许 2GB/2KB = 1M 个协程同时存在。

## M
表示 Machine（即操作系统线程）

# 基于信号的抢占式调度
sysmon 也叫监控线程，变动的周期性检查

Sysmon 操作系统线程，它会定期轮询网络，如果超过 10 毫秒没有轮询，则会将准备好的 G 添加到队列。

* 释放闲置超过 5 分钟的 span 物理内存
* 如果超过 2 分钟没有垃圾回收，强制执行
* 将⻓时间未处理的 netpoll 添加到全局队列
* 向⻓时间运行的 G 任务发出抢占调度(超过 10ms 的 g，会进行 retake)
* 收回因 syscall ⻓时间阻塞的 P

# panic

可以使用recover函数来捕获Goroutine中的panic，并进行相应的处理，如果没有对相应的Goroutine 进行异常处理，会导致主线程 panic 。

# 底层实现
go 关键字启动后编译器器会通过cmd/compile/internal/gc.state.stmt和cmd/compile/internal/gc.state.call 两个方法将该关键字转换成runtime.newproc函数调用。

启动一个新的 Goroutine 来执行任务时，会通过 runtime.newproc 初始化一个 g 来运行协程。

Goroutine 在 Go 运行时（runtime）系统中可以有以下 9 种状态
* Gidle：Goroutine 处于空闲状态，即没有被创建或者被回收；
* Grunnable：Goroutine 可以被调度器调度执行，但是还未被选中执行；
* Grunning：Goroutine 正在执行中，被赋予了M和P的资源；
* Gsyscall：Goroutine 发起了系统调用，进入系统调用阻塞状态；
* Gwaiting：Goroutine 被阻塞等待某个事件的发生，比如等待 I/O、等待锁、等待 channel 等；
* Gscan：GC正在扫描栈空间
* Gdead：没有正在执行的用户代码
* Gcopystack：栈正在被拷贝，没有正在执行的代码
* Gpreempted：Goroutine 被抢占，即在运行过程中被调度器中断。等待重新唤醒


Goroutine 的调度时机一般有以下几种情况
* 当前 Goroutine 主动让出执行权时，调度器会将 CPU 的执行权转交给其他 Goroutine。
* 当前 Goroutine 执行的时间超过了 Go 运行时所设置的阈值时，调度器会将当前 Goroutine 暂停，将 CPU 的执行权转交给其他 Goroutine。
* 当前 Goroutine 进行 I/O 操作、channel 操作或者其他系统调用时，调度器会将当前 Goroutine 暂停，将 CPU 的执行权转交给其他 Goroutine。
* 当前 Goroutine 被阻塞在同步原语（例如 sync.Mutex）时，调度器会将当前 Goroutine 暂停，将 CPU 的执行权转交给其他 Goroutine。

