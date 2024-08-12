
# 调度对象：进程和线程

时间片轮转：多个程序快速切换使用 CPU，这种切换策略被称为 调度。

调度器调度的是什么？

* 执行中的不同程序 (进程)
* 进程的子集，CPU 利用的基本单位 (线程)

然而，线程切换会带来一定的开销。

# 内核态和用户态
* ***内核态线程***：由操作系统管理和调度，CPU只负责处理内核态线程。
* ***用户态线程***：由用户程序管理，需绑定到内核态线程上执行，协程即为用户态线程的一种。


# goroutine

用户级线程由运行时系统 (用户级库) 管理，其切换开销几乎与函数调用相同。内核对用户级线程一无所知，将其视为单线程进程。
* 在 Go 中，用户级线程被称为 goroutine。
* goroutine 是 Go 运行时管理的轻量级线程，逻辑上代表着执行线程。
* Goroutine的调度由Go语言的运行时（runtime）负责，而不是操作系统。Go运行时在用户态进行调度，避免了频繁的上下文切换带来的开销，使得调度更加高效。
* 在Go程序中，main函数本身也是一个Goroutine，称为主Goroutine。当主Goroutine结束时，所有其他Goroutine也会随之终止。因此，需要确保主Goroutine等待所有子Goroutine执行完毕。
* 线程是运行Goroutine的实体，而调度器的功能是将可运行的Goroutine分配到工作线程上。

# 协作式的抢占式调度

在1.14 版本之前，程序只能依靠 Goroutine 主动让出 CPU 资源才能触发调度

* 某些 Goroutine 可以⻓时间占用线程，造成其它 Goroutine 的饥饿
* 垃圾回收需要暂停整个程序(Stop-the-world，STW)，最⻓可能需要几分钟的时间，导致整个程序无法工作。

# GMP

G 是 Goroutine 执行的实体，M 是 Goroutine 的承载者，P 是调度器。

## P（Processor 协程调度器）
* P代表执行上下文（Processor）。P管理着可运行的Goroutine队列，并负责与M进行绑定。P的数量决定了可以并行执行的Goroutine的数量。Go运行时会根据系统的CPU核数设置P的数量。
* 包含了运行 goroutine 的资源。如果线程想运行 goroutine，必须先获取 P，P 中还包含了可运行的 G 队列。
* P 的数量肯定不是越小越好，所以官方默认值就是 CPU 的核心数; 逻辑处理器 P 的数量始终固定 (默认为当前进程可用的逻辑 CPU 数量)
* 所有的 P 都在程序启动时创建，并保存在数组中，最多有 GOMAXPROCS (可配置)个。
P 的数量由 环境变量 $GOMAXPROCS 或者是由 runtime 的方法 GOMAXPROCS() 决定，当 P 的数量 n 确定以后，运行时系统会根据这个数量创建 n 个 P。
* 在任何情况下，Go运行时并行执行(注意，不是并发)的 goroutines 数量是小于等于 P 的数量的
* P 代表处理器，可以看作是在线程上运行的本地调度器。
* Go 运行时将首先根据机器的逻辑 CPU 数量 (或根据请求) 创建固定数量的逻辑处理器 P。
* 每个 Goroutine 都是由 Go 运行时调度器（Scheduler）进行调度的。调度器负责将 Goroutine 转换成线程上的执行上下文，并在多个线程之间分配 Goroutine 的执行。
* Go 调度器的调度策略是基于协作式调度的。 也就是说，调度器会在 Goroutine 主动让出执行权（例如在 I/O 操作、channel 操作、time.Sleep() 等操作中）时，将 CPU 的执行权转交给其他 Goroutine。这种调度策略可以保证 Goroutine 之间的调度是非常轻量级的。
* Goroutine 的调度是非确定性的，也就是说，Goroutine 之间的调度是不可预测的。这种调度策略可以保证 Goroutine 的执行具有随机性，可以充分利用多核 CPU 的性能。在实现上采用的是用户态调度，不需要进行内核态和用户态之间的切换，从而可以更快地切换和调度多个 Goroutine。相比之下，传统的线程需要占用更多的资源和时间，因此在多并发的情况下，Go 的 Goroutine 会更加高效。

