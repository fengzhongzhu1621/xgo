package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Go 1.7 标准库引入 context，中文译作“上下文”，准确说它是 goroutine 的上下文，包含 goroutine 的运行状态、环境、现场等信息
// context 主要用来在 goroutine 之间传递上下文信息，包括：取消信号、超时时间、截止时间、k-v 等。
// context 几乎成为了并发控制和超时控制的标准做法。
// context.Context 类型的值可以协调多个 groutine 中的代码执行“取消”操作，并且可以存储键值对。最重要的是它是并发安全的。
// 与它协作的 API 都可以由外部控制执行“取消”操作

// 在Go 里，我们不能直接杀死协程，协程的关闭一般会用 channel+select 方式来控制。但是在某些场景下，
// 例如处理一个请求衍生了很多协程，这些协程之间是相互关联的：需要共享一些全局变量、有共同的 deadline 等，
// 而且可以同时被关闭。再用 channel+select 就会比较麻烦，这时就可以通过 context 来实现。

// 一句话：context 用来解决 goroutine 之间退出通知、元数据传递的功能。

// context 取值
// Context 指向它的父节点，链表则指向下一个节点。通过 WithValue 函数，可以创建层层的 valueCtx，存储 goroutine 间可以共享的变量。
// 取值的过程，实际上是一个递归查找的过程：
// func (c *valueCtx) Value(key interface{}) interface{} {
//     if c.key == key {
//         return c.val
//     }
//     return c.Context.Value(key)
// }
// 它会顺着链路一直往上找，比较当前节点的 key 是否是要找的 key，如果是，则直接返回 value。否则，一直顺着 context 往前，最终找到根节点（一般是 emptyCtx），直接返回一个 nil。所以用 Value 方法的时候要判断结果是否为 nil。
// 因为查找方向是往上走的，所以，父节点没法获取子节点存储的值，子节点却可以获取父节点的值。
// WithValue 创建 context 节点的过程实际上就是创建链表节点的过程。两个节点的 key 值是可以相等的，但它们是两个不同的 context 节点。查找的时候，会向上查找到最后一个挂载的 context 节点，
// 也就是离得比较近的一个父节点 context。所以，整体上而言，用 WithValue 构造的其实是一个低效率的链表。

func process(ctx context.Context) (string, bool) {
	traceId, ok := ctx.Value("trace_id").(string)
	if ok {
		return traceId, ok
	} else {
		return traceId, ok
	}
}

func TestTransferData(t *testing.T) {
	rootCtx := context.Background()
	_, ok := process(rootCtx)
	assert.Equal(t, ok, false)

	// 创建子上下文
	childCtx := context.WithValue(rootCtx, "trace_id", "123456")
	expect := childCtx.Value("trace_id").(string)
	assert.Equal(t, expect, "123456")

	// 父上下文无法获取子上下文的数据
	expect2 := rootCtx.Value("trace_id")
	assert.Nil(t, expect2)
}

func Perform(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// * 被取消，直接返回
			// * 判断超时是否结束
			// 若不调用cancel函数，到了原先创建Contetx时的超时时间，它也会自动调用cancel()函数，即会往子Context的Done通道发送消息
			return
		case <-time.After(time.Second):
		}
	}
}

func TestCancel(t *testing.T) {
	// 返回一个子Context和一个取消函数CancelFunc
	childCtx, cancel := context.WithTimeout(context.Background(), time.Hour)
	go Perform(childCtx)

	cancel()
}

// gen 生成无限整数的协程
func gen(ctx context.Context) <-chan int {
	ch := make(chan int)
	go func() {
		var n int
		for {
			select {
			case <-ctx.Done():
				// 防止死循环
				return
			case ch <- n:
				n++
				time.Sleep(time.Second)
			}
		}
	}()
	return ch
}

func TestStopLoop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // 避免其他地方忘记 cancel，且重复调用不影响
	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			cancel()
			// 增加一个 context，在 break 前调用 cancel 函数，取消 goroutine。
			// gen 函数在接收到取消信号后，直接退出，系统回收资源。
			break
		}
	}
}
