package channel

// channel的基本特性
// * 类型安全： 可以传递任何类型的数据
// * 缓冲： 可以是带缓冲或无缓冲的
// * 同步： 提供同步机制，可以在数据发送和接受时同步 goroutines
// * 关闭： 可以呗关闭，一旦关闭就不能再发送数据

// 用途
// * 同步： 协调多个 goroutines 的执行
// * 通信： 在多个 goroutines 之间传递数据
// * 并行： 收集并发执行的结果
