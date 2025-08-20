package channel

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
)

// longRunningOperation 子协程，接收到截止信号后，停止运行。
func longRunningOperation(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			// 返回错误信息
			return ctx.Err()
		}
	}
}

func TestWithCancel(t *testing.T) {
	// 创建可取消的Context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("子协程收到上下文取消通知")
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}()

	time.Sleep(3 * time.Second)
	// 取消goroutine，子协成会接收到取消通知
	cancel()
	time.Sleep(2 * time.Second) // 等待 子goroutine 退出
}

func TestWithDeadline(t *testing.T) {
	// 创建带有截止时间的Context
	deadline := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("子协程收到上下文取消通知")
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}()

	time.Sleep(6 * time.Second) // 等待超过截止时间
}

func TestWithTimeout(t *testing.T) {
	// 创建带有超时的 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("子协程收到上下文取消通知")
				return
			default:
				time.Sleep(time.Second)
			}
		}
	}()

	// 超时控制
	time.Sleep(6 * time.Second) // 等待超过超时时间
}

func TestWithValue(t *testing.T) {
	// 创建携带键值对的 Context
	type key int
	const userKey key = 0
	ctx := context.WithValue(context.Background(), userKey, "xxx")

	go func(ctx context.Context) {
		// 获得父协程传递的值
		if userId, ok := ctx.Value(userKey).(string); ok {
			assert.Equal(t, userId, "xxx")
		} else {
			fmt.Println("没有获得父协程传递的值")
		}
	}(ctx)

	// 等待子协程执行完毕
	time.Sleep(time.Second)
}

// 启动多个协程执行多任务
func WorkerPool(ctx context.Context, tasks <-chan int) {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case task, ok := <-tasks:
					if !ok {
						return
					}
					// 执行子任务
					processTask(ctx, task)
				case <-ctx.Done():
					// 接收到父协程取消信号
					return
				}
			}
		}()
	}
	wg.Wait()
}

// 处理任务
func processTask(ctx context.Context, task int) {
	select {
	case <-time.After(time.Second):
		// 阻塞操作直到 1 秒之后执行
		fmt.Printf("完成任务 %d\n", task)
	case <-ctx.Done():
		fmt.Printf("任务 %d 被父协程取消\n", task)
	}
}

func TestWaitGroup(t *testing.T) {
	// 创建一个 3 秒后超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 创建多个任务，任务执行完毕后关闭管道
	tasks := make(chan int, 10)
	go func() {
		for i := 0; i < 20; i++ {
			tasks <- i
		}
		close(tasks)
	}()

	// 执行多个任务
	WorkerPool(ctx, tasks)
}

func TestGetPriority(t *testing.T) {
	ctx := WithPriority(context.Background(), 2)

	priority, ok := GetPriority(ctx)
	if !ok {
		priority = 0 // 默认优先级
	}

	// 根据优先级执行不同操作
	switch priority {
	case 1:
		fmt.Println("高优先级")
	case 2:
		fmt.Println("中优先级")
	default:
		fmt.Println("低优先级")
	}
}

func TestWithValues(t *testing.T) {
	type testKey struct{}
	testValue := "value"
	ctx := context.WithValue(context.TODO(), testKey{}, testValue)

	ctx1 := NewContextWithValues(context.TODO(), ctx)
	require.NotNil(t, ctx1.Value(testKey{}))
	type notExist struct{}
	require.Nil(t, ctx1.Value(notExist{}))

	// 原始 context（带超时）
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 值优先的 context（带键值对）
	valuesCtx := context.WithValue(context.Background(), "key1", "value1")

	// 合并后的 context
	mergedCtx := NewContextWithValues(ctx, valuesCtx)

	// 优先从 valuesCtx 获取值
	assert.Equal(t, "value1", mergedCtx.Value("key1"))

	// 超时行为仍由 ctx 控制
	select {
	case <-mergedCtx.Done():
		fmt.Println("timed out")
	}
}