## G（Goroutine）
* Goroutine是Go语言中的协程，代表一个独立的执行单元。Goroutine比线程更加轻量级，启动一个Goroutine的开销非常小。Goroutine的调度由Go运行时在用户态进行。
* 每个 goroutine (G) 都将在分配给逻辑 CPU (P) 的操作系统线程 (M) 上运行。
* 每个协程至少需要消耗 2KB 的空间，那么假设计算机的内存是 2GB，那么至多允许 2GB/2KB = 1M 个协程同时存在。
* 全局G队列：存放等待运行的G。当 P 的本地队列为空时，优先从全局队列获取，如果全局队列为空时则通过 work stealing 机制从其他P的本地队列偷取 G。
* P的本地G队列：存放不超过256个G。当新建协程时优先将G存放到本地队列，如果队列满了，会将本地队列的一半 G 和新创建的 G 打乱顺序，一起放入全局队列。
* G0：G0 是每次启动一个 M 都会第一个创建的 gourtine，G0 仅用于负责调度的 G，G0 不指向任何可执行的函数, 每个 M 都会有一个自己的 G0。在调度或系统调用时会使用 G0 的栈空间, 全局变量的 G0是 M0 的 G0。
* 新创建的协程优先保存在P的本地G队列，如果本地队列满了，会将P本地队列中的一半G打乱顺序移入全局队列。
* 唤醒：创建G时，当前运行的G会尝试唤醒其他PM组合执行。若唤醒的M绑定的P本地队列为空，M会尝试从全局队列获取G。
* 偷取(work stealing 机制)：如果P的本地队列和全局队列都为空，会从其他P偷取一半G到自己的本地队列执行。当本线程无可运行的 G 时，尝试从其他线程绑定 P 的队列中偷取 G，而不是销毁线程。
* 抢占：在 coroutine 中要等待一个协程主动让出 CPU 才执行下一个协程。在Go中，一个 goroutine 最多占用 CPU 10ms，防止其他 goroutine 被饿死。
* 自旋线程会占用CPU时间，但创建销毁线程也消耗CPU时间。系统最多有GOMAXPROCS个自旋线程，其他线程在休眠M队列里。
* 系统调用：当G进行系统调用时进入内核态被阻塞，M会释放绑定的P，把P转移给其他空闲的M执行。当系统调用结束，GM会尝试获取一个空闲的P。
* 阻塞处理(hand off 机制)：当G因channel或network I/O阻塞时，不会阻塞M。超过10ms时，M会寻找其他可运行的G。
* 公平性：调度器每调度61次时，会尝试从全局队列中取出待运行的Goroutine来运行。如果没有找到，就去其他P偷一些Goroutine来执行。
* P 的数量由 GOMAXPROCS 设置，最多有 GOMAXPROCS 个线程分布在多个 CPU上 同时运行。GOMAXPROCS 也限制了并发的程度，比如 GOMAXPROCS = 核数/2，则最多利用了一半的 CPU 核进行并行。


## M（Machine）
表示 Machine（即操作系统线程）
* M代表操作系统的线程。M负责实际执行Go代码。一个M可以执行多个Goroutine，但同一时间只能执行一个Goroutine。M与操作系统的线程直接对应，Go运行时通过M来利用多核CPU的并行计算能力。
* 当本线程因G进行系统调用等阻塞时，线程会释放绑定的P，把P转移给其他空闲的M执行。
* 线程想运行任务就得获取 P，从 P 的本地队列获取 G，P 队列为空时，M 也会尝试从全局队列拿一批 G 放到 P 的本地队列，或从其他 P 的本地队列偷一半放到自己 P 的本地队列。M 运行 G，G 执行之后，M 会从 P 获取下一个 G，不断重复下去。
* go 程序启动时，会设置 M 的最大数量，默认 10000，但是内核很难支持这么多的线程数。
* 当没有足够的 M 来关联 P 并运行其中的可运行的 G 时，比如所有的 M 此时都阻塞住了，而 P 中还有很多就绪任务，就会去寻找空闲的 M，而没有空闲的，就会去创建新的 M。
* M 与 P 的数量没有绝对关系，一个 M 阻塞，P 就会去创建或者切换另一个 M，所以，即使 P 的默认数量是 1，也有可能会创建很多个 M。
* M0：启动程序后的第一个主线程，负责执行初始化操作和启动第一个Goroutine，此后与其他M一样。这个 M 对应的实例会在全局变量 runtime.m0 中，不需要在 heap 上分配，M0 负责执行初始化操作和启动第一个 G， 在之后 M0 就和其他的 M 一样了。


# 基于信号的抢占式调度
* sysmon 也叫监控线程，变动的周期性检查
* Sysmon 操作系统线程，它会定期轮询网络，如果超过 10 毫秒没有轮询，则会将准备好的 G 添加到队列。

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

